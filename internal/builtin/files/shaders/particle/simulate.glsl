#ifdef _COMPUTE_
layout(local_size_x = 128) in;

const uint MODE_NORMAL = 0;
const uint MODE_REPELL = 1;
const uint MODE_BLACKHOLE = 2;
const uint MODE_GLOBAL = 3;

struct Particle {
    vec4 start_color;
    vec4 angular_velocity;
    vec4 direction;
    vec4 position;
    vec4 lifetime;
};

struct Attractor {
    vec4 position;
    vec4 direction;
    uint mode;
    float force;
    float range;
    float unused;
};

layout(std430, binding = 0) buffer particle_buffer {
    Particle particles[];
};

layout(std430, binding = 1) buffer alive_buffer {
    uint alive[];
};

layout(std430, binding = 2) buffer dead_buffer {
    uint dead[];
};

layout(std430, binding = 3) buffer index_buffer {
    uint alive_count;
    uint dead_count;
    uint alive_index;
};

layout(std430, binding = 4) buffer attractor_buffer {
    Attractor attractors[];
};

uniform uint u_invocations = 0;
uniform uint u_offset_out;
uniform uint u_attractor_count;
uniform int u_attractor_enable;
uniform float u_particle_mass = 0.1;
uniform float u_particle_inv_mass = 1.0 / 0.1;
uniform float u_delta_time = 0.0005;
uniform float u_max_distance = 45.0;

void main() {
    if (gl_GlobalInvocationID.x < u_invocations) {
        uint target = alive[gl_GlobalInvocationID.x+u_offset_out];

        Particle p = particles[target];

        vec3 force = vec3(0.0);

        if (u_attractor_enable != 0) {
            for (int i = 0; i < u_attractor_count; i++) {
                Attractor a = attractors[i];

                if (a.mode == MODE_GLOBAL) {
                    force += a.force;
                } else {
                    vec3 d = a.position.xyz - p.position.xyz;

                    if (a.range > 0) {
                        float dist = length(d);

                        if (a.mode == MODE_REPELL) {

                        } else {
                        
                        }
                    }

                }
            }
            // Temporary
            vec3 d = vec3(7, 0, 0) - p.position.xyz;
            float dist = length(d);

            if (dist < 0.2) {
                p.lifetime.y = 0;
                p.lifetime.z = 0;
            }

            force += (5000.0 / dist) * normalize(d);

            d = vec3(-5, 0, 0) - p.position.xyz;
            dist = length(d);

            if (dist < 0.2) {
                p.lifetime.y = 0;
                p.lifetime.z = 0;
            }

            force += (5000.0 / dist) * normalize(d);
        }

        vec3 a = force * u_particle_inv_mass;

        float d_time = u_delta_time / 100.0;

        p.position.xyz = p.position.xyz + p.angular_velocity.xyz * d_time + 0.5 * a * d_time * d_time;
        p.angular_velocity.xyz = p.angular_velocity.xyz + a * d_time;

        particles[target] = p;
    }
}

#endif