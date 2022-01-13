#version 300 es
// TODO: eventually pass vertex skinning into the shader.
precision highp float;

#ifdef VERT_SRC
layout (location = 0) in vec3 a_Position;
layout (location = 2) in vec2 a_UV1;
layout (location = 3) in vec4 a_InstanceColor;
layout (location = 4) in mat4 a_Transform;

out vec4 ColorV;
out vec2 TexCoords;

uniform mat4 view; 

void main() {
	ColorV = a_InstanceColor;
	TexCoords = a_UV1;
	gl_Position = view * a_Transform * vec4(a_Position, 1.0);
}
#endif

#ifdef FRAG_SRC
in vec4 ColorV;
in vec2 TexCoords;

// Pass albedo to check alpha too
uniform sampler2D albedoMap;
uniform float u_AlphaCutoff;

void main() {
	float alpha = texture(albedoMap, TexCoords).a * ColorV.a;
	if (alpha <= u_AlphaCutoff) {
		discard;
	}
	// we only need to write depth which is already doing
}
#endif

