#version 330

in vec2 fragTexCoord;
out vec4 fragColor;

uniform sampler2D sceneTexture;
uniform vec2 uResolution;
uniform float uTime;
uniform float uPixelSize = 4.0;

void main()
{
    // Pixelate
    vec2 pos = fragTexCoord * uResolution;
    // Map to grid
    pos = floor(pos / uPixelSize) * uPixelSize;
    // Back to UV
    vec2 uv = pos / uResolution;

    vec4 texel = texture(sceneTexture, uv);
    fragColor = texel;
}
