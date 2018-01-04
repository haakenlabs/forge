#ifdef _VERTEX_
layout(location = 0) in vec3 vertex;
layout(location = 1) in vec3 normal;
layout(location = 2) in vec2 uv;

out vec3 vo_eye;

uniform mat4 v_view_matrix;
uniform mat4 v_projection_matrix;

void main()
{
    mat4 inverse_projection = inverse(v_projection_matrix);
    mat3 inverse_view = transpose(mat3(v_view_matrix));

    vec3 unprojected = vec4(inverse_projection * vec4(vertex, 1.0)).xyz;
    vo_eye = inverse_view * unprojected;

    gl_Position = vec4(vertex, 1.0);
}

#endif

#ifdef _FRAGMENT_
in vec3 vo_eye;

out vec4 fo_color;

layout(binding = 0) uniform samplerCube f_environment;

void main()
{
    fo_color = texture(f_environment, vo_eye);
}

#endif
