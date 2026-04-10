#version 330

in vec2 fragTexCoord;
in vec4 fragColor;
in vec3 fragNormal;

out vec4 finalColor;

uniform sampler2D texture0;
uniform vec2 uScroll;
uniform vec2 uOffset;
uniform vec2 uFlip;
uniform float uTime;

void main()
{
    vec2 uv = fragTexCoord;
    
    // Support flipping
    if (uFlip.x != 0.0) { uv.x *= sign(uFlip.x); }
    if (uFlip.y != 0.0) { uv.y *= sign(uFlip.y); }
    
    // Support scrolling
    uv += (uScroll * uTime);
    uv += uOffset;

    vec4 texelColor = texture(texture0, uv);
    if (texelColor.a < 0.1) discard;
    
    // Basic diffuse ambient placeholder
    vec3 lightDir = normalize(vec3(0.5, 1.0, 0.5));
    float diff = max(dot(normalize(fragNormal), lightDir), 0.0);
    vec3 ambient = vec3(0.3);
    
    finalColor = vec4(texelColor.rgb * (ambient + diff), texelColor.a) * fragColor;
}
