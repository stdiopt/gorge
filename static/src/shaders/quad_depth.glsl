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

uniform int perspective;
uniform float near_plane;
uniform float far_plane;

float LinearizeDepth(float depth) {
    float z = depth * 2.0 - 1.0; // Back to NDC 
    return (2.0 * near_plane * far_plane) / (far_plane + near_plane - z * (far_plane - near_plane));
}

void main() {
    float depthValue = texture(albedoMap, TexCoords).r;

	if (perspective == 1) {
		FragColor = vec4(vec3(LinearizeDepth(depthValue) / far_plane), 1.0); // perspective
		//FragColor = vec4(TexCoords.x, 0.0, 0.0, 1.0);
	} else {
		FragColor = vec4(vec3(depthValue)*4.0, 1.0); // orthographic
		//FragColor = vec4(0.0, TexCoords.x,0.0,1.0);
	}

}  
#endif
