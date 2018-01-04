#ifdef _COMPUTE_
layout(local_size_x = 1000) in;

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

uniform uint c_max_particles;

uniform vec4 c_start_color;
uniform vec3 c_angular_velocity_3d;
uniform vec3 c_rotation;
uniform vec3 c_position;
uniform uint c_random_seed;
uniform float c_angular_velocity;
uniform float c_lifetime;
uniform float c_start_lifetime;
uniform float c_start_size;
uniform float c_velocity;

float rand(vec2 co){
    return fract(sin(dot(co.xy ,vec2(12.9898,78.233))) * 43758.5453);
}

void main()
{
    uint index = gl_GlobalInvocationID.x;

    if (index >= c_max_particles)
        return;

    float thing = float(index) / float(c_max_particles);

    thing = thing * 5.0 + rand(vec2(thing));

    Particle p;

    p.start_color = c_start_color;
    p.angular_velocity = vec4(c_angular_velocity_3d, c_angular_velocity);
    p.direction = vec4(c_rotation, 1.0f);
    p.position = vec4(vec3(thing), c_velocity);
    p.lifecycle = vec4(c_start_size, c_start_lifetime, c_lifetime, 1.0f);

    particles[index] = p;
}


#endif
