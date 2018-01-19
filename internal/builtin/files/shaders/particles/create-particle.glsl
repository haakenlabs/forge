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

float RadicalInverse_VdC(uint bits)
{
    bits = (bits << 16u) | (bits >> 16u);
    bits = ((bits & 0x55555555u) << 1u) | ((bits & 0xAAAAAAAAu) >> 1u);
    bits = ((bits & 0x33333333u) << 2u) | ((bits & 0xCCCCCCCCu) >> 2u);
    bits = ((bits & 0x0F0F0F0Fu) << 4u) | ((bits & 0xF0F0F0F0u) >> 4u);
    bits = ((bits & 0x00FF00FFu) << 8u) | ((bits & 0xFF00FF00u) >> 8u);

    return float(bits) * 2.3283064365386963e-10; // / 0x100000000
}

// ----------------------------------------------------------------------------
vec2 Hammersley(uint i, uint N)
{
    return vec2(float(i)/float(N), RadicalInverse_VdC(i));
}

void main()
{
    uint index = gl_GlobalInvocationID.x;

    if (index >= c_max_particles)
        return;

    vec2 s0 = Hammersley(index, c_max_particles);
    vec2 s1 = Hammersley(c_max_particles-index, c_max_particles);

    Particle p;

    p.start_color = c_start_color;
    p.angular_velocity = vec4(c_angular_velocity_3d, c_angular_velocity);
    p.direction = vec4(c_rotation, 1.0f);
    p.position = vec4(vec3(s0, s1.y)*5.0, c_velocity);
    p.lifecycle = vec4(c_start_size, c_start_lifetime, c_start_lifetime, 0.0f);

    particles[index] = p;
}


#endif
