#ifdef _VERTEX_
layout(location = 0) in vec4 color;
layout(location = 1) in vec4 angular_vel;
layout(location = 2) in vec4 rotation;
layout(location = 3) in vec4 position;
layout(location = 4) in vec4 lifecycle;

uniform mat4 v_view_matrix;
uniform mat4 v_model_matrix;

out flat uint go_vertex_id;

void main()
{
    go_vertex_id = gl_VertexID;
    gl_Position = v_view_matrix * v_model_matrix * vec4(position.xyz, 1.0);
}

#endif

#ifdef _GEOMETRY_
layout(points) in;
layout(triangle_strip, max_vertices = 4) out;

out vec2 go_texture;

in flat uint go_vertex_id[];
out flat uint vo_vertex_id;


uniform mat4 v_projection_matrix;
uniform float g_quad_length = 0.02f;

void main()
{
    mat4 m = v_projection_matrix;

    gl_Position = m * (vec4(-g_quad_length, -g_quad_length, 0.0, 0.0) + gl_in[0].gl_Position);
    go_texture = vec2(0.0, 0.0);
    EmitVertex();

    gl_Position = m * (vec4(g_quad_length, -g_quad_length, 0.0, 0.0) + gl_in[0].gl_Position);
    go_texture = vec2(1.0, 0.0);
    EmitVertex();

    gl_Position = m * (vec4(-g_quad_length, g_quad_length, 0.0, 0.0) + gl_in[0].gl_Position);
    go_texture = vec2(0.0, 1.0);
    EmitVertex();

    gl_Position = m * (vec4(g_quad_length, g_quad_length, 0.0, 0.0) + gl_in[0].gl_Position);
    go_texture = vec2(1.0, 1.0);
    EmitVertex();
}

#endif


#ifdef _FRAGMENT_
in vec2 go_texture;

layout(location = 0) out vec4 fo_attachment0;

layout(binding = 0) uniform sampler2D f_source_a;

uniform vec4 f_color;
uniform float f_time;

in flat uint vo_vertex_id;

struct Particle
{
    vec4 start_color;
    vec4 angular_velocity;
    vec4 direction;
    vec4 position;
    vec4 lifecycle;
};

layout(std430, binding = 0) buffer particle_buffer {
    Particle particles[];
};


void main()
{
    Particle p = particles[vo_vertex_id];

    if (p.lifecycle.y > 0.0)
        if (p.lifecycle.z <= 0.0)
            discard;

    vec4 color = texture(f_source_a, go_texture);

    color.r = color.r + -0.5 + cos(f_time * 0.4 + 1.5) * 0.5;
    color.g = color.g + -0.5 + cos(f_time * 0.6) * sin(f_time * 0.3) * 0.35;
    color.b = color.b + -0.5 + sin(f_time * 0.2) * 0.5;
    color.a *= 0.4;

    fo_attachment0 = color;
}

#endif
