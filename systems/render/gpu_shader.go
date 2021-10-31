package render

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/static"
	"github.com/stdiopt/gorge/systems/render/gl"
)

// REVAMPING shaders to use Defines

const (
	aPosition        = "a_Position"
	aNormal          = "a_Normal"
	aUV1             = "a_UV1"
	aInstanceColor   = "a_InstanceColor"
	aTransform       = "a_Transform"
	aNormalTransform = "a_NormalTransform"
)

// Spec Attrib positions
var (
	defAttrib = map[string]gl.Attrib{
		aPosition: gl.Attrib(0),
		aNormal:   gl.Attrib(1),
		aUV1:      gl.Attrib(2),

		aInstanceColor:   gl.Attrib(3),
		aTransform:       gl.Attrib(4),
		aNormalTransform: gl.Attrib(8),
	}
)

func shaderAttrib(s *Shader, k string) (gl.Attrib, bool) {
	var loc gl.Attrib
	var ok bool
	if s != nil {
		loc, ok = s.attribs[k]
	} else {
		loc, ok = defAttrib[k]
	}
	return loc, ok
}

type shaderManager struct {
	gorge *gorge.Context
	vbos  *vboManager
	// def    *Shader
	defRaw *gorge.ShaderData

	// Shader with variants
	hashedShaders map[*gorge.ShaderData]map[uint]*Shader

	count int
}

func newShaderManager(g *gorge.Context, vbos *vboManager) *shaderManager {
	m := &shaderManager{
		gorge:         g,
		vbos:          vbos,
		hashedShaders: map[*gorge.ShaderData]map[uint]*Shader{},
		defRaw:        static.Shaders.Default,
	}
	m.New(m.defRaw) // Bind a program to shader

	return m
}

func (m *shaderManager) New(r gpuResource) *Shader {
	s := &Shader{
		manager:  m,
		ubos:     map[string]struct{}{},
		uniforms: map[string]*uniform{},
		attribs:  map[string]gl.Attrib{},
	}

	// This is mostly when calling New only
	runtime.SetFinalizer(s, func(s *Shader) {
		m.gorge.RunInMain(func() {
			s.destroy()
		})
	})

	// We should only use shader
	if d, ok := r.(*gorge.ShaderData); ok {
		s.upload(d, nil)
		return s
	}
	return s
}

// We need to pass VBO so we can recompile whatever we need
// considering passing renderable or instance instead?
// func (m *shaderManager) GetX(mat gorge.Materialer, vbo *VBO) *Shader {
func (m *shaderManager) GetX(r *gorge.RenderableComponent) *Shader {
	mat := r.Material
	mesh := r.Mesh

	sd := m.defRaw

	res := mat.Resource()
	if s, ok := res.(*gorge.ShaderData); ok && s != nil {
		sd = s
	}

	defines := []string{}

	vbo, _ := m.vbos.Get(mesh)
	vertFormat := vbo.Format
	for _, f := range vertFormat {
		defines = append(defines, f.Define)
	}
	if mesh != nil {
		for d, v := range mesh.Defines() {
			defines = append(defines, d+" "+v)
		}
	}

	for d, v := range mat.Defines() {
		defines = append(defines, d+" "+v)
	}
	// Add Material shader definitions

	return m.Setup(sd, defines)
}

func (m *shaderManager) Setup(sd *gorge.ShaderData, defines []string) *Shader {
	hash := stringHash(sd.Name)
	for _, d := range defines {
		hash ^= stringHash(d)
	}

	shaderVar, ok := m.hashedShaders[sd]
	if !ok {
		shaderVar = map[uint]*Shader{}
		m.hashedShaders[sd] = shaderVar
	}
	// Check if exists.. If Not we compile it with specific defines
	if s, ok := shaderVar[hash]; ok {
		return s
	}

	s := &Shader{
		Name:     sd.Name,
		manager:  m,
		ubos:     map[string]struct{}{},
		uniforms: map[string]*uniform{},
		attribs:  map[string]gl.Attrib{},
	}
	log.Printf("\033[01;37mRecompile hash 0x%X with defines %v for: %v\033[0m", hash, defines, sd)
	s.upload(sd, defines)
	shaderVar[hash] = s
	return s
}

type uniform struct {
	loc     gl.Uniform
	ty      gl.Enum
	sampler bool
	value   interface{}
}

// Shader handle gl program, it holds a local State
type Shader struct {
	Name    string // debug purposes
	manager *shaderManager
	program gl.Program

	ubos map[string]struct{}
	// Extra attribs and uniforms
	attribsHash uint
	attribs     map[string]gl.Attrib
	uniforms    map[string]*uniform
	samplers    []string
}

func (s *Shader) destroy() {
	s.manager.count--
	gl.DeleteProgram(s.program)
}

// Also compiles and links
func (s *Shader) upload(data *gorge.ShaderData, defines []string) {
	if data == nil {
		return
	}
	s.Name = data.Name

	n := bytes.Index(data.Src, []byte("#version"))
	if n == -1 {
		// TODO: fix this one
		panic(fmt.Sprintf("error in %q\nmissing version thing for %q", data.Name, data.Src))
	}
	en := bytes.Index(data.Src[n:], []byte("\n"))

	versionStr := string(data.Src[:en+1])

	shaderSrc := string(data.Src[en:])
	defineSrc := ""
	for _, d := range defines {
		defineSrc += "#define " + d + "\n"
	}
	defineSrc += fmt.Sprintf("#line %d\n", 0)

	vertShader := gl.CreateShader(gl.VERTEX_SHADER)
	{
		// mark := time.Now()
		src := versionStr + defineSrc + "#define VERT_SRC\n" + shaderSrc
		gl.ShaderSource(vertShader, src)
		gl.CompileShader(vertShader)
		// log.Printf("Compiled Vert shader in: %v", time.Since(mark))
		if gl.GetShaderi(vertShader, gl.COMPILE_STATUS) == gl.FALSE {
			// Yep panics :/
			panic(fmt.Errorf("error in %q\nvertex: %s",
				data.Name, gl.GetShaderInfoLog(vertShader),
			))
		}
	}

	fragShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	{
		// mark := time.Now()
		src := versionStr + defineSrc + "#define FRAG_SRC\n" + shaderSrc
		gl.ShaderSource(fragShader, src)
		gl.CompileShader(fragShader)
		// log.Printf("Compiled Frag shader in: %v", time.Since(mark))
		if gl.GetShaderi(fragShader, gl.COMPILE_STATUS) == gl.FALSE {
			panic(fmt.Errorf("error in %q\nfragment: %s",
				data.Name, gl.GetShaderInfoLog(fragShader),
			))
		}
	}

	if !gl.IsValid(s.program) {
		s.manager.count++
		s.program = gl.CreateProgram()
	}

	gl.AttachShader(s.program, vertShader)
	gl.AttachShader(s.program, fragShader)
	gl.LinkProgram(s.program)
	if gl.GetProgrami(s.program, gl.LINK_STATUS) == gl.FALSE {
		log.Println("shader err:", gl.GetProgramInfoLog(s.program))
		return
	}
	gl.DeleteShader(vertShader)
	gl.DeleteShader(fragShader)

	// Setup uniforms and stuff
	// Load names from program

	nUniformBlocks := gl.GetProgrami(s.program, gl.ACTIVE_UNIFORM_BLOCKS)

	uboIndices := map[int32]struct{}{}
	for i := 0; i < nUniformBlocks; i++ {
		name := gl.GetActiveUniformBlockName(s.program, uint32(i))
		// Load indices so we can avoid mapping them in regular uniforms
		nUniforms := []int32{0}
		gl.GetActiveUniformBlockiv(
			s.program,
			uint32(i),
			gl.UNIFORM_BLOCK_ACTIVE_UNIFORMS,
			nUniforms,
		)
		indices := make([]int32, nUniforms[0])
		gl.GetActiveUniformBlockiv(
			s.program,
			uint32(i),
			gl.UNIFORM_BLOCK_ACTIVE_UNIFORM_INDICES,
			indices,
		)
		for _, in := range indices {
			uboIndices[in] = struct{}{}
		}

		/*		{ // Debug indices offsets
				// Debug capture the offsets

				var dataSz int32
				gl.GetActiveUniformBlockiv(
					s.program,
					uint32(i),
					gl.UNIFORM_BLOCK_DATA_SIZE,
					&dataSz,
				)
				log.Printf("Uniform block: name: \033[01;34m%s\033[0m", name)
				log.Printf("uniform data size: %d", dataSz)

				type indOff struct {
					name   string
					loc    int32
					sz     int
					offset int32
				}
				ind := []indOff{}

				for _, in := range indices {
					// Offset Inspect
					name, sz, _ := gl.GetActiveUniform(s.program, uint32(in))
					// uin := uint32(in)
					offset := gl.GetActiveUniformi(s.program, uint32(in), gl.UNIFORM_OFFSET)
					ind = append(ind, indOff{
						name:   name,
						loc:    in,
						offset: offset,
						sz:     sz,
					})
				}
				sort.Slice(ind, func(i, j int) bool { return ind[i].offset < ind[j].offset })
				for i, in := range ind {
					log.Printf("\t uniform offset: %d %s[%d] at %d", i, in.name, in.loc, in.offset)
				}

				log.Printf("active block uniforms: %d", nUniforms)
				log.Printf("active block indices: %d", indices)
				// log.Printf("active block index: %d", index)
			}*/

		s.ubos[name] = struct{}{}
	}

	// Load uniform names
	nUniforms := gl.GetProgrami(s.program, gl.ACTIVE_UNIFORMS)
	s.samplers = []string{}
	s.uniforms = map[string]*uniform{}
	for i := 0; i < nUniforms; i++ {
		if _, ok := uboIndices[int32(i)]; ok {
			continue
		}
		n, sz, ty := gl.GetActiveUniform(s.program, uint32(i))

		names := []string{n}
		// Solve uniform[0] stuff
		// if the size of uniform is bigger than 1 as far as I seen it's
		// registered as somename[0], this will generate names like
		// name: somename[0] size: 4
		//   - somename[0]
		//   - somename[1]
		//   - somename[2]
		//   - somename[3]
		if sz > 1 {
			n = strings.TrimRight(n, "[0]")
			names = names[:0]
			for j := 0; j < sz; j++ {
				names = append(names, fmt.Sprintf("%s[%d]", n, j))
			}
		}
		for _, name := range names {
			sampler := false
			switch ty {
			case gl.SAMPLER_2D, gl.SAMPLER_3D, gl.SAMPLER_CUBE:
				sampler = true
				s.samplers = append(s.samplers, name)
			}
			loc := gl.GetUniformLocation(s.program, name)
			// log.Println("Registering uniform:", name)
			s.uniforms[name] = &uniform{loc: loc, ty: ty, sampler: sampler}
		}
	}
	// Load attrib names
	nAttribs := gl.GetProgrami(s.program, gl.ACTIVE_ATTRIBUTES)
	s.attribs = map[string]gl.Attrib{}
	s.attribsHash = 0
	for i := 0; i < nAttribs; i++ {
		name, _, _ := gl.GetActiveAttrib(s.program, uint32(i))
		loc := gl.GetAttribLocation(s.program, name)
		s.attribs[name] = loc
		s.attribsHash ^= stringHash(name) ^ uint(loc)
	}
}

// Attrib returns the attribute for name
func (s *Shader) Attrib(k string) (gl.Attrib, bool) {
	a, ok := s.attribs[k]
	return a, ok
}

// Set sets a uniform value
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
	if v == nil {
		zeroMat := m32.Mat4{}
		// TODO: more types
		switch u.ty {
		case gl.BOOL:
			gl.Uniform1i(u.loc, 0)
		case gl.FLOAT:
			gl.Uniform1f(u.loc, 0)
		case gl.FLOAT_VEC2, gl.INT_VEC2:
			gl.Uniform2f(u.loc, 0, 0)
		case gl.FLOAT_VEC3, gl.INT_VEC3:
			gl.Uniform3f(u.loc, 0, 0, 0)
		case gl.FLOAT_VEC4, gl.INT_VEC4:
			gl.Uniform4f(u.loc, 0, 0, 0, 0)
		case gl.FLOAT_MAT4:
			gl.UniformMatrix4fv(u.loc, zeroMat[:])
		case gl.SAMPLER_2D, gl.SAMPLER_CUBE:
			// ignore samplers

			//	gl.Uniform1i(u.loc, 0)
			// ignore
			// default:
			//	gl.Uniform1i(u.loc, 0)
		}
		u.value = v
		return
	}

	switch v := v.(type) {
	case bool:
		n := 0
		if v {
			n = 1
		}
		gl.Uniform1i(u.loc, n)
	case uint32:
		gl.Uniform1i(u.loc, int(v))
	case int:
		gl.Uniform1i(u.loc, v)
	case float32:
		gl.Uniform1f(u.loc, v)
	case m32.Vec2:
		gl.Uniform2fv(u.loc, v[:])
	case m32.Vec3:
		gl.Uniform3fv(u.loc, v[:])
	case m32.Vec4:
		gl.Uniform4fv(u.loc, v[:])
	case m32.Mat3:
		gl.UniformMatrix3fv(u.loc, v[:])
	case m32.Mat4:
		gl.UniformMatrix4fv(u.loc, v[:])

	// Pointer support
	case *uint32:
		gl.Uniform1i(u.loc, int(*v))
	case *int:
		gl.Uniform1i(u.loc, *v)
	case *float32:
		gl.Uniform1f(u.loc, *v)
	case *m32.Vec2:
		gl.Uniform2fv(u.loc, (*v)[:])
	case *m32.Vec3:
		gl.Uniform3fv(u.loc, (*v)[:])
	case *m32.Vec4:
		gl.Uniform4fv(u.loc, (*v)[:])
	case *m32.Mat3:
		gl.UniformMatrix3fv(u.loc, (*v)[:])
	case *m32.Mat4:
		gl.UniformMatrix4fv(u.loc, (*v)[:])
	default:
		panic(fmt.Sprintf("not implemented: %T", v))
	}
	u.value = v
}

// Bind binds the shader's program
func (s *Shader) Bind() {
	gl.UseProgram(s.program)
}
