#version 300 es
layout (location = 0) in vec3 aPosition;

out vec3 TexCoords;

uniform mat4 projection;
uniform mat4 view;

void main()
{
	TexCoords = aPosition;
	mat4 lview = mat4(mat3(view));
	vec4 pos =  projection * lview * vec4(aPosition, 1.0);
	gl_Position = pos.xyww;
}  
