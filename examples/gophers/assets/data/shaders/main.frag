#version 300 es
precision mediump float;

in vec4 ColorV;
in vec3 FragPos;
in vec3 Normal;
in vec2 TexCoords;
out vec4 FragColor;

// Just cause
uniform float time;

uniform vec3 ambient;

// Diffuse
uniform sampler2D albedo;
uniform bool albedoEnabled;
uniform vec2 albedoTile;

#define NLIGHTS 1
// lights
uniform vec3 lightPos[NLIGHTS];
uniform	vec3 lightColor[NLIGHTS];

// Material related
uniform float specularStrength;
uniform float shininess;


uniform vec3 viewPos;

void main() {
	vec4 inColor = ColorV;

	// Diffuse improve branching
	if (albedoEnabled) {
		vec2 tx = TexCoords.xy;
		vec2 phase = fract(tx.xy / albedoTile);
		inColor *= texture(albedo,phase);
	}

	vec3 lightDir = normalize(lightPos[0] - FragPos);
	vec3 viewDir = normalize(viewPos - FragPos);
	vec3 norm = normalize(Normal);

	// Specular
	vec3 specular = vec3(0);
	if (shininess > 0.0) {
		// Phong:
		//vec3 reflectDir = reflect(-lightDir, norm);
		//float spec = pow(max(dot(viewDir, reflectDir), 0.0), shininess);

		// Blinn-Phong:
		vec3 halfwayDir = normalize(lightDir + viewDir);
		float spec = pow(max(dot(norm, halfwayDir), 0.0), shininess/2.0);


		specular = specularStrength * spec * lightColor[0];
	}


	// Diffuse
	float diff = max(dot(norm, lightDir), 0.0);
	vec3 diffuse = diff * lightColor[0];
	vec3 result = (ambient + diffuse + specular) * inColor.rgb;


	// Our incoming color should be the texture
	// Premultiplied alpha

	FragColor = vec4(result.rgb * inColor.a, inColor.a); 
}

