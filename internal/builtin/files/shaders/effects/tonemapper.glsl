#ifdef _FRAGMENT_

uniform mat3 rgb2xyz = mat3(
  0.4124564, 0.2126729, 0.0193339,
  0.3575761, 0.7151522, 0.1191920,
  0.1804375, 0.0721750, 0.9503041 );

uniform mat3 xyz2rgb = mat3(
  3.2404542, -0.9692660, 0.0556434,
  -1.5371385, 1.8760108, -0.2040259,
  -0.4985314, 0.0415560, 1.0572252 );

uniform vec3 lum_factor = vec3(0.2126, 0.7152, 0.0722);

uniform float f_exposure = 0.35;
uniform float f_average_gray;
uniform float f_white = 0.928;
uniform float f_adaptation_speed;
uniform vec2 f_image_size = vec2(1280, 720);

layout(binding = 2) uniform sampler2D u_avg_luminance;


vec3 gamma_correct(vec3 c) {
    vec3 color = c / (c + vec3(1.0));
    color = pow(color, vec3(1.0/2.2));

    return color;
}

subroutine(RenderPassType)
vec4 pass_avg_luminance()
{
    float lum = dot(texture(u_source, vo_texture).rgb, lum_factor);

    return vec4(log(lum + 0.00001));
}

subroutine(RenderPassType)
vec4 pass_basic()
{
    vec4 color = texture(u_source, vo_texture);
    float average_luminance = textureLod(u_avg_luminance, vo_texture, 8).r;
    average_luminance = exp(average_luminance / (f_image_size.x * f_image_size.y));

    // Convert to XYZ
    vec3 xyzCol = rgb2xyz * color.rgb;

    // Convert to xyY
    float xyzSum = xyzCol.x + xyzCol.y + xyzCol.z;
    vec3 xyYCol = vec3( xyzCol.x / xyzSum, xyzCol.y / xyzSum, xyzCol.y);

    // Apply the tone mapping operation to the luminance (xyYCol.z or xyzCol.y)
    float L = (f_exposure * xyYCol.z) / average_luminance;
    L = (L * (1 + L / (f_white * f_white) )) / (1 + L);

    // Using the new luminance, convert back to XYZ
    xyzCol.x = (L * xyYCol.x) / (xyYCol.y);
    xyzCol.y = L;
    xyzCol.z = (L * (1 - xyYCol.x - xyYCol.y))/xyYCol.y;

    return vec4(gamma_correct(xyz2rgb * xyzCol), 1.0);
}

subroutine(RenderPassType)
vec4 pass_basic_reinhard()
{
    // Not yet implemented
    return texture(u_source, vo_texture);
}

subroutine(RenderPassType)
vec4 pass_adaptive_reinhard()
{
    // Not yet implemented
    return texture(u_source, vo_texture);
}

#endif
