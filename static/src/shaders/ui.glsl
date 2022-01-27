#version 300 es
#ifdef GL_ES
// to have functions like fwidth for OpenGL ES or WebGL, the extension should
// be explicitly enabled.
#extension GL_OES_standard_derivatives : enable
#endif

// Based on unlit but with attempt to draw borders on shaders

precision highp float;

layout(std140) uniform Camera {
	mat4 VP;
	vec3 ambient;
	vec3 viewPos;
};

uniform vec4 rect;

#ifdef VERT_SRC
in vec3 a_Position;
in vec2 a_UV1;
layout (location = 3) in vec4 a_InstanceColor;
layout (location = 4) in mat4 a_Transform;
layout (location = 8) in mat4 a_NormalTransform;


// Experimental
out vec4 ColorV;
out vec3 FragPos;
out vec2 TexCoords;

#ifdef HAS_BORDER
out vec4 Border;

uniform vec4 border;
#endif

void main() {

#ifdef HAS_BORDER
	//Border = border;
	Border = vec4(
			border[0]/rect[2],
			border[1]/rect[3],
			border[2]/rect[2],
			border[3]/rect[3]
	);
#endif
	ColorV = a_InstanceColor;
	FragPos = vec3(a_Transform * vec4(a_Position, 1.0));
	TexCoords = a_UV1;
	gl_PointSize = 5.0;
	gl_Position = VP * vec4(FragPos, 1.0);
}
#endif

#ifdef FRAG_SRC

#define linearstep(edge0, edge1, x) clamp((x - (edge0)) / (edge1 - (edge0)), 0.0, 1.0)

out vec4 FragColor;

in vec4 ColorV;
in vec3 FragPos;
in vec2 TexCoords;
in vec4 Border;

uniform sampler2D albedoMap;

#ifdef HAS_BORDER
uniform vec4 borderColor;
#endif

// http://jeremt.github.io/pages/anti-aliased_shapes_in_glsl.html

void main() {
	vec4 color = texture(albedoMap, TexCoords);
	color *= ColorV;
#ifdef HAS_BORDER
	//vec2 uvPixel  = fwidth(TexCoords);
	float border = 0.0;
	//vec2 bb = step(Border.xy, TexCoords) * step(Border.zw, vec2(1) - TexCoords);

	//vec2 bb = linearstep(vec2(0), Border.xy, TexCoords) *
	//	  linearstep(vec2(0), Border.zw, vec2(1) - TexCoords);
	
	vec2 bb = smoothstep(vec2(0), Border.xy, TexCoords) *
		  smoothstep(vec2(0), Border.zw, vec2(1) - TexCoords);
	border = bb.x * bb.y;

	color = mix(borderColor, color, border);
#endif
	if (color.a <= 0.0) {
		discard;
	}
	FragColor = vec4(color.rgb * color.a, color.a);
}
#endif
