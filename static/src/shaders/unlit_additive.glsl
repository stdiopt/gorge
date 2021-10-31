#version 300 es

precision mediump float;

layout(std140) uniform Camera {
	mat4 VP;
	vec3 ambient;
	vec3 viewPos;
};

#ifdef VERT_SRC
in vec3 a_Position;
in vec3 a_Normal;
in vec2 a_UV1;
in vec4 a_InstanceColor;
in mat4 a_Transform;
in mat4 a_NormalTransform;

out vec4 ColorV;
out vec3 FragPos;
out vec3 Normal;
out vec2 TexCoords;


void main() {
	ColorV = a_InstanceColor;
	FragPos = vec3(a_Transform * vec4(a_Position, 1.0));
	Normal = mat3(a_Transform) * a_Normal;
	TexCoords = a_UV1;

	gl_Position = VP * vec4(FragPos, 1.0);
}
#endif

#ifdef FRAG_SRC
out vec4 FragColor;

in vec4 ColorV;
in vec3 Normal;
in vec2 TexCoords;

uniform sampler2D albedoMap;

void main() {
	vec4 tex = texture(albedoMap, TexCoords);
	if (tex.a <= 0.01) {
		discard;
	}
	tex *= ColorV;
	FragColor = vec4(tex.rgb, tex.a);
}
#endif
