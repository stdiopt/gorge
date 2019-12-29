#version 300 es
precision highp float;


in vec4 ColorV;
in vec3 FragPos;
in vec3 Normal;
in vec2 TexCoords;

out vec4 FragColor;

uniform vec3 viewPos;
uniform samplerCube envMap;

void main() {

	vec3 N = normalize(Normal);
	// {lpf} new
	// Reflections here
	vec3 I = normalize(FragPos - viewPos);
	vec3 R = reflect(I, N);
	//vec3 R = refract(I, normalize(Normal),1.0/2.42);
	vec3 ambient = texture(envMap, R).rgb;
	FragColor = vec4(ambient,1.0);
}
