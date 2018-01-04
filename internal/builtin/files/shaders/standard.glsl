#ifdef _VERTEX_
layout(location = 0) in vec3 vertex;
layout(location = 1) in vec3 normal;
layout(location = 2) in vec2 uv;

out vec3 vo_position;
out vec3 vo_normal;
out vec3 vo_eye;
out vec3 vo_ws_position;
out vec3 vo_ws_normal;
out vec2 vo_texture;

uniform mat4 v_mvp_matrix;
uniform mat4 v_projection_matrix;
uniform mat4 v_view_matrix;
uniform mat4 v_model_matrix;
uniform mat3 v_normal_matrix;

void main()
{
    vo_texture = uv;
    vo_normal = normal;// normalize(v_normal_matrix * normal);
    vo_position = vertex;
    vo_ws_position = vec3(v_model_matrix * vec4(vertex, 1.0));
    vo_ws_normal = vec3(v_model_matrix * vec4(normal, 1.0));

    gl_Position = v_projection_matrix * v_view_matrix * v_model_matrix * vec4(vertex, 1.0);
}

#endif

#ifdef _FRAGMENT_
subroutine void RenderPassType();
subroutine uniform RenderPassType RenderPass;

in vec3 vo_position;
in vec3 vo_normal;
in vec3 vo_eye;
in vec3 vo_ws_position;
in vec3 vo_ws_normal;
in vec2 vo_texture;

layout(location = 0) out vec4 fo_attachment0;
layout(location = 1) out uvec4 fo_attachment1;

layout(binding = 0) uniform sampler2D f_attachment0;
layout(binding = 1) uniform usampler2D f_attachment1;
layout(binding = 2) uniform sampler2D f_depth;
layout(binding = 3) uniform samplerCube f_environment;
layout(binding = 4) uniform samplerCube f_irradiance;
layout(binding = 5) uniform sampler2D f_albedo_map;
layout(binding = 6) uniform sampler2D f_metallic_map;
layout(binding = 7) uniform sampler2D f_normal_map;

uniform vec3 f_camera;
uniform vec3 f_albedo;
uniform float f_roughness;
uniform float f_metallic;

#define PI   3.1415926535897932384626433832795
#define PI2  6.2831853071795864769252867665590

vec3 get_position(vec4 data)
{
    return data.xyz;
}

vec3 get_normal(uvec4 data)
{
    vec3 normal;

    normal.xy = unpackHalf2x16(data.x);
    normal.z = unpackHalf2x16(data.y).x;

    return normal;
}

vec3 get_albedo(uvec4 data)
{
    return unpackUnorm4x8(data.z).rgb;
}

vec3 get_reflection(vec3 dir)
{
    return texture(f_environment, dir).rgb;
}

float get_metallic(uvec4 data)
{
    return unpackHalf2x16(data.w).y;
}

float get_roughness(uvec4 data)
{
    return unpackHalf2x16(data.w).x;
}

float DistributionGGX(vec3 N, vec3 H, float a)
{
    float a2     = a*a;
    float NdotH  = max(dot(N, H), 0.0);
    float NdotH2 = NdotH*NdotH;

    float nom    = a2;
    float denom  = (NdotH2 * (a2 - 1.0) + 1.0);
    denom        = PI * denom * denom;

    return nom / denom;
}

float GeometrySchlickGGX(float NdotV, float k)
{
    float nom   = NdotV;
    float denom = NdotV * (1.0 - k) + k;

    return nom / denom;
}

float GeometrySmith(vec3 N, vec3 V, vec3 L, float k)
{
    float NdotV = max(dot(N, V), 0.0);
    float NdotL = max(dot(N, L), 0.0);
    float ggx1 = GeometrySchlickGGX(NdotV, k);
    float ggx2 = GeometrySchlickGGX(NdotL, k);

    return ggx1 * ggx2;
}

vec3 fresnelSchlick(float cosTheta, vec3 F0)
{
    return F0 + (1.0 - F0) * pow(1.0 - cosTheta, 5.0);
}

subroutine(RenderPassType)
void forward_pass()
{
    vec3 N = vo_ws_normal;

    fo_attachment0 = vec4(N, 1.0);
}

subroutine(RenderPassType)
void deferred_pass_geometry()
{
    fo_attachment0.xyz = vo_ws_position;

    fo_attachment1.x = packHalf2x16(vo_normal.xy);
    fo_attachment1.y = packHalf2x16(vec2(vo_normal.z, 0.0));
    fo_attachment1.z = packUnorm4x8(vec4(f_albedo, 1.0));
    fo_attachment1.w = packHalf2x16(vec2(f_roughness, f_metallic));
}

subroutine(RenderPassType)
void deferred_pass_ambient()
{
     float depth = texture(f_depth, vo_texture).r;
     if (depth == 1.0)
         discard;

    vec4 data0 = texture(f_attachment0, vo_texture);
    uvec4 data1 = texture(f_attachment1, vo_texture);

    vec3 albedo = get_albedo(data1);
    vec3 P = get_position(data0);
    vec3 V = normalize(f_camera - P);
    vec3 N = get_normal(data1);
    vec3 L = normalize(-reflect(V, N));

    vec3 irradiance = texture(f_irradiance, L).rgb;

    fo_attachment0 = vec4(irradiance, 1.0);
}

void main()
{
    RenderPass();
}

#endif
