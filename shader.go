package gorge

// This just contains data for renderer

import (
	"fmt"
	"strings"

	"github.com/stdiopt/gorge/m32"
)

// ShaderResourcer is the interface to return a shader resource.
type ShaderResourcer interface {
	Resource() ShaderResource
}

// ShaderResource is a shader resource interface.
type ShaderResource interface {
	isShader()
	isGPU()
}

// ShaderData contains shaders sources
type ShaderData struct {
	GPU
	Name string
	Src  []byte
}

// Resource implements a resourcer.
func (d *ShaderData) Resource() ShaderResource { return d }

func (d *ShaderData) isShader() {}

func (d *ShaderData) String() string {
	return fmt.Sprintf("(shader name: %q)", d.Name)
}

// shaderProps something add definitions and uniform in shaders
// mostly to be used in material and mesh
type shaderProps struct {
	samplers    map[string]*Texture
	props       map[string]interface{}
	defines     map[string]string
	definesHash uint
	updates     int
}

func (s shaderProps) copy() shaderProps {
	r := shaderProps{
		updates: s.updates,
	}
	if s.props != nil {
		r.props = map[string]interface{}{}
		for k, v := range s.props {
			r.props[k] = v
		}

	}
	if s.defines != nil {
		r.defines = map[string]string{}
		for k, v := range s.defines {
			r.defines[k] = v
		}
	}
	return r
}

// ResetProps will reset uniforms and texture properties from material
func (s *shaderProps) ResetProps() {
	s.defines = nil
	s.props = nil
	s.definesHash = 0
}

func (s *shaderProps) SetTexture(name string, t Texturer) {
	if s.samplers == nil {
		s.samplers = map[string]*Texture{}
	}

	if t == nil {
		delete(s.samplers, name)
		return
	}
	s.samplers[name] = t.Texture()
}

// Set properties by name
func (s *shaderProps) Set(name string, v interface{}) {
	if t, ok := v.(Texturer); ok {
		s.SetTexture(name, t.Texture())
		return
	}
	if s.props == nil {
		s.props = map[string]interface{}{}
	}

	// Extra case
	if f, ok := v.(float64); ok {
		v = float32(f)
	}
	s.props[name] = v
}

func (s shaderProps) DefinesHash() uint {
	if s.definesHash == 0 {
		for _, d := range s.defines {
			s.definesHash ^= StringHash(d)
		}
	}
	return s.definesHash
}

func (s *shaderProps) Update() {
	s.updates++
}

func (s *shaderProps) Updates() int {
	return s.updates
}

// Define add defines
func (s *shaderProps) Define(defs ...string) {
	if s.defines == nil {
		s.defines = map[string]string{}
	}
	for _, d := range defs {
		p := strings.SplitN(d, " ", 2)
		val := ""
		if len(p) > 1 {
			val = p[1]
		}
		s.defines[p[0]] = val
	}
	s.updates++
}

// Undefine removes definitions
func (s *shaderProps) Undefine(defs ...string) {
	if s.defines == nil {
		return
	}
	for _, d := range defs {
		delete(s.defines, d)
	}
	s.updates++
}

// Defines declare some shader defines
func (s *shaderProps) Defines() map[string]string {
	return s.defines
}

// Get return named property
func (s *shaderProps) Get(name string) interface{} {
	if s.props == nil {
		return nil
	}
	return s.props[name]
}

func (s *shaderProps) GetTexture(name string) *Texture {
	if s.samplers == nil {
		return nil
	}
	return s.samplers[name]
}

// Props returns the properties of this material
func (s *shaderProps) Props() map[string]interface{} {
	return s.props
}

// SetFloat32 XXX testing sets a float32
func (s *shaderProps) SetFloat32(name string, v float32) {
	s.Set(name, v)
}

func (s *shaderProps) SetVec3(name string, v1, v2, v3 float32) {
	s.Set(name, m32.Vec3{v1, v2, v3})
}

func (s *shaderProps) SetVec4(name string, v1, v2, v3, v4 float32) {
	s.Set(name, m32.Vec4{v1, v2, v3, v4})
}
