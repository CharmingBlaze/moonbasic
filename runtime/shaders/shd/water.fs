#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

out vec4 finalColor;

uniform sampler2D texture0;
uniform float uTime;
uniform vec2 uScroll;

void main()
{
    vec2 uv = fragTexCoord;
    
    // UV distortion based on sine waves mimicking water displacement
    uv.x += sin(uv.y * 10.0 + uTime) * 0.02;
    uv.y += cos(uv.x * 10.0 + uTime) * 0.02;
    
    uv += (uScroll * uTime);

    vec4 texelColor = texture(texture0, uv);
    
    // Slight transparency scaling from alpha
    finalColor = vec4(texelColor.rgb, texelColor.a * 0.8) * fragColor;
}
