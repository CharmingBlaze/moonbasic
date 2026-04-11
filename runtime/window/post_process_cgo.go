//go:build cgo || (windows && !cgo)

package window

import (
	"fmt"
	"strings"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmodel3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

const postVS = `#version 330
in vec3 vertexPosition;
in vec2 vertexTexCoord;
out vec2 fragTexCoord;
uniform mat4 mvp;
void main() {
    fragTexCoord = vertexTexCoord;
    gl_Position = mvp * vec4(vertexPosition, 1.0);
}
`

const postFS = `#version 330
in vec2 fragTexCoord;
out vec4 fragColor;
uniform sampler2D sceneTexture;
uniform sampler2D depthTexture;
uniform vec2 resolution;

uniform int useChromatic;
uniform float chromaticOffset;
uniform int useSSAO;
uniform float ssaoRadius;
uniform float ssaoIntensity;
uniform int useSSR;
uniform int ssrSteps;
uniform float ssrStride;
uniform int useMotionBlur;
uniform float motionBlurStrength;
uniform int useDOF;
uniform float dofFocus;
uniform float dofRange;
uniform int useBloom;
uniform float bloomThreshold;
uniform float bloomIntensity;
uniform int useVignette;
uniform float vignetteStrength;
uniform int tonemapMode;
uniform int useSharpen;
uniform float sharpenAmount;
uniform int useGrain;
uniform float grainAmount;
uniform int useFXAA;

float luminance(vec3 c) {
    return dot(c, vec3(0.299, 0.587, 0.114));
}

float sampleDepth(vec2 uv) {
    return texture(depthTexture, uv).r;
}

vec3 fetchScene(vec2 uv) {
    return texture(sceneTexture, uv).rgb;
}

vec3 tonemapReinhard(vec3 c) {
    return c / (c + vec3(1.0));
}

vec3 tonemapFilmic(vec3 c) {
    vec3 x = max(c - vec3(0.004), vec3(0.0));
    return (x * (6.2 * x + vec3(0.5))) / (x * (6.2 * x + vec3(1.7)) + vec3(0.06));
}

vec3 tonemapACES(vec3 x) {
    const float a = 2.51, b = 0.03, c = 2.43, d = 0.59, e = 0.14;
    return clamp((x * (a * x + b)) / (x * (c * x + d) + e), 0.0, 1.0);
}

void main() {
    vec2 uv = fragTexCoord;
    vec2 px = vec2(1.0 / max(resolution.x, 1.0), 1.0 / max(resolution.y, 1.0));
    vec3 col;
    if (useChromatic != 0) {
        float o = chromaticOffset * 0.001;
        col.r = texture(sceneTexture, uv + vec2(o, 0.0)).r;
        col.g = texture(sceneTexture, uv).g;
        col.b = texture(sceneTexture, uv - vec2(o, 0.0)).b;
    } else {
        col = fetchScene(uv);
    }

    float d = sampleDepth(uv);

    if (useSSR != 0 && d < 0.999) {
        float cx = sampleDepth(uv + vec2(px.x, 0.0)) - sampleDepth(uv - vec2(px.x, 0.0));
        float cy = sampleDepth(uv + vec2(0.0, py.y)) - sampleDepth(uv - vec2(0.0, py.y));
        vec2 refDir = normalize(vec2(cx, cy) + vec2(0.0001));
        int steps = clamp(ssrSteps, 1, 32);
        for (int i = 1; i <= 32; i++) {
            if (i > steps) break;
            vec2 suv = uv + refDir * ssrStride * px * float(i);
            if (suv.x < 0.0 || suv.x > 1.0 || suv.y < 0.0 || suv.y > 1.0) break;
            col = mix(col, fetchScene(suv), 0.12);
        }
    }

    if (useSSAO != 0 && d < 0.999) {
        float ao = 0.0;
        const int SPI = 12;
        float noise = fract(sin(dot(uv, vec2(12.9898, 78.233))) * 43758.5453);
        for (int i = 0; i < SPI; i++) {
            float ang = 6.2831853 * (float(i) + noise) / float(SPI);
            float dist = fract(noise * float(i + 1)) * ssaoRadius;
            vec2 off = vec2(cos(ang), sin(ang)) * (dist * px.x);
            float nd = sampleDepth(uv + off);
            float delta = nd - d;
            ao += smoothstep(0.0, 0.02, delta) * ssaoIntensity;
        }
        col *= clamp(1.0 - ao / float(SPI), 0.0, 1.0);
    }

    if (useDOF != 0) {
        float blurAmt = clamp(abs(d - dofFocus) / max(dofRange, 0.0001), 0.0, 1.0);
        vec3 blur = (
            fetchScene(uv + vec2(px.x * 3.0, 0.0)) +
            fetchScene(uv - vec2(px.x * 3.0, 0.0)) +
            fetchScene(uv + vec2(0.0, py.y * 3.0)) +
            fetchScene(uv - vec2(0.0, py.y * 3.0))
        ) * 0.25;
        col = mix(col, blur, blurAmt);
    }

    if (useMotionBlur != 0) {
        vec2 dir = vec2(motionBlurStrength * 0.008, 0.0);
        vec3 acc = col;
        acc += fetchScene(uv + dir);
        acc += fetchScene(uv + dir * 2.0);
        acc += fetchScene(uv + dir * 3.0);
        acc += fetchScene(uv - dir);
        col = acc * 0.2;
    }

    if (useBloom != 0) {
        float lum = luminance(col);
        float bright = max(lum - bloomThreshold, 0.0);
        col += col * bright * bloomIntensity;
    }

    if (useFXAA != 0) {
        float lumaM = luminance(col);
        float lumaNW = luminance(fetchScene(uv + vec2(-px.x, -px.y)));
        float lumaNE = luminance(fetchScene(uv + vec2(px.x, -px.y)));
        float lumaSW = luminance(fetchScene(uv + vec2(-px.x, px.y)));
        float lumaSE = luminance(fetchScene(uv + vec2(px.x, px.y)));
        
        float lumaMin = min(lumaM, min(min(lumaNW, lumaNE), min(lumaSW, lumaSE)));
        float lumaMax = max(lumaM, max(max(lumaNW, lumaNE), max(lumaSW, lumaSE)));
        
        vec2 dir;
        dir.x = -((lumaNW + lumaNE) - (lumaSW + lumaSE));
        dir.y =  ((lumaNW + lumaSW) - (lumaNE + lumaSE));
        
        float dirReduce = max((lumaNW + lumaNE + lumaSW + lumaSE) * (0.25 * 0.125), 0.00001);
        float rcpDirMin = 1.0 / (min(abs(dir.x), abs(dir.y)) + dirReduce);
        dir = min(vec2(8.0, 8.0), max(vec2(-8.0, -8.0), dir * rcpDirMin)) * px;
        
        vec3 rgbA = 0.5 * (fetchScene(uv + dir * (1.0/3.0 - 0.5)) + fetchScene(uv + dir * (2.0/3.0 - 0.5)));
        vec3 rgbB = rgbA * 0.5 + 0.25 * (fetchScene(uv + dir * -0.5) + fetchScene(uv + dir * 0.5));
        float lumaB = luminance(rgbB);
        if ((lumaB < lumaMin) || (lumaB > lumaMax)) col = rgbA;
        else col = rgbB;
    }

    if (tonemapMode == 1) {
        col = tonemapReinhard(max(col, vec3(0.0)));
    } else if (tonemapMode == 2) {
        col = tonemapFilmic(max(col, vec3(0.0)));
    } else if (tonemapMode == 3) {
        col = tonemapACES(max(col, vec3(0.0)));
    }

    if (useSharpen != 0) {
        vec3 blur = (
            fetchScene(uv + vec2(px.x, 0.0)) +
            fetchScene(uv - vec2(px.x, 0.0)) +
            fetchScene(uv + vec2(0.0, py.y)) +
            fetchScene(uv - vec2(0.0, py.y))
        ) * 0.25;
        col = col + (col - blur) * sharpenAmount;
    }

    if (useGrain != 0) {
        float n = fract(sin(dot(uv * resolution, vec2(12.9898, 78.233))) * 43758.5453);
        col += (n - 0.5) * grainAmount;
    }

    if (useVignette != 0) {
        vec2 p = uv * 2.0 - 1.0;
        float v = 1.0 - dot(p, p) * vignetteStrength * 0.2;
        col *= clamp(v, 0.0, 1.0);
    }

    col = clamp(col, 0.0, 1.0);
    fragColor = vec4(col, 1.0);
}
`

var (
	postMu             sync.Mutex
	postActive         bool
	deferredPipeline   bool
	postBloom          bool
	postVignette       bool
	postChromatic      bool
	postSSAO           bool
	postSSR            bool
	postMotionBlur     bool
	postDOF            bool
	postSharpen        bool
	postGrain          bool
	postFXAA           bool
	postCustom         rl.Shader
	postCustomOn       bool
	postSceneRT        rl.RenderTexture2D
	postRTW, postRTH   int32
	postBuiltIn        rl.Shader
	postBuiltInLoaded  bool
	postTonemapMode    int32 // 0 none, 1 reinhard, 2 filmic, 3 aces
	postKV             = map[string]float32{
		"bloom.threshold":      0.8,
		"bloom.intensity":      1.2,
		"vignette.strength":    0.6,
		"chromatic.offset":     3.0,
		"ssao.radius":          12.0,
		"ssao.intensity":       0.7,
		"ssr.steps":            12.0,
		"ssr.stride":           4.0,
		"motionblur.strength":  1.0,
		"dof.focus":            0.5,
		"dof.range":            0.25,
		"sharpen.amount":       0.35,
		"grain.amount":         0.04,
	}
	postCapturing bool
	uFloat1       = make([]float32, 1)
	uFloat2       = make([]float32, 2)
)

func setRenderPipelineMode(mode string) {
	postMu.Lock()
	defer postMu.Unlock()
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "deferred":
		deferredPipeline = true
	case "forward":
		deferredPipeline = false
	}
	postRebuildActiveLocked()
}

func postCaptureEnabled() bool {
	return postActive || deferredPipeline
}

// postRebuildActiveLocked sets postActive from flags. "Tonemap mode 0" alone must not force the
// offscreen post pass — that path (RT + fullscreen shader) has caused black frames on some Intel
// drivers when no actual effect is enabled. Hold postMu.
func postRebuildActiveLocked() {
	postActive = deferredPipeline || postCustomOn || postBloom || postVignette || postChromatic ||
		postSSAO || postSSR || postMotionBlur || postDOF || postSharpen || postGrain || postFXAA ||
		postTonemapMode != 0
}

func (m *Module) registerPostCommands(r runtime.Registrar) {
	r.Register("POST.ADD", "post", m.postAdd)
	r.Register("POST.REMOVE", "post", m.postRemove)
	r.Register("POST.SETPARAM", "post", m.postSetParam)
	r.Register("POST.SETTONEMAP", "post", m.postSetTonemap)
	r.Register("RENDER.SETTONEMAPPING", "render", m.postSetTonemap)
	r.Register("POST.ADDSHADER", "post", m.postAddShader)

	r.Register("POST.BLOOM", "post", m.postBloomShorthand)
	r.Register("POST.VIGNETTE", "post", m.postVignetteShorthand)
	r.Register("POST.CHROMATIC", "post", m.postChromaticShorthand)
	r.Register("RENDER.SETPOSTPROCESS", "render", runtime.AdaptLegacy(m.rSetPostProcess))
}

func (m *Module) rSetPostProcess(args []value.Value) (value.Value, error) {
	if len(args) < 1 { return value.Nil, fmt.Errorf("RENDER.SETPOSTPROCESS (id, [enable])") }
	id, ok := args[0].ToInt()
	if !ok { return value.Nil, fmt.Errorf("invalid post process id") }
	on := true
	if len(args) > 1 {
		on = valueTruthy(args[1])
	}
	// Use embedded shader library enum IDs mapping to post passes
	_ = on
	_ = id
	return value.Nil, nil
}

func (m *Module) postBloomShorthand(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("POST.BLOOM expects (enable, threshold, intensity)")
	}
	on := valueTruthy(args[0])
	postMu.Lock()
	postBloom = on
	if on {
		if f, ok := args[1].ToFloat(); ok {
			postKV["bloom.threshold"] = float32(f)
		}
		if f, ok := args[2].ToFloat(); ok {
			postKV["bloom.intensity"] = float32(f)
		}
	}
	postRebuildActiveLocked()
	if postActive {
		ensureBuiltInPostShader()
	}
	postMu.Unlock()
	return value.Nil, nil
}

func (m *Module) postVignetteShorthand(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("POST.VIGNETTE expects (enable, strength)")
	}
	on := valueTruthy(args[0])
	postMu.Lock()
	postVignette = on
	if on {
		if f, ok := args[1].ToFloat(); ok {
			postKV["vignette.strength"] = float32(f)
		}
	}
	postRebuildActiveLocked()
	if postActive {
		ensureBuiltInPostShader()
	}
	postMu.Unlock()
	return value.Nil, nil
}

func (m *Module) postChromaticShorthand(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("POST.CHROMATIC expects (enable, offset, [unused])")
	}
	on := valueTruthy(args[0])
	postMu.Lock()
	postChromatic = on
	if on {
		if f, ok := args[1].ToFloat(); ok {
			postKV["chromatic.offset"] = float32(f) * 1000.0 // scale back because inner shader uses 0.001
		}
	}
	postRebuildActiveLocked()
	if postActive {
		ensureBuiltInPostShader()
	}
	postMu.Unlock()
	return value.Nil, nil
}

func (m *Module) postAdd(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("POST.ADD expects 1 string (bloom, vignette, chromatic)")
	}
	name, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	postMu.Lock()
	defer postMu.Unlock()
	postCustomOn = false
	switch name {
	case "bloom":
		postBloom = true
	case "vignette":
		postVignette = true
	case "chromatic":
		postChromatic = true
	case "ssao":
		postSSAO = true
	case "ssr":
		postSSR = true
	case "motionblur":
		postMotionBlur = true
	case "dof":
		postDOF = true
	case "sharpen":
		postSharpen = true
	case "grain":
		postGrain = true
	case "fxaa":
		postFXAA = true
	default:
		return value.Nil, fmt.Errorf("POST.ADD: unknown effect %q", name)
	}
	postRebuildActiveLocked()
	if postActive {
		ensureBuiltInPostShader()
	}
	return value.Nil, nil
}

func (m *Module) postRemove(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("POST.REMOVE expects 1 string (bloom, vignette, etc)")
	}
	name, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	postMu.Lock()
	defer postMu.Unlock()
	switch name {
	case "bloom":
		postBloom = false
	case "vignette":
		postVignette = false
	case "chromatic":
		postChromatic = false
	case "ssao":
		postSSAO = false
	case "ssr":
		postSSR = false
	case "motionblur":
		postMotionBlur = false
	case "dof":
		postDOF = false
	case "sharpen":
		postSharpen = false
	case "grain":
		postGrain = false
	case "fxaa":
		postFXAA = false
	case "all":
		postBloom, postVignette, postChromatic, postSSAO, postSSR = false, false, false, false, false
		postMotionBlur, postDOF, postSharpen, postGrain, postFXAA = false, false, false, false, false
		postTonemapMode = 0
	default:
		return value.Nil, fmt.Errorf("POST.REMOVE: unknown effect %q", name)
	}
	postRebuildActiveLocked()
	return value.Nil, nil
}

func (m *Module) postSetTonemap(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("POST.SETTONEMAP expects (mode): 0=None, 1=Reinhard, 2=Filmic, 3=ACES")
	}
	mode, ok := args[0].ToInt()
	if !ok {
		return value.Nil, fmt.Errorf("POST.SETTONEMAP: mode must be numeric")
	}
	postMu.Lock()
	postTonemapMode = int32(mode)
	postRebuildActiveLocked()
	if postActive {
		ensureBuiltInPostShader()
	}
	postMu.Unlock()
	return value.Nil, nil
}

func (m *Module) postSetParam(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("POST.SETPARAM expects (pass, key, value)")
	}
	pass, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	key, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	var v float32
	if f, ok := args[2].ToFloat(); ok {
		v = float32(f)
	} else if i, ok := args[2].ToInt(); ok {
		v = float32(i)
	} else {
		return value.Nil, fmt.Errorf("POST.SETPARAM: value must be numeric")
	}
	postMu.Lock()
	postKV[pass+"."+key] = v
	postMu.Unlock()
	return value.Nil, nil
}

func (m *Module) postAddShader(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if rt == nil || rt.Heap == nil {
		return value.Nil, fmt.Errorf("POST.ADDSHADER: heap not available")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("POST.ADDSHADER expects shader handle")
	}
	sh, err := mbmodel3d.ShaderRaylib(rt.Heap, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	postMu.Lock()
	defer postMu.Unlock()
	postActive = true
	postCustomOn = true
	postCustom = sh
	postBloom, postVignette, postChromatic = false, false, false
	return value.Nil, nil
}

func ensureBuiltInPostShader() {
	if postBuiltInLoaded {
		return
	}
	postBuiltIn = rl.LoadShaderFromMemory(postVS, postFS)
	postBuiltInLoaded = true
}

func postRenderTargetBegin(c rl.Color) bool {
	postMu.Lock()
	defer postMu.Unlock()
	postCapturing = false
	if !postCaptureEnabled() {
		return false
	}
	w, h := int32(rl.GetRenderWidth()), int32(rl.GetRenderHeight())
	if w <= 0 || h <= 0 {
		return false
	}
	if postSceneRT.Texture.ID == 0 || postRTW != w || postRTH != h {
		if postSceneRT.Texture.ID != 0 {
			rl.UnloadRenderTexture(postSceneRT)
		}
		postSceneRT = rl.LoadRenderTexture(w, h)
		postRTW, postRTH = w, h
	}
	rl.BeginTextureMode(postSceneRT)
	rl.ClearBackground(c)
	postCapturing = true
	return true
}

func postRenderTargetPresent() {
	postMu.Lock()
	if !postCaptureEnabled() || !postCapturing {
		postMu.Unlock()
		return
	}
	postCapturing = false
	custom := postCustomOn
	bi := postBuiltIn
	bloom := postBloom
	vig := postVignette
	chr := postChromatic
	ssao := postSSAO
	ssr := postSSR
	mb := postMotionBlur
	dof := postDOF
	shrp := postSharpen
	grain := postGrain
	tm := postTonemapMode
	th := postKV["bloom.threshold"]
	inten := postKV["bloom.intensity"]
	vstr := postKV["vignette.strength"]
	coff := postKV["chromatic.offset"]
	ssaoR := postKV["ssao.radius"]
	ssaoI := postKV["ssao.intensity"]
	ssrS := int32(postKV["ssr.steps"])
	if ssrS < 1 {
		ssrS = 1
	}
	ssrSt := postKV["ssr.stride"]
	mbStr := postKV["motionblur.strength"]
	dofF := postKV["dof.focus"]
	dofR := postKV["dof.range"]
	shAmt := postKV["sharpen.amount"]
	gAmt := postKV["grain.amount"]
	rt := postSceneRT
	postMu.Unlock()

	rl.EndTextureMode()

	w := float32(rl.GetRenderWidth())
	h := float32(rl.GetRenderHeight())
	tex := rt.Texture
	depthTex := rt.Depth

	if custom {
		rl.BeginShaderMode(postCustom)
		loc := rl.GetShaderLocation(postCustom, "texture0")
		if loc >= 0 {
			rl.SetShaderValueTexture(postCustom, loc, tex)
		}
		rl.DrawTexturePro(tex, rl.NewRectangle(0, 0, float32(tex.Width), -float32(tex.Height)), rl.NewRectangle(0, 0, w, h), rl.Vector2Zero(), 0, rl.White)
		rl.EndShaderMode()
		return
	}

	ensureBuiltInPostShader()
	sh := bi
	rl.BeginShaderMode(sh)
	loc := rl.GetShaderLocation(sh, "sceneTexture")
	if loc >= 0 {
		rl.SetShaderValueTexture(sh, loc, tex)
	}
	locD := rl.GetShaderLocation(sh, "depthTexture")
	if locD >= 0 {
		rl.SetShaderValueTexture(sh, locD, depthTex)
	}
	setPI := func(n string, v int32) {
		l := rl.GetShaderLocation(sh, n)
		if l >= 0 {
			uFloat1[0] = float32(v)
			rl.SetShaderValue(sh, l, uFloat1, rl.ShaderUniformInt)
		}
	}
	setPF := func(n string, v float32) {
		l := rl.GetShaderLocation(sh, n)
		if l >= 0 {
			uFloat1[0] = v
			rl.SetShaderValue(sh, l, uFloat1, rl.ShaderUniformFloat)
		}
	}
	setPV2 := func(n string, vx, vy float32) {
		l := rl.GetShaderLocation(sh, n)
		if l >= 0 {
			uFloat2[0] = vx
			uFloat2[1] = vy
			rl.SetShaderValue(sh, l, uFloat2, rl.ShaderUniformVec2)
		}
	}

	setPV2("resolution", w, h)

	setPI("useBloom", boolAsInt(bloom))
	setPF("bloomThreshold", th)
	setPF("bloomIntensity", inten)
	setPI("useVignette", boolAsInt(vig))
	setPF("vignetteStrength", vstr)
	setPI("useChromatic", boolAsInt(chr))
	setPF("chromaticOffset", coff)
	setPI("useSSAO", boolAsInt(ssao))
	setPF("ssaoRadius", ssaoR)
	setPF("ssaoIntensity", ssaoI)
	setPI("useSSR", boolAsInt(ssr))
	setPI("ssrSteps", ssrS)
	setPF("ssrStride", ssrSt)
	setPI("useMotionBlur", boolAsInt(mb))
	setPF("motionBlurStrength", mbStr)
	setPI("useDOF", boolAsInt(dof))
	setPF("dofFocus", dofF)
	setPF("dofRange", dofR)
	setPI("tonemapMode", tm)
	setPI("useSharpen", boolAsInt(shrp))
	setPF("sharpenAmount", shAmt)
	setPI("useGrain", boolAsInt(grain))
	setPF("grainAmount", gAmt)
	setPI("useFXAA", boolAsInt(postFXAA))

	rl.DrawTexturePro(tex, rl.NewRectangle(0, 0, float32(tex.Width), -float32(tex.Height)), rl.NewRectangle(0, 0, w, h), rl.Vector2Zero(), 0, rl.White)
	rl.EndShaderMode()
}

func boolAsInt(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func valueTruthy(v value.Value) bool {
	switch v.Kind {
	case value.KindBool:
		return v.IVal != 0
	case value.KindInt:
		return v.IVal != 0
	case value.KindFloat:
		return v.FVal != 0
	default:
		return false
	}
}
