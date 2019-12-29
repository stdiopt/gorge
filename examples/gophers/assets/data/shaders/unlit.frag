#version 300 es

precision mediump float;

in vec4 ColorV;
in vec3 FragPos;
in vec3 Normal;
in vec2 TexCoords;
out vec4 FragColor;

uniform sampler2D albedoMap;

void main() {
	vec4 tex = texture(albedoMap,TexCoords);
	if (tex.a <= 0.0) {
		discard;
	}
	FragColor = vec4(tex.rgb * tex.a, tex.a);
}
