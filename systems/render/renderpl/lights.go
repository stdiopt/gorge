package renderpl

import (
	"fmt"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/static"
	"github.com/stdiopt/gorge/systems/render"
	"github.com/stdiopt/gorge/systems/render/bufutil"
	"github.com/stdiopt/gorge/systems/render/gl"
)

// Default light stuff.
const (
	MaxLights      = 8
	pointShadowSz  = 1024
	dirShadowSz    = 2048 // too much?
	dirNearPlane   = 100
	dirFarPlane    = -100
	dirSz          = 40
	pointNearPlane = 1
)

// lights prepare lights and process any depth Map for shadows
type lights struct {
	renderer        *render.Context
	depthFBO        gl.Framebuffer
	depthDirShader  *render.Shader // Or Spot
	depthCubeShader *render.Shader

	// PerNLight
	depthCubeTex []*render.Texture
	depth2DTex   []*render.Texture
}

// PrepareLights prepares a rendering stage to setup light related uniforms and draw
// geometry to depht maps for shadowing.
func PrepareLights(r *render.Context, next render.StepFunc) render.StepFunc {
	lightNames := []lightName{}
	depthNames := []depthName{}

	lightSpec := bufutil.OffsetSpec{"nLights": 0}
	for i := 0; i < MaxLights; i++ {
		depthNames = append(depthNames, depthName{
			Depth2D:   fmt.Sprintf("depth2D[%d]", i),
			DepthCube: fmt.Sprintf("depthCube[%d]", i),
		})
		lightNames = append(lightNames, lightName{
			Position:     fmt.Sprintf("u_Lights[%d].position", i),
			Direction:    fmt.Sprintf("u_Lights[%d].direction", i),
			Range:        fmt.Sprintf("u_Lights[%d].range", i),
			Color:        fmt.Sprintf("u_Lights[%d].color", i),
			Intensity:    fmt.Sprintf("u_Lights[%d].intensity", i),
			InnerConeCos: fmt.Sprintf("u_Lights[%d].innerConeCos", i),
			OuterConeCos: fmt.Sprintf("u_Lights[%d].outerConeCos", i),
			Matrix:       fmt.Sprintf("u_Lights[%d].matrix", i),
			Type:         fmt.Sprintf("u_Lights[%d].type", i),
			DepthIndex:   fmt.Sprintf("u_Lights[%d].depthIndex", i),
		})
		lightOff := (i * 144)
		lightSpec[lightNames[i].Position] = lightOff + 16
		lightSpec[lightNames[i].Direction] = lightOff + 32
		lightSpec[lightNames[i].Range] = lightOff + 44
		lightSpec[lightNames[i].Color] = lightOff + 48
		lightSpec[lightNames[i].Intensity] = lightOff + 60
		lightSpec[lightNames[i].InnerConeCos] = lightOff + 64
		lightSpec[lightNames[i].OuterConeCos] = lightOff + 68
		lightSpec[lightNames[i].Type] = lightOff + 72
		lightSpec[lightNames[i].Matrix] = lightOff + 80
		lightSpec[lightNames[i].DepthIndex] = lightOff + 144
	}

	lightsUBO := bufutil.NewNamedOffset(
		bufutil.NewCached[byte](
			r.NewBuffer(gl.UNIFORM_BUFFER, gl.DYNAMIC_DRAW),
		),
		1168,
		// sz
		lightSpec,
	)

	// This is nearly the same except we have some transformation here
	// we can setup this as material and add some defines
	depthDirShader := r.NewShader(static.Shaders.Depth)
	depthCubeShader := r.NewShader(static.Shaders.DepthCube)

	// POINT CUBE Textures
	depthCubeTex := []*render.Texture{}
	for ti := 0; ti < MaxLights; ti++ {
		depthTex := &render.Texture{
			ID:   gl.CreateTexture(),
			Type: gl.TEXTURE_CUBE_MAP,
		}
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, depthTex.ID)
		for i := 0; i < 6; i++ {
			gl.TexImage2D(
				gl.Enum(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i), 0,
				gl.DEPTH_COMPONENT16, pointShadowSz, pointShadowSz,
				gl.DEPTH_COMPONENT, gl.UNSIGNED_SHORT, nil) // Webgl
		}
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
		depthCubeTex = append(depthCubeTex, depthTex)
	}

	// DIRECTIONAL 2D Textures
	depth2DTex := []*render.Texture{}
	for ti := 0; ti < MaxLights; ti++ {
		depthTex := &render.Texture{
			ID:   gl.CreateTexture(),
			Type: gl.TEXTURE_2D,
		}
		gl.BindTexture(gl.TEXTURE_2D, depthTex.ID)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT32F, dirShadowSz, dirShadowSz, gl.DEPTH_COMPONENT, gl.FLOAT, nil) // Webgl
		// gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT, shadowSz, shadowSz, gl.DEPTH_COMPONENT, gl.FLOAT, nil) // Glfw
		// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

		// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

		// Experiment in some gl versions

		if gl.Global().Impl() != "wasm" { // TODO: experiment and improve this string comparison to const
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)
			borderColor := []float32{1.0, 1.0, 1.0, 1.0}
			gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, borderColor)
		}

		depth2DTex = append(depth2DTex, depthTex)
	}

	depthFBO := gl.CreateFramebuffer()

	ls := &lights{
		renderer:        r,
		depthFBO:        depthFBO,
		depthDirShader:  depthDirShader,
		depthCubeShader: depthCubeShader,

		depthCubeTex: depthCubeTex,
		depth2DTex:   depth2DTex,

		// lightsUBO: lightsUBO,

		// lightNames: lightNames,
	}

	return func(s *render.Step) {
		depthCubeIndex := 0
		depth2DIndex := 0
		lights := r.Lights.Items()
		for ti := 0; ti < len(lights); ti++ {
			lightDepthIndex := -1
			light := lights[ti]
			l := light.Light()
			// t := light.Transform()
			mat4 := light.Mat4()
			dir := mat4.MulV4(gm.Vec4{0, 0, -1, 0}).Vec3()
			pos := mat4.Col(3).Vec3()

			switch l.Type {
			case gorge.LightDirectional:
				lightsUBO.WriteOffset(lightNames[ti].Type, int32(0))
				if l.DisableShadow {
					break
				}
				lightDepthIndex = depth2DIndex
				mat4 = ls.processDepth2D(s, light, depth2DIndex)
				s.Samplers[depthNames[depth2DIndex].Depth2D] = ls.depth2DTex[depth2DIndex]
				depth2DIndex++
			case gorge.LightPoint:
				lightsUBO.WriteOffset(lightNames[ti].Type, int32(1))
				if l.DisableShadow {
					break
				}
				lightDepthIndex = depthCubeIndex
				ls.processDepthCube(s, light, depthCubeIndex)
				s.Samplers[depthNames[depthCubeIndex].DepthCube] = ls.depthCubeTex[depthCubeIndex]
				depthCubeIndex++
			case gorge.LightSpot:
				// Same as Directional but with a projection tex
				lightsUBO.WriteOffset(lightNames[ti].Type, int32(2))
			}
			lightsUBO.WriteOffset(lightNames[ti].Position, pos)
			lightsUBO.WriteOffset(lightNames[ti].Direction, dir)
			lightsUBO.WriteOffset(lightNames[ti].Color, l.Color)
			lightsUBO.WriteOffset(lightNames[ti].Intensity, l.Intensity)
			lightsUBO.WriteOffset(lightNames[ti].Range, l.Range)
			lightsUBO.WriteOffset(lightNames[ti].InnerConeCos, l.InnerConeCos)
			lightsUBO.WriteOffset(lightNames[ti].OuterConeCos, l.OuterConeCos)
			lightsUBO.WriteOffset(lightNames[ti].DepthIndex, int32(lightDepthIndex))
			lightsUBO.WriteOffset(lightNames[ti].Matrix, mat4)
		}
		// This should be on define
		lightsUBO.WriteOffset("nLights", int32(r.Lights.Len()))
		lightsUBO.Flush()
		s.Ubos["Lights"] = lightsUBO.ID()
		next(s)
	}
}

// This render a depth cube based on light to target DepthIndex
func (s *lights) processDepthCube(ri *render.Step, light render.Light, di int) {
	// Check cached light and render if needed
	pos := light.Mat4().Col(3).Vec3()
	farPlane := light.Light().Range
	lightMat := []gm.Mat4{
		gm.LookAt(pos, pos.Add(gm.Vec3{1, 0, 0}), gm.Vec3{0, -1, 0}),
		gm.LookAt(pos, pos.Add(gm.Vec3{-1, 0, 0}), gm.Vec3{0, -1, 0}),

		gm.LookAt(pos, pos.Add(gm.Vec3{0, 1, 0}), gm.Vec3{0, 0, 1}),
		gm.LookAt(pos, pos.Add(gm.Vec3{0, -1, 0}), gm.Vec3{0, 0, -1}),

		gm.LookAt(pos, pos.Add(gm.Vec3{0, 0, 1}), gm.Vec3{0, -1, 0}),
		gm.LookAt(pos, pos.Add(gm.Vec3{0, 0, -1}), gm.Vec3{0, -1, 0}),
	}
	lproj := gm.Perspective(90, 1, pointNearPlane, farPlane)

	// 1 - RENDER DEPTH TO FBO
	gl.Disable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(true)
	gl.Viewport(0, 0, pointShadowSz, pointShadowSz)
	gl.BindFramebuffer(gl.FRAMEBUFFER, s.depthFBO)
	shdr := s.depthCubeShader
	shdr.Bind()

	shdr.Set("u_AlphaCutoff", 0.01)
	shdr.Set("farPlane", farPlane)
	shdr.Set("lightPos", pos)
	// 6 Pass that could be reduced to one with geometry shader
	for i := 0; i < 6; i++ {
		shdr.Set("view", lproj.Mul(lightMat[i]))
		face := gl.Enum(gl.TEXTURE_CUBE_MAP_POSITIVE_X + i)

		/*gl.FramebufferTextureLayer(
			gl.FRAMEBUFFER,
			gl.DEPTH_ATTACHMENT,
			s.depthCubeArrayTex.ID,
			0, (di*6)+i,
		)*/

		gl.FramebufferTexture2D(
			gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, face, s.depthCubeTex[di].ID, 0,
		)

		gl.Clear(gl.DEPTH_BUFFER_BIT)
		// THIS need to figure out
		s.renderer.PassShadow(ri, shdr)

	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, gl.Null)
}

// Returns the mat4 used to render stuff
func (s *lights) processDepth2D(ri *render.Step, light render.Light, di int) gm.Mat4 {
	m4 := ri.Camera.Mat4()
	camPos := m4.Col(3).Vec3()
	camForward := m4.MulV4(gm.Vec4{0, 0, -1, 0}).Vec3()
	dir := light.Mat4().MulV4(gm.Vec4{0, 0, -1, 0}).Vec3()

	trans := gorge.TransformIdent()
	trans.LookDir(dir, gm.Up())
	trans.SetPositionv(camPos.Add(camForward).Mul(3))

	// Depends on light
	proj := gm.Ortho(-dirSz, dirSz, -dirSz, dirSz, dirNearPlane, dirFarPlane)

	lightMatrix := proj.Mul(trans.Inv())

	gl.Disable(gl.CULL_FACE)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(true)
	gl.BindFramebuffer(gl.FRAMEBUFFER, s.depthFBO)

	gl.Viewport(0, 0, dirShadowSz, dirShadowSz)
	shdr := s.depthDirShader
	shdr.Bind()
	shdr.Set("view", lightMatrix)
	shdr.Set("u_AlphaCutoff", 0.01)
	{
		gl.FramebufferTexture2D(
			gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, s.depth2DTex[di].ID, 0,
		)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
		s.renderer.PassShadow(ri, shdr)
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, gl.Null)

	return lightMatrix
}

type lightName struct {
	Position     string
	Direction    string
	Range        string
	Color        string
	Intensity    string
	InnerConeCos string
	OuterConeCos string
	Type         string
	Matrix       string
	DepthIndex   string
}

type depthName struct {
	Depth2D   string
	DepthCube string
}
