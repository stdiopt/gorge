// Copyright 2019 Luis Figueiredo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package renderer

import (
	"errors"
	"fmt"

	"github.com/stdiopt/gorge/asset"
	"github.com/stdiopt/gorge/gl"
)

// uniform unfortunately webgl uses its own type for gl.Uniform rather than
// uint32 so we track index and loc
type uniform struct {
	loc   gl.Uniform
	value interface{}
}

// ShaderManager globalized shader states
type shaderManager struct {
	g       gl.Context3
	assets  *asset.System
	shaders map[string]*Shader
}

// Get it will return a cached program or load a shader using asset package
// Shaders should be named {name}.vert and {name}.frag
func (sm *shaderManager) Get(name string) (*Shader, error) {
	if sm.shaders == nil {
		sm.shaders = map[string]*Shader{}
	}
	// Already loaded
	if shader, ok := sm.shaders[name]; ok {
		return shader, nil
	}
	var vertSrc string
	var fragSrc string
	if name == "pbr" {
		vertSrc = shaderPBRVert
		fragSrc = shaderPBRFrag
	} else {
		var err error
		vertSrc, err = sm.assets.LoadString("shaders/" + name + ".vert")
		if err != nil {
			return nil, err
		}

		fragSrc, err = sm.assets.LoadString("shaders/" + name + ".frag")
		if err != nil {
			return nil, err
		}
	}

	shader, err := sm.CompileShader(vertSrc, fragSrc)
	if err != nil {
		return nil, err
	}
	sm.shaders[name] = shader
	return shader, nil
}

func (sm *shaderManager) New(p gl.Program) *Shader {
	g := sm.g
	// Load names from program
	nUniforms := g.GetProgrami(p, gl.ACTIVE_UNIFORMS)
	samplers := []string{}
	uniforms := map[string]*uniform{}
	nAttribs := g.GetProgrami(p, gl.ACTIVE_ATTRIBUTES)
	attribs := map[string]gl.Attrib{}

	for i := 0; i < nUniforms; i++ {
		name, _, ty := g.GetActiveUniform(p, uint32(i))
		switch ty {
		case gl.SAMPLER_2D, gl.SAMPLER_3D, gl.SAMPLER_CUBE:
			//	value = uint32(0)
			samplers = append(samplers, name)
		}
		loc := g.GetUniformLocation(p, name)
		uniforms[name] = &uniform{loc: loc}
	}
	for i := 0; i < nAttribs; i++ {
		name, _, _ := g.GetActiveAttrib(p, uint32(i))
		loc := g.GetAttribLocation(p, name)
		attribs[name] = loc
	}

	return &Shader{
		g:        sm.g,
		program:  p,
		attribs:  attribs,
		uniforms: uniforms,
		samplers: samplers,
	}
}

func (sm *shaderManager) CompileShader(vertSrc, fragSrc string) (*Shader, error) {
	program, err := sm.compileProgram(vertSrc, fragSrc)
	if err != nil {
		return nil, err
	}
	s := sm.New(program)
	return s, nil
}
func (sm *shaderManager) compileProgram(vertSrc, fragSrc string) (gl.Program, error) {
	g := sm.g
	var program gl.Program
	vertShader := g.CreateShader(gl.VERTEX_SHADER)

	g.ShaderSource(vertShader, vertSrc)
	g.CompileShader(vertShader)
	if g.GetShaderi(vertShader, gl.COMPILE_STATUS) == gl.FALSE {
		return program, fmt.Errorf("vertex: %s", g.GetShaderInfoLog(vertShader))
	}

	fragShader := g.CreateShader(gl.FRAGMENT_SHADER)
	g.ShaderSource(fragShader, fragSrc)
	g.CompileShader(fragShader)

	if g.GetShaderi(fragShader, gl.COMPILE_STATUS) == gl.FALSE {
		return program, fmt.Errorf("fragment: %s", g.GetShaderInfoLog(fragShader))
	}

	// Based on material
	program = g.CreateProgram()
	g.AttachShader(program, vertShader)
	g.AttachShader(program, fragShader)
	g.LinkProgram(program)
	if g.GetProgrami(program, gl.LINK_STATUS) == gl.FALSE {
		return program, errors.New(g.GetProgramInfoLog(program))
	}

	return program, nil
}

///////////////////////////////////////////////////////////////////////////////
// SHADER
///////////////////////////////////////////////////////////////////////////////

// Shader handle gl program, it holds a local State
type Shader struct {
	g       gl.Context3
	program gl.Program

	//UBO gl.Buffer

	// Extra attribs and uniforms
	attribs  map[string]gl.Attrib
	uniforms map[string]*uniform
	samplers []string
}

// Attrib returns the attribute for name
func (s *Shader) Attrib(k string) (gl.Attrib, bool) {
	a, ok := s.attribs[k]
	return a, ok
}

// Set sets a uniform value
// TODO: Try UBO's uniform buffer objects here too
func (s *Shader) Set(k string, v interface{}) {
	u, ok := s.uniforms[k]
	if !ok {
		return
	}
	if u.value == v {
		return
	}
	s.set(u, v)
}

func (s *Shader) set(u *uniform, v interface{}) {
	g := s.g
	// Lastly do a type switch
	switch v := v.(type) {
	case uint32:
		g.Uniform1i(u.loc, int(v))
	case int:
		g.Uniform1i(u.loc, v)
	case float32:
		g.Uniform1f(u.loc, v)
	case vec2:
		g.Uniform2fv(u.loc, v[:])
	case vec3:
		g.Uniform3fv(u.loc, v[:])
	case vec4:
		g.Uniform4fv(u.loc, v[:])
	case mat4:
		g.UniformMatrix4fv(u.loc, v[:])
	default:
		// Clear uniform
		panic(fmt.Sprintf("not implemented: %T", v))
	}
	u.value = v
}

func (s *Shader) bind() {
	s.g.UseProgram(s.program)
}

// TODO: make the default shader more simple
const (
	shaderPBRFrag = `#version 300 es
precision mediump float;

in vec4 ColorV;
in vec3 FragPos;
in vec3 Normal;
in vec2 TexCoords;
out vec4 FragColor;

uniform samplerCube envMap;

// material parameters
uniform sampler2D albedoMap;
// AlbedoTile

//uniform vec3 albedo;
uniform float metallic;
uniform float roughness;
uniform float ao;

#define NLIGHTS 1
// lights
uniform vec3 lightPos[NLIGHTS];
uniform vec3 lightColors[NLIGHTS];

uniform vec3 viewPos;

const float PI = 3.14159265359;
// ----------------------------------------------------------------------------
float DistributionGGX(vec3 N, vec3 H, float roughness) {
    float a = roughness*roughness;
    float a2 = a*a;
    float NdotH = max(dot(N, H), 0.0);
    float NdotH2 = NdotH*NdotH;

    float nom   = a2;
    float denom = (NdotH2 * (a2 - 1.0) + 1.0);
    denom = PI * denom * denom;

    return nom / max(denom, 0.001); // prevent divide by zero for roughness=0.0 and NdotH=1.0
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
// ----------------------------------------------------------------------------
void main() {		
	float alpha = texture(albedoMap,TexCoords).a * ColorV.a;

	if (alpha <= 0.0) {
		discard;
	}

	vec3 albedo = ColorV.rgb * pow(texture(albedoMap, TexCoords).rgb, vec3(2.2));

    vec3 N = normalize(Normal);
    vec3 V = normalize(viewPos - FragPos);

    // calculate reflectance at normal incidence; if dia-electric (like plastic) use F0 
    // of 0.04 and if it's a metal, use the albedo color as F0 (metallic workflow)    
    vec3 F0 = vec3(0.04); 
    F0 = mix(F0, albedo, metallic);

    // reflectance equation
    vec3 Lo = vec3(0.0);
    for(int i = 0; i < NLIGHTS; ++i) 
    {
        // calculate per-light radiance
        vec3 L = normalize(lightPos[i] - FragPos);
        vec3 H = normalize(V + L);
        float distance = length(lightPos[i] - FragPos);
        float attenuation = 1.0 / (distance * distance);
        vec3 radiance = lightColors[i] * attenuation;

        // Cook-Torrance BRDF
        float NDF = DistributionGGX(N, H, roughness);   
        float G   = GeometrySmith(N, V, L, roughness);      
        vec3 F    = fresnelSchlick(clamp(dot(H, V), 0.0, 1.0), F0);

        vec3 nominator    = NDF * G * F; 
        float denominator = 4.0 * max(dot(N, V), 0.0) * max(dot(N, L), 0.0);
        vec3 specular = nominator / max(denominator, 0.001); // prevent divide by zero for NdotV=0.0 or NdotL=0.0

        // kS is equal to Fresnel
        vec3 kS = F;
        // for energy conservation, the diffuse and specular light can't
        // be above 1.0 (unless the surface emits light); to preserve this
        // relationship the diffuse component (kD) should equal 1.0 - kS.
        vec3 kD = vec3(1.0) - kS;
        // multiply kD by the inverse metalness such that only non-metals 
        // have diffuse lighting, or a linear blend if partly metal (pure metals
        // have no diffuse light).
        kD *= 1.0 - metallic;	  

        // scale light by NdotL
        float NdotL = max(dot(N, L), 0.0);        

        // add to outgoing radiance Lo
        Lo += (kD * albedo / PI + specular) * radiance * NdotL;  // note that we already multiplied the BRDF by the Fresnel (kS) so we won't multiply by kS again
    }   

	// {lpf} new
	// Reflections here
	//vec3 I = normalize(FragPos - viewPos);
	//vec3 R = reflect(I, normalize(Normal));
	//vec3 R = refract(I, normalize(Normal),1.0/2.42);
	//vec3 ambient = texture(envMap, N).rgb;
	//FragColor = vec4(ambient,1);
	//
    // ambient lighting (note that the next IBL tutorial will replace 
    // this ambient lighting with environment lighting).
	
    vec3 ambient = vec3(0.03) * albedo * ao;
	


    vec3 color = ambient + Lo;

    // HDR tonemapping
    color = color / (color + vec3(1.0));
    // gamma correct
    color = pow(color, vec3(1.0/2.2)); 

    FragColor = vec4(color * alpha, alpha);
	/**/
}`
	shaderPBRVert = `#version 300 es
layout (location = 0) in vec3 aPosition;
layout (location = 1) in vec3 aNormal;
layout (location = 2) in vec2 aTexCoords;

layout (location = 3) in vec4 aColor;
layout (location = 4) in mat4 aTransform;

out vec4 ColorV;
out vec3 FragPos;
out vec3 Normal;
out vec2 TexCoords;

uniform mat4 projection;
uniform mat4 view;

void main() {
	ColorV = aColor;

    FragPos = vec3(aTransform * vec4(aPosition, 1.0));
    Normal = mat3(aTransform) * aNormal;   
    TexCoords = aTexCoords;

    gl_Position = projection * view * vec4(FragPos,1.0);
}
`
)
