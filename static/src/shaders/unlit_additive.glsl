#version 300 es

precision highp float;

layout(std140) uniform Camera {
	mat4 VP;
	vec3 ambient;
	vec3 viewPos;
};

#ifdef VERT_SRC
layout (location = 0) in vec3 a_Position;
layout (location = 1) in vec3 a_Normal;
layout (location = 2) in vec2 a_UV1;

// Instance stuff
layout (location = 3) in vec4 a_InstanceColor;
layout (location = 4) in mat4 a_Transform;
// location = 5 aTransform
// location = 6 aTransform
// location = 7 aTransform
layout (location = 8) in mat4 a_NormalTransform;

out vec4 ColorV;
//out vec3 FragPos;
out vec3 Normal;
out vec2 TexCoords;

void main() {
	ColorV = a_InstanceColor;
	vec3 FragPos = vec3(a_Transform * vec4(a_Position, 1.0));
	Normal = normalize(vec3(a_NormalTransform * vec4(a_Normal, 0.0)));
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
	
	float alpha = tex.a * ColorV.a;
	vec3 color = tex.rgb * ColorV.rgb;

	if (alpha <= 0.001) {
		discard;
	}
	//tex *= ColorV;
	FragColor = vec4(color * alpha , alpha);
}
#endif
