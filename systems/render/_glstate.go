package render

import "github.com/stdiopt/gorge/systems/render/gl"

type stencilFunc struct {
	Func gl.Enum
	Ref  int
	Mask uint32
}

type glstate struct {
	toggles map[gl.Enum]bool
	// Toggles
	enableBlend    bool
	enableDepth    bool
	enableCullFace bool
	enableStencil  bool
	enableScissor  bool

	blendFunc [2]gl.Enum

	depthMask bool
	depthFunc gl.Enum

	colorMask [4]bool

	cullFace  gl.Enum
	frontFace gl.Enum

	stencilMask uint32
	stencilFunc stencilFunc
	stencilOp   [3]gl.Enum

	scissor [4]int32

	viewport [4]int
}

func (s *glstate) EnableBlend(blend bool) {
	if s.enableBlend == blend {
		return
	}
	s.enableBlend = blend
	if blend {
		gl.Enable(gl.BLEND)
	} else {
		gl.Disable(gl.BLEND)
	}
}

func (s *glstate) SetBlendFunc(src, dst gl.Enum) {
	if s.blendFunc[0] == src && s.blendFunc[1] == dst {
		return
	}

	s.blendFunc[0] = src
	s.blendFunc[1] = dst
	gl.BlendFunc(src, dst)
}

func (s *glstate) EnableDepth(depth bool) {
	if s.enableDepth == depth {
		return
	}
	s.enableDepth = depth
	if depth {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
}

func (s *glstate) SetDepthMask(depthMask bool) {
	if s.depthMask == depthMask {
		return
	}
	s.depthMask = depthMask
	gl.DepthMask(depthMask)
}

func (s *glstate) SetDepthFunc(depthFunc gl.Enum) {
	if s.depthFunc == depthFunc {
		return
	}
	s.depthFunc = depthFunc
	gl.DepthFunc(depthFunc)
}

func (s *glstate) SetColorMask(r, g, b, a bool) {
	if s.colorMask[0] == r && s.colorMask[1] == g && s.colorMask[2] == b && s.colorMask[3] == a {
		return
	}
	s.colorMask[0] = r
	s.colorMask[1] = g
	s.colorMask[2] = b
	s.colorMask[3] = a
	gl.ColorMask(r, g, b, a)
}

func (s *glstate) EnableCullFace(cullFace bool) {
	if s.enableCullFace == cullFace {
		return
	}
	s.enableCullFace = cullFace
	if cullFace {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
}

func (s *glstate) SetCullFace(cullFace gl.Enum) {
	if s.cullFace == cullFace {
		return
	}
	s.cullFace = cullFace
	gl.CullFace(cullFace)
}

func (s *glstate) SetFrontFace(frontFace gl.Enum) {
	if s.frontFace == frontFace {
		return
	}
	s.frontFace = frontFace
	gl.FrontFace(frontFace)
}

func (s *glstate) EnableStencil(stencil bool) {
	if s.enableStencil == stencil {
		return
	}
	s.enableStencil = stencil
	if stencil {
		gl.Enable(gl.STENCIL_TEST)
	} else {
		gl.Disable(gl.STENCIL_TEST)
	}
}

func (s *glstate) SetStencilMask(mask uint32) {
	if s.stencilMask == mask {
		return
	}
	s.stencilMask = mask
	gl.StencilMask(mask)
}

func (s *glstate) SetStencilFunc(f gl.Enum, ref int, mask uint32) {
	if s.stencilFunc.Func == f && s.stencilFunc.Ref == ref && s.stencilFunc.Mask == mask {
		return
	}
	s.stencilFunc.Func = f
	s.stencilFunc.Ref = ref
	s.stencilFunc.Mask = mask
	gl.StencilFunc(f, ref, mask)
}

func (s *glstate) SetStencilOp(fail, zfail, zpass gl.Enum) {
	if s.stencilOp[0] == fail && s.stencilOp[1] == zfail && s.stencilOp[2] == zpass {
		return
	}
	s.stencilOp = [3]gl.Enum{fail, zfail, zpass}
	gl.StencilOp(fail, zfail, zpass)
}

func (s *glstate) EnableScissor(scissor bool) {
	if s.enableScissor == scissor {
		return
	}
	s.enableScissor = scissor
	if scissor {
		gl.Enable(gl.SCISSOR_TEST)
	} else {
		gl.Disable(gl.SCISSOR_TEST)
	}
}

func (s *glstate) SetScissor(x, y, width, height int32) {
	if s.scissor[0] == x && s.scissor[1] == y && s.scissor[2] == width && s.scissor[3] == height {
		return
	}
	s.scissor = [4]int32{x, y, width, height}
	gl.Scissor(x, y, width, height)
}

// SetViewport is not cached
func (s *glstate) SetViewport(x, y, width, height int32) {
	if s.viewport[0] == x && s.viewport[1] == y && s.viewport[2] == width && s.viewport[3] == height {
		return
	}
	s.viewport = [4]int32{x, y, width, height}
	gl.Viewport(x, y, width, height)
}

func (s *glstate) Apply() {
	if s.enableBlend {
		gl.Enable(gl.BLEND)
	} else {
		gl.Disable(gl.BLEND)
	}
	gl.BlendFunc(s.blendFunc[0], s.blendFunc[1])

	if s.enableDepth {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
	gl.DepthMask(s.depthMask)

	gl.ColorMask(s.colorMask[0], s.colorMask[1], s.colorMask[2], s.colorMask[3])

	if s.enableCullFace {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
	gl.CullFace(s.cullFace)
	gl.FrontFace(s.frontFace)

	if s.enableStencil {
		gl.Enable(gl.STENCIL_TEST)
	} else {
		gl.Disable(gl.STENCIL_TEST)
	}
	gl.StencilMask(s.stencilMask)
	gl.StencilFunc(s.stencilFunc.Func, s.stencilFunc.Ref, s.stencilFunc.Mask)
	gl.StencilOp(s.stencilOp[0], s.stencilOp[1], s.stencilOp[2])

	if s.enableScissor {
		gl.Enable(gl.SCISSOR_TEST)
	} else {
		gl.Disable(gl.SCISSOR_TEST)
	}

	// gl.Scissor(s.scissor[0], s.scissor[1], s.scissor[2], s.scissor[3])
}
