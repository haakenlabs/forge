#ifdef _VERTEX_
layout(location = 0) in vec3 vertex;
layout(location = 1) in vec3 normal;
layout(location = 2) in vec2 uv;

out vec3 vo_position;

uniform mat4 v_projection_matrix;
uniform mat4 v_view_matrix;

void main()
{
    vo_position = vec3(v_projection_matrix * v_view_matrix * vec4(vertex, 1.0));

    gl_Position = vec4(vec3(vertex.x * -1.0, vertex.yz), 1.0);
}

#endif

#ifdef _FRAGMENT_
#define M_PI 3.141592653589

in vec3 vo_position;

out vec4 fo_color;

layout(binding = 0) uniform sampler2D f_attachment0;

const vec2 invAtan = vec2(0.1591, 0.3183);

vec2 SampleSphericalMap(vec3 v) {
    vec2 uv = vec2(atan(v.z, v.x), asin(v.y));
    uv *= invAtan;
    uv += 0.5;

    return uv;
}

void main()
{
    vec2 uv = SampleSphericalMap(normalize(vo_position));
    vec3 color = texture(f_attachment0, uv).rgb;

    fo_color = vec4(color, 1.0);
}

#endif
