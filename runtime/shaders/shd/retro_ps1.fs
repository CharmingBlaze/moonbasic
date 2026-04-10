#version 330

in vec2 fragTexCoord;
in vec4 fragColor;
in float zDepth;

out vec4 finalColor;

uniform sampler2D texture0;

void main()
{
    // PS1 Affine mapping inverse:
    // Our vertex shader sent (uv * pos.w) and (pos.w).
    // The GPU interpolator bilinearly interpolated both linearly across the primitive.
    // By dividing the UV by the interpolated W here, we perfectly emulate the lack of perspective correction!
    vec2 affineUV = fragTexCoord / zDepth;
    
    vec4 texelColor = texture(texture0, affineUV);
    
    // Nearest-neighbor emulation if Bilinear filter is accidentally on
    // In actual implementation we just set the texture to FILTER_POINT in Raylib.
    
    finalColor = texelColor * fragColor;
}
