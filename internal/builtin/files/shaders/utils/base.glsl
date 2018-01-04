#ifdef _VERTEX_
layout(location = 0) in vec3 vertex;
layout(location = 1) in vec3 normal;
layout(location = 2) in vec2 uv;

out vec2 vo_texture;

void main()
{
    vo_texture = uv;

    gl_Position = vec4(vertex, 1.0);
}

#endif

#ifdef _FRAGMENT_

subroutine vec4 RenderPassType();
subroutine uniform RenderPassType RenderPass;

in vec2 vo_texture;

layout(location = 0) out vec4 fo_destination;

layout(binding = 0) uniform sampler2D u_source;
layout(binding = 1) uniform sampler2D u_depth;

uniform vec2 u_resolution;

void main()
{
    fo_destination = RenderPass();
}

#endif
