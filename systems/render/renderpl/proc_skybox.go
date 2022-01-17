package renderpl

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/static"
	"github.com/stdiopt/gorge/systems/render"
	"github.com/stdiopt/gorge/systems/render/gl"
)

// ProceduralSkybox renders a procedural skybox into a cube
// So we can reuse it in env mapping
func ProceduralSkybox(r *render.Context, next render.StepFunc) render.StepFunc {
	targetFBO := gl.CreateFramebuffer()

	// skybox Object
	shader := r.NewShader(static.Shaders.ProceduralSkybox)
	cubeMesh := gorge.NewMesh(&gorge.MeshData{
		Format:   gorge.VertexFormatP(),
		Vertices: skyboxVert,
	})
	cubeVBO := r.GetVBO(cubeMesh)

	skyboxVAO := gl.CreateVertexArray()
	gl.BindVertexArray(skyboxVAO)
	cubeVBO.BindAttribs(shader)
	gl.BindVertexArray(gl.Null)

	const cubeTexSz = 1024
	// POINT CUBE Textures
	cubeTex := &render.Texture{
		ID:   gl.CreateTexture(),
		Type: gl.TEXTURE_CUBE_MAP,
	}

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, cubeTex.ID)
	for i := 0; i < 6; i++ {
		gl.TexImage2D(
			gl.Enum(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i), 0,
			gl.RGBA, cubeTexSz, cubeTexSz,
			gl.RGBA, gl.UNSIGNED_BYTE, nil) // Webgl
	}
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	// static Position 0,0,0
	// we could add Height
	camTargets := []m32.Mat4{
		m32.LookAt(m32.Vec3{}, m32.Vec3{1, 0, 0}, m32.Vec3{0, -1, 0}),
		m32.LookAt(m32.Vec3{}, m32.Vec3{-1, 0, 0}, m32.Vec3{0, -1, 0}),

		m32.LookAt(m32.Vec3{}, m32.Vec3{0, 1, 0}, m32.Vec3{0, 0, 1}),
		m32.LookAt(m32.Vec3{}, m32.Vec3{0, -1, 0}, m32.Vec3{0, 0, -1}),

		m32.LookAt(m32.Vec3{}, m32.Vec3{0, 0, 1}, m32.Vec3{0, -1, 0}),
		m32.LookAt(m32.Vec3{}, m32.Vec3{0, 0, -1}, m32.Vec3{0, -1, 0}),
	}
	camProj := m32.Perspective(90, 1, .1, 10)

	prevLightDir := m32.Vec3{0, 0, 0}
	return func(s *render.Step) {
		// Find first directional light
		// Do this elsewhere like on prepare
		lightDir := m32.Vec3{-3, -3, -3} // Else we will use this
		for _, rl := range r.Lights.Items() {
			if rl.Light().Type == gorge.LightDirectional {
				// lightDir = m32.Vec3{}.Sub(rl.TransformComponent().WorldPosition()).Normalize()
				// lightDir = rl.Transform().Forward() // WHAT again?
				lightDir = rl.Mat4().MulV4(m32.Vec4{0, 0, -1, 0}).Vec3()
			}
		}
		if lightDir == prevLightDir {
			next(s)
			return
		}
		prevLightDir = lightDir
		// grab Env into cube
		oldVP := s.Viewport
		gl.Viewport(0, 0, cubeTexSz, cubeTexSz)
		gl.BindFramebuffer(gl.FRAMEBUFFER, targetFBO)
		gl.BindVertexArray(skyboxVAO)
		gl.FrontFace(gl.CCW)
		gl.Disable(gl.BLEND)
		shader.Bind()
		shader.Set("lightDir", lightDir)
		for i := 0; i < 6; i++ {
			shader.Set("VP", camProj.Mul(camTargets[i]))

			face := gl.Enum(gl.TEXTURE_CUBE_MAP_POSITIVE_X + i)
			gl.FramebufferTexture2D(
				gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, face, cubeTex.ID, 0,
			)

			gl.Clear(gl.COLOR_BUFFER_BIT)
			r.Draw(gl.TRIANGLES, cubeVBO, 1)
		}
		gl.FrontFace(gl.CW)
		gl.BindVertexArray(gl.Null)
		gl.BindFramebuffer(gl.FRAMEBUFFER, gl.Null)
		gl.Viewport(int(oldVP[0]), int(oldVP[1]), int(oldVP[2]), int(oldVP[3]))

		s.Props["hasEnvMap"] = true
		s.Samplers["envMap"] = cubeTex
		next(s)
	}
}

var skyboxVert = []float32{
	// positions
	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,

	-1.0, -1.0, 1.0,
	-1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0,

	1.0, -1.0, -1.0,
	1.0, -1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, -1.0,
	1.0, -1.0, -1.0,

	-1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, -1.0, 1.0,
	-1.0, -1.0, 1.0,

	-1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,

	-1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, 1.0,
}
