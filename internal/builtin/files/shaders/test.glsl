#ifdef _VERTEX_
layout(location = 0) in vec3 vertex;
layout(location = 1) in vec3 normal;
layout(location = 2) in vec2 uv;

out vec3 vo_position;
out vec3 vo_normal;
out vec2 vo_texture;

uniform mat4 v_mvp_matrix;
uniform mat4 v_projection_matrix;
uniform mat4 v_view_matrix;
uniform mat4 v_model_matrix;
uniform mat3 v_normal_matrix;

void main()
{
    vo_position = vertex;
    vo_normal = normal;
    vo_texture = uv;

    gl_Position = vec4(vertex, 1.0);
}

#endif

#ifdef _FRAGMENT_
in vec3 vo_position;
in vec3 vo_normal;
in vec2 vo_texture;

out vec4 fo_color;

layout(binding = 0) uniform sampler2D f_attachment0;

void main()
{
    fo_color = vec4(vo_texture, 0.0, 1.0);
}

#endif
