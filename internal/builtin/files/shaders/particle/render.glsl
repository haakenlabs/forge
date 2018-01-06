#ifdef _VERTEX_
uniform mat4 v_view_matrix;
uniform mat4 v_model_matrix;
uniform uint v_offset;

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

layout(std430, binding = 1) buffer alive_buffer {
    uint alive[];
};

out flat uint vo_index;

void main()
{
    vo_index = alive[v_offset+gl_VertexID];
    Particle p = particles[vo_index];

//    vo_color = vec4(p.position.xyz, 1.0) * 0.5 + 0.5;

    gl_Position = v_view_matrix * v_model_matrix * vec4(p.position.xyz, 1.0);
}

#endif

#ifdef _GEOMETRY_
layout(points) in;
layout(triangle_strip, max_vertices = 4) out;

out vec2 go_texture;

in flat uint vo_index[];
out flat uint go_index;

uniform mat4 v_projection_matrix;
uniform float g_quad_length = 0.02f;

void main()
{
    go_index = vo_index[0];

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

in flat uint go_index;

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

void main() {
//    vec3 color = vec3(1.0, 0.5, 0.1);

    Particle p = particles[go_index];

    vec3 color = mix(vec3(0.0, 0.0, 1.0), vec3(1.0, 0.0, 0.0), p.lifecycle.z);

    fo_attachment0 = vec4(color, texture(f_source_a, go_texture).a);
}

#endif
