#version 300 es

// Based on unlit but with attempt to draw borders on shaders

precision highp float;

layout(std140) uniform Camera {
	mat4 VP;
	vec3 ambient;
	vec3 viewPos;
};

uniform vec4 rect;

#ifdef VERT_SRC
in vec3 a_Position;
in vec2 a_UV1;
layout (location = 3) in vec4 a_InstanceColor;
layout (location = 4) in mat4 a_Transform;
layout (location = 8) in mat4 a_NormalTransform;


// Experimental
out vec4 ColorV;
out vec3 FragPos;
out vec2 TexCoords;

#ifdef HAS_BORDER
out vec4 Border;
uniform vec4 border;
#endif

void main() {

#ifdef HAS_BORDER
	Border = vec4(
			border[0]/rect[2], 
			border[1]/rect[3], 
			1.0-border[2]/rect[2], 
			1.0-border[3]/rect[3]
	);
#endif
	ColorV = a_InstanceColor;
	FragPos = vec3(a_Transform * vec4(a_Position, 1.0));
	TexCoords = a_UV1;
	gl_PointSize = 5.0;
	gl_Position = VP * vec4(FragPos, 1.0);
}
#endif

#ifdef FRAG_SRC
out vec4 FragColor;

in vec4 ColorV;
in vec3 FragPos;
in vec2 TexCoords;
in vec4 Border;

uniform sampler2D albedoMap;

#ifdef HAS_BORDER
uniform vec4 borderColor;
#endif

void main() {
#ifdef HAS_BORDER
	if (TexCoords.x < Border[0] || TexCoords.x > Border[2]) {
		FragColor = vec4(borderColor.rgb * borderColor.a, borderColor.a);
		return;
	}
	if (TexCoords.y < Border[1] || TexCoords.y > Border[3]) {
		FragColor = vec4(borderColor.rgb * borderColor.a, borderColor.a);
		return;
	}
#endif
	vec4 tex = texture(albedoMap, TexCoords);
	if (tex.a <= 0.0) {
		discard;
	}
	tex *= ColorV;
	FragColor = vec4(tex.rgb * tex.a, tex.a);
}
#endif
