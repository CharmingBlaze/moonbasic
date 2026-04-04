//go:build cgo

package mbmodel3d

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Raylib LoadShaderFromMemory binds texture0..texture2; we patch locs for higher map slots.
const pbrVertexShader = `#version 330
in vec3 vertexPosition;
in vec2 vertexTexCoord;
in vec3 vertexNormal;
in vec4 vertexColor;
in vec4 vertexTangent;
uniform mat4 mvp;
uniform mat4 matModel;
uniform mat4 matNormal;
out vec3 fragPos;
out vec2 fragUV;
out vec4 fragCol;
out vec3 fragN;
out vec3 fragT;
out vec3 fragB;

void main() {
    vec4 wp = matModel * vec4(vertexPosition, 1.0);
    fragPos = wp.xyz;
    fragUV = vertexTexCoord;
    fragCol = vertexColor;
    mat3 nmat = mat3(matNormal);
    fragN = normalize(nmat * vertexNormal);
    fragT = normalize(nmat * vec3(vertexTangent));
    fragB = cross(fragN, fragT) * vertexTangent.w;
    gl_Position = mvp * vec4(vertexPosition, 1.0);
}
`

// Fragment shader: texture0 albedo, texture1 metalness, texture2 normal, texture3 roughness,
// texture11 depth shadow map (MATERIAL_MAP_BRDF slot). Custom uniforms set from Go before DrawMesh.
const pbrFragmentShader = `#version 330
in vec3 fragPos;
in vec2 fragUV;
in vec4 fragCol;
in vec3 fragN;
in vec3 fragT;
in vec3 fragB;
out vec4 finalColor;

uniform vec4 colDiffuse;
uniform sampler2D texture0;
uniform sampler2D texture1;
uniform sampler2D texture2;
uniform sampler2D texture3;
uniform sampler2D texture11;

uniform float roughnessValue;
uniform float metalnessValue;
uniform vec3 camPos;
uniform vec3 lightDir;
uniform vec3 lightColor;
uniform int useNormalMap;
uniform mat4 lightVP;
uniform int shadowEnabled;

float shadowFactor(vec3 N, vec3 L) {
    if (shadowEnabled == 0) return 1.0;
    vec4 lv = lightVP * vec4(fragPos, 1.0);
    vec3 proj = lv.xyz / lv.w;
    proj.xy = proj.xy * 0.5 + 0.5;
    if (proj.z > 1.0 || proj.x < 0.0 || proj.x > 1.0 || proj.y < 0.0 || proj.y > 1.0)
        return 1.0;
    float bias = max(0.0008 * (1.0 - dot(N, L)), 0.00015);
    float d = texture(texture11, proj.xy).r;
    return (proj.z - bias > d) ? 0.35 : 1.0;
}

void main() {
    vec4 albS = texture(texture0, fragUV);
    vec3 albedo = albS.rgb * colDiffuse.rgb * fragCol.rgb;
    float met = clamp(texture(texture1, fragUV).r * metalnessValue, 0.0, 1.0);
    float rough = clamp(texture(texture3, fragUV).r * roughnessValue, 0.04, 1.0);

    vec3 N = normalize(fragN);
    if (useNormalMap != 0) {
        vec3 t = normalize(fragT);
        vec3 b = normalize(fragB);
        mat3 TBN = mat3(t, b, N);
        vec3 tn = texture(texture2, fragUV).xyz * 2.0 - 1.0;
        N = normalize(TBN * tn);
    }

    vec3 V = normalize(camPos - fragPos);
    vec3 L = normalize(-lightDir);
    float NdotL = max(dot(N, L), 0.0);
    float NdotV = max(dot(N, V), 0.001);
    vec3 H = normalize(V + L);
    float NdotH = max(dot(N, H), 0.0);

    float alpha = rough * rough;
    float alpha2 = alpha * alpha;
    float denom = NdotH * NdotH * (alpha2 - 1.0) + 1.0;
    float D = alpha2 / (3.14159265 * denom * denom);
    float k = (rough + 1.0);
    k = k * k / 8.0;
    float G = (NdotV / (NdotV * (1.0 - k) + k)) * (NdotL / (NdotL * (1.0 - k) + k));
    vec3 F0 = mix(vec3(0.04), albedo, met);
    float VoH = max(dot(V, H), 0.0);
    vec3 F = F0 + (1.0 - F0) * pow(1.0 - VoH, 5.0);
    vec3 spec = (D * G) * F / (4.0 * NdotV * NdotL + 0.001);
    vec3 kS = F;
    vec3 kD = (1.0 - kS) * (1.0 - met);

    float sh = shadowFactor(N, L);
    vec3 radiance = lightColor * NdotL * sh;
    vec3 Lo = (kD * albedo / 3.14159265 + spec) * radiance;
    vec3 ambient = albedo * 0.06;
    vec3 color = ambient + Lo;
    color = color / (color + vec3(1.0));
    color = pow(color, vec3(1.0 / 2.2));
    finalColor = vec4(color, colDiffuse.a * fragCol.a);
}
`

func makePBRMaterial() rl.Material {
	mat := rl.LoadMaterialDefault()
	sh := rl.LoadShaderFromMemory(pbrVertexShader, pbrFragmentShader)
	if !rl.IsShaderValid(sh) {
		return mat
	}
	mat.Shader = sh
	patchStandardMapTextureLocs(&mat.Shader)
	mat.GetMap(rl.MapRoughness).Value = 1
	mat.GetMap(rl.MapMetalness).Value = 1
	return mat
}

func patchStandardMapTextureLocs(sh *rl.Shader) {
	for i := int32(3); i <= 11; i++ {
		name := fmt.Sprintf("texture%d", i)
		loc := rl.GetShaderLocation(*sh, name)
		if loc < 0 {
			continue
		}
		mapLoc := rl.ShaderLocMapAlbedo + i
		if mapLoc < rl.MaxShaderLocations {
			sh.UpdateLocation(mapLoc, loc)
		}
	}
}
