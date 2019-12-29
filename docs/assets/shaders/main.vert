#version 300 es
layout (location = 0) in vec3 aPosition;
layout (location = 1) in vec3 aNormal; 
layout (location = 2) in vec2 aTexCoords;

// instanced
layout (location = 3) in vec4 aColor;
layout (location = 4) in mat4 aTransform;

out vec4 ColorV;
out vec3 FragPos;
out vec3 Normal;
out vec2 TexCoords;

uniform mat4 projection;


void main() {
	ColorV = aColor; 

	// This should be the model, and not the transform which contains the view?
	FragPos = vec3(aTransform * vec4(aPosition, 1.0));
	Normal = mat3(aTransform) * aNormal;
	TexCoords = aTexCoords;
	
	gl_Position = projection * vec4(FragPos, 1.0);
}

