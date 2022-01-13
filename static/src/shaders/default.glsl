#version 300 es

// it was 4, experiment huge amount of lights
//#define MAX_POINT_LIGHTS 8
//#define MAX_DIR_LIGHTS 8
#define LIGHT_COUNT 4
#define SAMPLERARRAY_HACKERY

precision highp float;
precision highp samplerCubeShadow;
precision highp sampler2DShadow;

layout(std140) uniform Camera {
	mat4 VP;
	vec3 ambient;
	vec3 viewPos;
};

#ifdef VERT_SRC
	layout (location = 0) in vec3 a_Position;
	layout (location = 1) in vec3 a_Normal;
	layout (location = 2) in vec2 a_UV1;

	// Instance stuff
	layout (location = 3) in vec4 a_InstanceColor;
	layout (location = 4) in mat4 a_Transform;
	layout (location = 8) in mat4 a_NormalTransform;
	// location = 5 aTransform
	// location = 6 aTransform
	// location = 7 aTransform

	out vec4 ColorV;
	out vec3 FragPos;
	out vec3 Normal;
	out vec2 TexCoords;

	void main() {
		ColorV = a_InstanceColor;
		FragPos = vec3(a_Transform * vec4(a_Position, 1.0));
		Normal = normalize(vec3(a_NormalTransform * vec4(a_Normal, 0.0)));
		TexCoords = a_UV1;

		gl_PointSize = 5.0;
		gl_Position = VP * vec4(FragPos,1.0);
	}
#endif

#ifdef FRAG_SRC
	out vec4 FragColor;

	in vec4 ColorV;
	in vec3 FragPos;
	in vec3 Normal;
	in vec2 TexCoords;

	// if envMap
	uniform bool hasEnvMap;
	uniform samplerCube envMap;

	// material parameters
	uniform sampler2D albedoMap;
	uniform sampler2D normalMap;
	uniform sampler2D metallicMap;
	uniform sampler2D roughnessMap;
	uniform sampler2D aoMap;

	uniform float metallic;
	uniform float roughness;
	uniform float ao;

	// Engine automagically sets these (wrongly)
	uniform bool has_normalMap;
	uniform bool has_metallicMap;
	uniform bool has_roughnessMap;
	uniform bool has_aoMap;

	// IBL
	uniform samplerCube irradianceMap;
	uniform samplerCube prefilterMap;
	uniform sampler2D brdfLUT;

	uniform bool has_irradianceMap;
	uniform bool has_prefilterMap;

	const int LightType_Directional = 0;
	const int LightType_Point = 1;
	const int LightType_Spot = 2;

	struct Light {
		vec3 position;
		vec3 direction;

		float range;

		vec3 color;
		float intensity;

		float innerConeCos;

		float outerConeCos;
		int type;

		mat4 matrix;
		int depthIndex;
	};


	layout(std140) uniform Lights {
		int u_nLights;
		Light u_Lights[LIGHT_COUNT];
	};
	
#ifdef SAMPLERARRAY_HACKERY
	uniform samplerCube depthCube[LIGHT_COUNT];
	uniform sampler2D depth2D[LIGHT_COUNT];
#else
	struct DepthCube { samplerCube depthMap; };
	struct Depth2D { sampler2D depthMap; };
	uniform DepthCube depthCube[LIGHT_COUNT];
	uniform Depth2D depth2D[LIGHT_COUNT];
#endif
	

	//uniform samplerCube depthCube[LIGHT_COUNT];

	// https://github.com/KhronosGroup/glTF/blob/master/extensions/2.0/Khronos/KHR_lights_punctual/README.md#range-property
	float getRangeAttenuation(float range, float distance) {
		if (range <= 0.0) {
			// negative range means unlimited
			return 1.0;
		}
		return max(min(1.0 - pow(distance / range, 4.0), 1.0), 0.0) / pow(distance, 2.0);
	}

	// https://github.com/KhronosGroup/glTF/blob/master/extensions/2.0/Khronos/KHR_lights_punctual/README.md#inner-and-outer-cone-angles
	float getSpotAttenuation(vec3 pointToLight, vec3 spotDirection, float outerConeCos, float innerConeCos) {
		float actualCos = dot(normalize(spotDirection), normalize(-pointToLight));
		if (actualCos > outerConeCos) {
			if (actualCos < innerConeCos) {
				return smoothstep(outerConeCos, innerConeCos, actualCos);
			}
			return 1.0;
		}
		return 0.0;
	}
	const float PI = 3.14159265359;
	const float minFloat = 0.0001;
	
	// ----------------------------------------------------------------------------
	// Easy trick to get tangent-normals to world-space to keep PBR code simplified.
	// Don't worry if you don't get what's going on; you generally want to do normal
	// mapping the usual way for performance anways; I do plan make a note of this
	// technique somewhere later in the normal mapping tutorial.
	vec3 getNormalFromMap() {
		vec3 tangentNormal = texture(normalMap, TexCoords).xyz * 2.0 - 1.0;

		vec3 Q1  = dFdx(FragPos);
		vec3 Q2  = dFdy(FragPos);
		vec2 st1 = dFdx(TexCoords);
		vec2 st2 = dFdy(TexCoords);

		vec3 N   = normalize(Normal);
		vec3 T  = normalize(Q1*st2.t - Q2*st1.t);
		vec3 B  = -normalize(cross(N, T));
		mat3 TBN = mat3(T, B, N);

		return normalize(TBN * tangentNormal);
	}

	// ----------------------------------------------------------------------------
	float DistributionGGX(vec3 N, vec3 H, float roughness) {
		float a = roughness*roughness;
		float a2 = a*a;
		float NdotH = max(dot(N, H), 0.0);
		float NdotH2 = NdotH*NdotH;

		float nom   = a2;
		float denom = (NdotH2 * (a2 - 1.0) + 1.0);
		denom = PI * denom * denom;

		return nom / max(denom, minFloat); // prevent divide by zero for roughness=0.0 and NdotH=1.0
	}
	// ----------------------------------------------------------------------------
	float GeometrySchlickGGX(float NdotV, float roughness) {
		float r = (roughness + 1.0);
		float k = (r*r) / 8.0;

		float nom   = NdotV;
		float denom = NdotV * (1.0 - k) + k;

		return nom / denom;
	}
	// ----------------------------------------------------------------------------
	float GeometrySmith(vec3 N, vec3 V, vec3 L, float roughness) {
		float NdotV = max(dot(N, V), 0.0);
		float NdotL = max(dot(N, L), 0.0);
		float ggx2 = GeometrySchlickGGX(NdotV, roughness);
		float ggx1 = GeometrySchlickGGX(NdotL, roughness);

		return ggx1 * ggx2;
	}
	// ----------------------------------------------------------------------------
	vec3 fresnelSchlick(float cosTheta, vec3 F0) {
		return F0 + (1.0 - F0) * pow(1.0 - cosTheta, 5.0);
	}

	vec3 fresnelSchlickRoughness(float cosTheta, vec3 F0, float roughness) {
		return F0 + (max(vec3(1.0 - roughness), F0) - F0) * pow(1.0 - cosTheta, 5.0);
	}

	// ----------------------------------------------------------------------------
	const vec3 gridSamplingDisk[20] = vec3[] (
	   vec3(1, 1,  1), vec3( 1, -1,  1), vec3(-1, -1,  1), vec3(-1, 1,  1),
	   vec3(1, 1, -1), vec3( 1, -1, -1), vec3(-1, -1, -1), vec3(-1, 1, -1),
	   vec3(1, 1,  0), vec3( 1, -1,  0), vec3(-1, -1,  0), vec3(-1, 1,  0),
	   vec3(1, 0,  1), vec3(-1,  0,  1), vec3( 1,  0, -1), vec3(-1, 0, -1),
	   vec3(0, 1,  1), vec3( 0, -1,  1), vec3( 0, -1, -1), vec3( 0, 1, -1)
	);

	float getDepth2D(int n, vec2 coord) {
#ifdef SAMPLERARRAY_HACKERY
		switch (n) {
			case 0: return texture(depth2D[0], coord).r;
			case 1: return texture(depth2D[1], coord).r;
			case 2: return texture(depth2D[2], coord).r;
			case 3: return texture(depth2D[3], coord).r;
			/*
			case 4: return texture(depth2D[4],coord).r;
			case 5: return texture(depth2D[5],coord).r;
			case 6: return texture(depth2D[6],coord).r;
			case 7: return texture(depth2D[7],coord).r;
			*/
		}
#else
		return texture(depth2D[n].depthMap, coord).r;
#endif
	}
	float getDepthCube(int n, vec3 coord) {
#ifdef SAMPLERARRAY_HACKERY
		switch (n) {
			case 0: return texture(depthCube[0], coord).r;
			case 1: return texture(depthCube[1], coord).r;
			case 2: return texture(depthCube[2], coord).r;
			case 3: return texture(depthCube[3], coord).r;
			/*
			case 4: return texture(depthCube[4],coord).r;
			case 5: return texture(depthCube[5],coord).r;
			case 6: return texture(depthCube[6],coord).r;
			case 7: return texture(depthCube[7],coord).r;
			*/
		}
#else
		return texture(depthCube[n].depthMap, coord).r;
#endif
	}

	float shadowCubeCalculation(Light light, vec3 fragPos, vec3 pointToLight, float distance) {
		int di = light.depthIndex;
		// This should be range?
		// We already have this else where
		float shadow = 0.0;

		vec3 fragToLight = -pointToLight;
		float bias = 0.25;
		int samples = 20;
		float viewDistance = length(viewPos - fragPos);
		float diskRadius = (1.0 + (viewDistance / light.range)) / 25.0;
		
		for(int i = 0; i < samples; ++i) {
			float closestDepth = getDepthCube(di, fragToLight + gridSamplingDisk[i] * diskRadius);
			closestDepth *= light.range;   // undo mapping [0;1]
			if(distance - bias > closestDepth)
				shadow += 1.0;
		}
		
		shadow /= float(samples);
		return shadow;
		
	}

	float shadow2DCalculation(Light light, vec3 fragPos) {
		int di = light.depthIndex;
		
		float shadow = 0.0;
		
		// This could be done at vertex shader
		vec4 fragPosLightSpace = light.matrix * vec4(FragPos, 1.0);
		vec3 projCoords = fragPosLightSpace.xyz / fragPosLightSpace.w;
		projCoords = projCoords * 0.5 + 0.5;
		if (projCoords.z > 1.0)
			return 0.0;
		//float closestDepth = 1.-texture(depth[ti].DepthMap2D, projCoords.xy).r;
		float currentDepth = projCoords.z;
		float bias = 0.005;
		//vec2 texelSize = vec2(1.0f) / vec2(textureSize(depth2D[di].depthMap, 0));
		vec2 texelSize = vec2(1.0f) / vec2(2048,2048);
		for(int x = -1; x <= 1; ++x) {
			for(int y = -1; y <= 1; ++y) {
				float pcfDepth = getDepth2D(di, projCoords.xy + vec2(x, y) * texelSize);
				shadow += currentDepth - bias > pcfDepth ? 1.0 : 0.0;
			}
		}
		shadow /= 9.0;

		return shadow;
	}

	float ShadowCalculation(Light light, vec3 fragPos, vec3 pointToLight, float distance) {
		if (light.depthIndex < 0) {
			return 0.0;
		}

		if (light.type == LightType_Point) {
			return shadowCubeCalculation(light, fragPos, pointToLight, distance);
		}
		return shadow2DCalculation(light, fragPos);
	}

	// ----------------------------------------------------------------------------
	void main() {
		//vec3 ambient = Ambient;
		float alpha = texture(albedoMap, TexCoords).a * ColorV.a;

		if (alpha <= 0.001) {
			discard;
		}

		vec3 albedo = ColorV.rgb * pow(texture(albedoMap, TexCoords).rgb, vec3(2.2));

		// material properties
		float metallic = metallic;
		if (has_metallicMap) {
			metallic = texture(metallicMap, TexCoords).r;
		}
		float roughness = roughness;
		if (has_roughnessMap) {
			roughness = texture(roughnessMap, TexCoords).r;
		}
		float ao = ao;
		if (has_aoMap) {
			ao = texture(aoMap, TexCoords).r;
		}
		vec3 N = normalize(Normal);
		if (has_normalMap) {
			N = getNormalFromMap();
		}
		
		vec3 V = normalize(viewPos - FragPos);
		vec3 R = reflect(-V, N);

		vec3 F0 = vec3(0.04);
		F0 = mix(F0, albedo, metallic);

		vec3 Lo = vec3(0.0);

		int cnlights = min(u_nLights, LIGHT_COUNT);
		for(int i = 0; i < cnlights; ++i) {
			Light light = u_Lights[i];


			vec3 pointToLight = -light.direction;
			float rangeAttenuation = 1.0;
			float spotAttenuation = 1.0;

			if (light.type != LightType_Directional) {
				pointToLight = light.position - FragPos;
			}
			float distance = length(pointToLight);
			if (light.type != LightType_Directional) {
				rangeAttenuation = getRangeAttenuation(light.range, distance);
			}
			if (light.type == LightType_Spot) {
				spotAttenuation = getSpotAttenuation(pointToLight, light.direction, light.outerConeCos, light.innerConeCos);
			}
			vec3 radiance = rangeAttenuation * spotAttenuation * light.intensity * light.color;

			vec3 L = normalize(pointToLight);
			vec3 H = normalize(V + L);

			// Cook-Torrance BRDF
			float NDF = DistributionGGX(N, H, roughness);
			float G   = GeometrySmith(N, V, L, roughness);
			vec3 F    = fresnelSchlick(clamp(dot(H, V), 0.0, 1.0), F0);
			
			vec3 nominator		= NDF * G * F;
			float denominator	= float(cnlights) * max(dot(N, V), 0.0) * max(dot(N, L), 0.0) + minFloat;
			vec3 specular		= nominator / denominator;

			vec3 kS = F;
			vec3 kD = vec3(1.0) - kS;
			kD *= 1.0 - metallic;
			
			float NdotL = max(dot(N, L), 0.0);

			float shadow = ShadowCalculation(light,FragPos,pointToLight,distance);

			// note that we already multiplied the BRDF by the Fresnel (kS) so we won't multiply by kS again
			Lo += (kD * albedo / PI + specular) * radiance * NdotL * (1.0 - shadow);
		}
		
		vec3 F = fresnelSchlickRoughness(max(dot(N, V), 0.0), F0, roughness);

		vec3 kS = F;
		vec3 kD = 1.0 - kS;
		kD *= 1.0 - metallic;

		vec3 irradiance = ambient;
		if (has_irradianceMap) {
			irradiance = texture(irradianceMap, N).rgb;
		}
		vec3 diffuse      = irradiance * albedo;

		vec3 prefilteredColor = vec3(0,0,0);
		if (has_prefilterMap) {
			const float MAX_REFLECTION_LOD = 4.0;
			prefilteredColor = textureLod(prefilterMap, R,  roughness * MAX_REFLECTION_LOD).rgb;
		} else {
			// Dumb way
			//prefilteredColor = texture(envMap, R).rgb * (1.0 - roughness);
		}

		// sample both the pre-filter map and the BRDF lut and combine them together as per the Split-Sum approximation to get the IBL specular part.
		vec2 brdf  = texture(brdfLUT, vec2(max(dot(N, V), 0.0), roughness)).rg;
		vec3 specular = prefilteredColor * (F * brdf.x + brdf.y);

		vec3 ambient = (kD * diffuse + specular) * ao;

		vec3 color = ambient + Lo;

		// HDR tonemapping
		color = color / (color + vec3(1.0));
		// gamma correct
		color = pow(color, vec3(1.0/2.2));

		//FragColor = vec4(color , 1.0);
		FragColor = vec4(color * alpha, alpha);
		

		// {lpf} old
		// Reflections here
		/*vec3 amb = ambient;
		if (hasEnvMap) {
			vec3 I = normalize(FragPos - viewPos);
			vec3 R = reflect(I, N);
			//vec3 R = refract(I, N, 1.0/1.32);
			//vec3 R = refract(I, normalize(Normal),1.0/2.42);
			amb = texture(envMap, R).rgb;
		}
		vec3 color = vec3(0.3) * amb * albedo * ao;
		color = color + Lo;
		color = color / (color + vec3(1.0));
		color = pow(color, vec3(1.0/2.2)); 
		FragColor = vec4(color * alpha, alpha);*/
}
#endif
