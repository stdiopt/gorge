#version 300 es
precision highp float;

#ifdef VERT_SRC
	layout (location = 0) in vec3 a_Position;

	out vec3 TexCoords;

	uniform mat4 VP;

	void main() {
		TexCoords = a_Position;
		vec4 pos =  VP * vec4(a_Position, 1.0);
		gl_Position = pos;
	}
#endif

#ifdef FRAG_SRC
	out vec4 FragColor;

	in vec3 TexCoords;

	uniform samplerCube skybox;
	uniform float lod;

	void main() {
		FragColor = textureLod(skybox, TexCoords, lod);
	}
#endif

