#version 330

in vec3 vertexPosition;
in vec2 vertexTexCoord;
in vec4 vertexColor;

out vec2 fragTexCoord;
out vec4 fragColor;
out float zDepth;

uniform mat4 mvp;
uniform vec2 uResolution; // Used for vertex snapping

void main()
{
    vec4 pos = mvp * vec4(vertexPosition, 1.0);
    
    // PS1 Vertex Snapping (Jitter)
    // Snap vertices to a "grid" based on the screen resolution
    // Note: To emulate PS1 perfectly, we perform perspective division manually here and snap.
    if (uResolution.x > 0.0) {
        vec2 snapped = pos.xy / pos.w;
        snapped = floor(snapped * uResolution * 0.5) / (uResolution * 0.5);
        pos.xy = snapped * pos.w;
    }

    gl_Position = pos;

    // PS1 Affine Texture Mapping
    // This removes the perspective-correction from the interpolator by multiplying UV by W,
    // which the fragment shader undoes to cause the "wobbling" effect.
    zDepth = pos.w;
    fragTexCoord = vertexTexCoord * pos.w; 
    fragColor = vertexColor;
}
