//go:build cgo || (windows && !cgo)

package mbmodel3d

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ShaderLibrary manages the "Essential Five" studio shaders.
type ShaderLibrary struct {
	shaders map[string]rl.Shader
	mu      sync.Mutex
}

var globalShaderLib = &ShaderLibrary{
	shaders: make(map[string]rl.Shader),
}

const (
	EffectStandard  = "STANDARD"
	EffectWater     = "WATER"
	EffectToon      = "TOON"
	EffectAtmosphere = "ATMOSPHERE"
	EffectDissolve  = "DISSOLVE"
)

// GetShader returns a shared shader instance for the given effect.
func (l *ShaderLibrary) GetShader(effect string) rl.Shader {
	l.mu.Lock()
	defer l.mu.Unlock()

	if sh, ok := l.shaders[effect]; ok {
		return sh
	}

	var sh rl.Shader
	switch effect {
	case EffectStandard:
		sh = rl.LoadShaderFromMemory(pbrVertexShader, pbrFragmentShader)
	case EffectWater:
		sh = rl.LoadShaderFromMemory(pbrVertexShader, waterFragmentShader)
	case EffectToon:
		sh = rl.LoadShaderFromMemory(pbrVertexShader, toonFragmentShader)
	case EffectAtmosphere:
		sh = rl.LoadShaderFromMemory(pbrVertexShader, atmosphereFragmentShader)
	case EffectDissolve:
		sh = rl.LoadShaderFromMemory(pbrVertexShader, dissolveFragmentShader)
	default:
		return rl.Shader{}
	}

	if rl.IsShaderValid(sh) {
		patchStandardMapTextureLocs(&sh)
		l.shaders[effect] = sh
	}
	return sh
}

const waterFragmentShader = `
#version 330
in vec3 fragPos;
in vec2 fragUV;
in vec4 fragCol;
in vec3 fragN;
out vec4 finalColor;

uniform vec4 colDiffuse;
uniform sampler2D texture0; // Normals
uniform float time;
uniform vec3 camPos;

void main() {
    vec2 uv = fragUV + vec2(time * 0.05, time * 0.02);
    vec3 n = texture(texture0, uv).rgb * 2.0 - 1.0;
    vec3 V = normalize(camPos - fragPos);
    float fresnel = pow(1.0 - max(dot(n, V), 0.0), 5.0);
    
    vec3 waterCol = mix(colDiffuse.rgb, vec3(0.1, 0.4, 0.8), 0.5);
    finalColor = vec4(mix(waterCol, vec3(1.0), fresnel * 0.5), 0.8);
}
`

const toonFragmentShader = `
#version 330
in vec3 fragPos;
in vec2 fragUV;
in vec4 fragCol;
in vec3 fragN;
out vec4 finalColor;

uniform vec4 colDiffuse;
uniform sampler2D texture0;
uniform vec3 lightDir;

void main() {
    vec4 alb = texture(texture0, fragUV) * colDiffuse * fragCol;
    float intensity = dot(normalize(fragN), normalize(-lightDir));
    
    if (intensity > 0.95) intensity = 1.0;
    else if (intensity > 0.5) intensity = 0.7;
    else if (intensity > 0.2) intensity = 0.3;
    else intensity = 0.1;
    
    finalColor = vec4(alb.rgb * intensity, alb.a);
}
`

const atmosphereFragmentShader = `
#version 330
in vec3 fragPos;
in vec2 fragUV;
in vec4 fragCol;
in vec3 fragN;
out vec4 finalColor;

uniform vec4 colDiffuse;
uniform sampler2D texture0;
uniform vec3 camPos;
uniform float fogDensity;

void main() {
    vec4 alb = texture(texture0, fragUV) * colDiffuse * fragCol;
    float dist = length(camPos - fragPos);
    float fog = exp(-pow(dist * fogDensity, 2.0));
    
    finalColor = mix(vec4(0.5, 0.6, 0.7, 1.0), alb, clamp(fog, 0.0, 1.0));
}
`

const dissolveFragmentShader = `
#version 330
in vec3 fragPos;
in vec2 fragUV;
in vec4 fragCol;
in vec3 fragN;
out vec4 finalColor;

uniform vec4 colDiffuse;
uniform sampler2D texture0;
uniform sampler2D texture1; // Noise
uniform float dissolveAmount;

void main() {
    float noise = texture(texture1, fragUV).r;
    if (noise < dissolveAmount) discard;
    
    vec4 alb = texture(texture0, fragUV) * colDiffuse * fragCol;
    float edge = smoothstep(dissolveAmount, dissolveAmount + 0.05, noise);
    finalColor = mix(vec4(1.0, 0.5, 0.0, 1.0), alb, edge);
}
`
