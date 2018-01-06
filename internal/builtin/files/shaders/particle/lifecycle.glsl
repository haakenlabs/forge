#ifdef _COMPUTE_
layout(local_size_x = 128) in;

subroutine void TaskType();
subroutine uniform TaskType Task;

uniform vec4 u_start_color;
uniform vec3 u_angular_velocity_3d;
uniform vec3 u_rotation;
uniform vec3 u_position;
uniform uint u_random_seed = 1234;
uniform uint u_offset_in;
uniform uint u_offset_out;
uniform uint u_invocations;
uniform uint u_max_particles;
uniform float u_delta_time;
uniform float u_angular_velocity;
uniform float u_lifetime;
uniform float u_start_lifetime;
uniform float u_start_size;
uniform float u_velocity;

struct Particle {
    vec4 start_color;
    vec4 angular_velocity;
    vec4 direction;
    vec4 position;
    vec4 lifetime;
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


uint hashInt1D(uint x) {
    x += x >> 11;
    x ^= x << 7;
    x += x >> 15;
    x ^= x << 5;
    x += x >> 12;
    x ^= x << 9;

    return x;
}

uint hashInt2D(uint x, uint y) {
    x += x >> 11;
    x ^= x << 7;
    x += y;
    x ^= x << 6;
    x += x >> 15;
    x ^= x << 5;
    x += x >> 12;
    x ^= x << 9;

    return x;
}

uint hashInt3D(uint x, uint y, uint z) {
    x += x >> 11;
    x ^= x << 7;
    x += y;
    x ^= x << 3;
    x += z ^ ( x >> 14 );
    x ^= x << 6;
    x += x >> 15;
    x ^= x << 5;
    x += x >> 12;
    x ^= x << 9;

    return x;
}

float random(float f, uint seed) {
    const uint mantissa_mask = 0x007FFFFFu;
    const uint one = 0x3F800000u;

    uint h = hashInt2D(floatBitsToUint(f), seed);
    h &= mantissa_mask;
    h |= one;

    float  r2 = uintBitsToFloat(h);
    return r2 - 1.0;
}

vec3 random_vec3(float f, uint seed) {
    float k = random(f, seed);

    return normalize(vec3(random(k, seed), random(k+2, seed), random(k+1, seed)) * 2.0 - 1.0);
}

subroutine(TaskType)
void task_emit() {
    uint a_idx = atomicAdd(alive_count, 1) - 1;
    uint d_idx = atomicAdd(dead_count, -1);

    alive[u_offset_out+a_idx] = dead[d_idx];
    uint idx = alive[a_idx+u_offset_out];

    particles[idx].start_color = u_start_color;
    particles[idx].direction.xyz = vec3(0.0);
    particles[idx].direction.w = 10.0;
    particles[idx].position.xyz = vec3(0.0);
    particles[idx].angular_velocity.xyz = random_vec3(gl_GlobalInvocationID.x, u_random_seed) * 100.0;
    particles[idx].angular_velocity.w = 0;
    particles[idx].lifetime.x = u_start_lifetime;
    particles[idx].lifetime.y = u_start_lifetime;
    particles[idx].lifetime.z = 1.0;
}

subroutine(TaskType)
void task_lifetime() {
    uint index = gl_GlobalInvocationID.x;
    uint target = alive[u_offset_in+index];

    Particle p = particles[target];

    if (p.lifetime.x == 0) {
        uint a_idx = atomicAdd(alive_index, 1) - 1;
        alive[u_offset_out+a_idx] = target;
        return;
    }

    p.lifetime.y -= u_delta_time;
    p.lifetime.z = p.lifetime.y / p.lifetime.x;

    if (p.lifetime.y > 0) {
        uint a_idx = atomicAdd(alive_index, 1) - 1;
        alive[u_offset_out+a_idx] = target;
    } else {
        p.lifetime.y = 0;
        p.lifetime.z = 0;

        uint d_idx = atomicAdd(dead_count, 1) - 1;
        atomicAdd(alive_count, -1);

        dead[d_idx] = target;
    }

    particles[target] = p;
}

void main() {
    if (gl_GlobalInvocationID.x < u_invocations) {
        Task();
    }
}

#endif
