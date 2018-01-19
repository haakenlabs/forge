#ifdef _COMPUTE_
layout(local_size_x = 1000) in;

struct Particle
{
    vec4 start_color;
    vec4 angular_velocity;
    vec4 direction; // xyz: rotation, w: unused
    vec4 position; // xyz: position, w: velocity
    vec4 lifecycle; // x: start_size, y: start_lifetime, z: lifetime, w: unused
};

struct Attractor
{
    vec3 position;
    float force;
};

layout(std430, binding = 0) buffer particle_buffer 
{
    Particle particles[];
};
layout(std430, binding = 1) buffer alive_list
{
    uint alive[];
};
layout(std430, binding = 2) buffer dead_list
{
    uint dead[];
};

// layout(std430, binding = 1) buffer attractor_buffer 
// {
//     Attractor attractors[];
// };

uniform uint c_max_particles;
uniform uint c_attractors = 0;

uniform float c_particle_mass = 0.1;
uniform float c_particle_inv_mass = 1.0 / 0.1;
uniform float c_delta_time = 0.0005;
uniform float c_max_distance = 45.0;

uniform bool c_reset = false;

void main()
{
    uint index = gl_GlobalInvocationID.x;

    if (index >= c_max_particles)
        return;

    Particle p = particles[index];

    if (p.lifecycle.y > 0.0)
    {
        if (p.lifecycle.z < 0.0)
        {
            p.lifecycle.z = 0.0;
        } else {
            p.lifecycle.z -= c_delta_time;
        }
    }

    if (p.lifecycle.z > 0) {
        vec3 force = vec3(0.0);

        // Temporary
        vec3 d = vec3(5, 0, 0) - p.position.xyz;
        float dist = length(d);
        force += (5000.0 / dist) * normalize(d);

        d = vec3(-5, 0, 0) - p.position.xyz;
        dist = length(d);
        force += (5000.0 / dist) * normalize(d);

        if (c_reset)
        {
            p.position.xyz = vec3(0.0);
        }
        else
        {
            vec3 a = force * c_particle_inv_mass;
            float d_time = c_delta_time / 10.0;
            p.position.xyz = p.position.xyz + p.angular_velocity.xyz * d_time + 0.5 * a * d_time * d_time;
            p.angular_velocity.xyz = p.angular_velocity.xyz + a * d_time;
        }
    }

    // for (uint i = 0; i < c_attractors; i++)
    // {
    //     vec3 d = attractors[i].position - p.position;
    //     float dist = length(d);
    //     force += (attractors[i].force / dist) * normalize(d);
    // }

    particles[index] = p;
}

#endif
