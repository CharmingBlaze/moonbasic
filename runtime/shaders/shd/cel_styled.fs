#version 330

in vec2 fragTexCoord;
in vec4 fragColor;
in vec3 fragNormal;

out vec4 finalColor;

uniform sampler2D texture0;
uniform vec3 uLightDir = vec3(-0.5, -1.0, -0.5);

void main()
{
    vec4 texel = texture(texture0, fragTexCoord);
    if (texel.a < 0.1) discard;

    vec3 n = normalize(fragNormal);
    vec3 l = normalize(-uLightDir);
    float NdotL = max(dot(n, l), 0.0);
    
    // Light stepping for toon cel-shading
    float intensity;
    if (NdotL > 0.95) intensity = 1.0;
    else if (NdotL > 0.5) intensity = 0.7;
    else if (NdotL > 0.25) intensity = 0.4;
    else intensity = 0.2;
    
    finalColor = vec4(texel.rgb * intensity, texel.a) * fragColor;
}
