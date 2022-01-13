#version 300 es
// TODO: eventually pass vertex skinning into the shader.
precision highp float;
#ifdef VERT_SRC
layout (location = 0) in vec3 a_Position;
layout (location = 2) in vec2 a_UV1;
layout (location = 3) in vec4 a_InstanceColor;
layout (location = 4) in mat4 a_Transform; // Aka model

out vec3 FragPos;
out vec4 ColorV;
out vec2 TexCoords;

uniform mat4 view;

void main() {
	ColorV = a_InstanceColor;
	TexCoords = a_UV1;
	FragPos = vec3(a_Transform * vec4(a_Position, 1.0));
	gl_Position = view * a_Transform * vec4(a_Position, 1.0);
}
#endif

#ifdef FRAG_SRC
precision mediump float;

in vec4 ColorV;
in vec3 FragPos;
in vec2 TexCoords;

// Pass albedo to check alpha too
uniform sampler2D albedoMap;
uniform float u_AlphaCutoff;

uniform vec3 lightPos;
uniform float farPlane;

void main() {

	float alpha = texture(albedoMap, TexCoords).a * ColorV.a;
	if (alpha <= u_AlphaCutoff) {
		discard;
	}

	// get distance between fragment and light source
	float lightDistance = length(FragPos.xyz - lightPos);
	// map to [0;1] range by dividing by far_plane
	lightDistance = lightDistance / farPlane;

	// write this as modified depth
	gl_FragDepth = lightDistance;
}
#endif
