#ifdef _COMPUTE_
layout(local_size_x = 1000) in;

subroutine vec3 ShapeType(uint i);
subroutine uniform ShapeType EmitShape;

struct Particle
{
    vec4 start_color;
    vec4 angular_velocity;
    vec4 direction;
    vec4 position;
    vec4 lifetime;
};

layout(std430, binding = 0) buffer particle_buffer {
    Particle particles[];
};

uniform uint c_max_particles;

uniform vec4 c_start_color;
uniform vec3 c_angular_velocity_3d;
uniform vec3 c_rotation;
uniform vec3 c_position;
uniform float c_angular_velocity;
uniform float c_lifetime;
uniform float c_start_lifetime;
uniform float c_start_size;
uniform float c_velocity;

subroutine(ShapeType)
vec3 shape_sphere(uint i) {
    return vec3(0, 0, 0);
}

subroutine(ShapeType)
vec3 shape_cone(uint i) {
    return vec3(0, 0, 0);
}

subroutine(ShapeType)
vec3 shape_box(uint i) {
    return vec3(0, 0, 0);
}

void main() {
    uint index = gl_GlobalInvocationID.x;

    Particle p;

    if (index >= c_max_particles)
        return;

    p = particles[index];

    if (p.lifetime.z > 0) {
        return;
    }

    p.start_color = c_start_color;
    p.angular_velocity = vec4(c_angular_velocity_3d, c_angular_velocity);
    p.direction = vec4(EmitShape(index), 1.0f);
    p.position = vec4(c_position, c_velocity);
    p.lifetime = vec4(c_start_size, c_start_lifetime, c_start_lifetime, 0.0f);

    particles[index] = p;
}

#endif
