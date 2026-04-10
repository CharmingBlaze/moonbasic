#version 330

in vec2 fragTexCoord;
out vec4 fragColor;

uniform sampler2D sceneTexture;
uniform vec2 uResolution;
uniform float uTime;

void main()
{
    vec2 q = fragTexCoord;
    vec2 uv = q;
    
    // Curvature
    uv = uv * 2.0 - 1.0;
    vec2 offset = abs(uv.yx) / vec2(6.0, 4.0);
    uv = uv + uv * offset * offset;
    uv = uv * 0.5 + 0.5;

    // Bounds check
    if (uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0) {
        fragColor = vec4(0.0, 0.0, 0.0, 1.0);
        return;
    }
    
    // Chromatic aberration / Color bleeding
    vec3 col;
    col.r = texture(sceneTexture, vec2(uv.x + 0.002, uv.y)).r;
    col.g = texture(sceneTexture, uv).g;
    col.b = texture(sceneTexture, vec2(uv.x - 0.002, uv.y)).b;
    
    // Scanlines
    float s1 = sin(uv.y * uResolution.y * 3.1415);
    vec3 scanColor = col * vec3(s1 * 0.15 + 0.85);

    // Vignette
    float vig = (uv.x * (1.0 - uv.x) * uv.y * (1.0 - uv.y)) * 15.0;
    vig = pow(vig, 0.25);
    
    scanColor *= vig;
    fragColor = vec4(scanColor, 1.0);
}
