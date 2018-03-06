#ifdef _VERTEX_
layout(location = 0) in vec3 vertex;
layout(location = 1) in vec3 normal;
layout(location = 2) in vec2 uv;

out vec3 vo_position;
out vec3 vo_normal;
out vec2 vo_texture;

uniform mat4 v_ortho_matrix;
uniform mat4 v_model_matrix;

void main()
{
    vo_position = vertex;
    vo_normal = normal;
    vo_texture = uv;

    gl_Position = v_ortho_matrix * v_model_matrix * vec4(vertex, 1.0);
}

#endif

#ifdef _FRAGMENT_
in vec3 vo_position;
in vec3 vo_normal;
in vec2 vo_texture;

out vec4 fo_color;

layout(binding = 0) uniform sampler2D f_source_a;

uniform vec4 f_color;
uniform float f_alpha;

void main()
{
    fo_color = vec4(f_color.rgb, texture(f_source_a, vo_texture).r * f_color.a * f_alpha);
}


#endif
