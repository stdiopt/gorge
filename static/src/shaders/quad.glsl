#version 300 es
precision highp float;

#ifdef VERT_SRC
out vec2 TexCoords;

void main() {
	float x = float(((uint(gl_VertexID) + 2u) / 3u)%2u);
	float y = float(((uint(gl_VertexID) + 1u) / 3u)%2u);

	gl_Position = vec4(-1.0f + x*2.0f, -1.0f+y*2.0f, 0.0f, 1.0f);
	TexCoords = vec2(x, y);
}
#endif

#ifdef FRAG_SRC
in vec2 TexCoords;

out vec4 FragColor;

uniform sampler2D albedoMap;

void main() {
	FragColor = vec4(texture(albedoMap, TexCoords).xyz,1);
}
#endif
