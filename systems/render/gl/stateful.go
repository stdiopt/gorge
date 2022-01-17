package gl

type wrapper = Wrapper

type stencilFunc struct {
	Func Enum
	Ref  int
	Mask uint32
}

type stateful struct {
	*wrapper

	// Toggles
	enableBlend    bool
	enableDepth    bool
	enableCullFace bool
	enableStencil  bool
	enableScissor  bool

	blendFunc [2]Enum

	depthMask bool
	depthFunc Enum

	clearColor [4]float32
	colorMask  [4]bool

	cullFace  Enum
	frontFace Enum

	stencilMask uint32
	stencilFunc stencilFunc
	stencilOp   [3]Enum

	scissor [4]int32

	viewport [4]int32
}

func (s *stateful) Enable(cap Enum) {
	switch cap {
	case BLEND:
		if s.enableBlend {
			return
		}
		s.enableBlend = true
	case DEPTH_TEST:
		if s.enableDepth {
			return
		}
		s.enableDepth = true
	case CULL_FACE:
		if s.enableCullFace {
			return
		}
		s.enableCullFace = true
	case STENCIL_TEST:
		if s.enableStencil {
			return
		}
		s.enableStencil = true
	case SCISSOR_TEST:
		if s.enableScissor {
			return
		}
		s.enableScissor = true
	}
	s.wrapper.Enable(cap)
}

func (s *stateful) Disable(cap Enum) {
	switch cap {
	case BLEND:
		if !s.enableBlend {
			return
		}
		s.enableBlend = false
	case DEPTH_TEST:
		if !s.enableDepth {
			return
		}
		s.enableDepth = false
	case CULL_FACE:
		if !s.enableCullFace {
			return
		}
		s.enableCullFace = false
	case STENCIL_TEST:
		if !s.enableStencil {
			return
		}
		s.enableStencil = false
	case SCISSOR_TEST:
		if !s.enableScissor {
			return
		}
		s.enableScissor = false
	}
	s.wrapper.Disable(cap)
}

func (s *stateful) BlendFunc(sfactor, dfactor Enum) {
	if s.blendFunc[0] == sfactor && s.blendFunc[1] == dfactor {
		return
	}
	s.blendFunc[0] = sfactor
	s.blendFunc[1] = dfactor
	s.wrapper.BlendFunc(sfactor, dfactor)
}

func (s *stateful) DepthMask(flag bool) {
	if s.depthMask == flag {
		return
	}
	s.depthMask = flag
	s.wrapper.DepthMask(flag)
}

func (s *stateful) DepthFunc(f Enum) {
	if s.depthFunc == f {
		return
	}
	s.depthFunc = f
	s.wrapper.DepthFunc(f)
}

func (s *stateful) ClearColor(r, g, b, a float32) {
	if s.clearColor[0] == r && s.clearColor[1] == g && s.clearColor[2] == b && s.clearColor[3] == a {
		return
	}
	s.clearColor = [4]float32{r, g, b, a}
	s.wrapper.ClearColor(r, g, b, a)
}

func (s *stateful) ColorMask(r, g, b, a bool) {
	if s.colorMask[0] == r && s.colorMask[1] == g && s.colorMask[2] == b && s.colorMask[3] == a {
		return
	}
	s.colorMask[0] = r
	s.colorMask[1] = g
	s.colorMask[2] = b
	s.colorMask[3] = a
	s.wrapper.ColorMask(r, g, b, a)
}

func (s *stateful) CullFace(f Enum) {
	if s.cullFace == f {
		return
	}
	s.cullFace = f
	s.wrapper.CullFace(f)
}

func (s *stateful) FrontFace(f Enum) {
	if s.frontFace == f {
		return
	}
	s.frontFace = f
	s.wrapper.FrontFace(f)
}

func (s *stateful) StencilMask(mask uint32) {
	if s.stencilMask == mask {
		return
	}
	s.stencilMask = mask
	s.wrapper.StencilMask(mask)
}

func (s *stateful) StencilFunc(f Enum, ref int, mask uint32) {
	if s.stencilFunc.Func == f && s.stencilFunc.Ref == ref && s.stencilFunc.Mask == mask {
		return
	}
	s.stencilFunc = stencilFunc{f, ref, mask}
	s.wrapper.StencilFunc(f, ref, mask)
}

func (s *stateful) StencilOp(sfail, dfail, dpass Enum) {
	if s.stencilOp[0] == sfail && s.stencilOp[1] == dfail && s.stencilOp[2] == dpass {
		return
	}
	s.stencilOp = [3]Enum{sfail, dfail, dpass}
	s.wrapper.StencilOp(sfail, dfail, dpass)
}

func (s *stateful) Scissor(x, y, width, height int32) {
	if s.scissor[0] == x && s.scissor[1] == y && s.scissor[2] == width && s.scissor[3] == height {
		return
	}
	s.scissor = [4]int32{x, y, width, height}
	s.wrapper.Scissor(x, y, width, height)
}

func (s *stateful) Viewport(x, y, width, height int32) {
	if s.viewport[0] == x && s.viewport[1] == y && s.viewport[2] == width && s.viewport[3] == height {
		return
	}
	s.viewport = [4]int32{x, y, width, height}
	s.wrapper.Viewport(x, y, width, height)
}

func (s *stateful) init() {
	s.enableDepth = true
	s.enableCullFace = true
	s.enableBlend = false
	s.enableStencil = false
	s.enableScissor = false
	s.blendFunc = [2]Enum{ONE, ONE_MINUS_SRC_ALPHA}
	s.depthMask = true
	s.depthFunc = LESS
	s.clearColor = [4]float32{0, 0, 0, 0}
	s.colorMask = [4]bool{true, true, true, true}
	s.cullFace = BACK
	s.frontFace = CCW
	s.stencilMask = 0xFF
	s.stencilFunc = stencilFunc{NEVER, 0, 0xFF}
	s.stencilOp = [3]Enum{KEEP, KEEP, KEEP}
	s.scissor = [4]int32{0, 0, 0, 0}
	s.viewport = [4]int32{0, 0, 0, 0}

	s.wrapper.Enable(DEPTH_TEST)
	s.wrapper.Enable(CULL_FACE)
	s.wrapper.Disable(BLEND)
	s.wrapper.Disable(STENCIL_TEST)
	s.wrapper.Disable(SCISSOR_TEST)

	s.wrapper.BlendFunc(s.blendFunc[0], s.blendFunc[1])
	s.wrapper.DepthMask(s.depthMask)
	s.wrapper.DepthFunc(s.depthFunc)
	s.wrapper.ClearColor(s.clearColor[0], s.clearColor[1], s.clearColor[2], s.clearColor[3])
	s.wrapper.ColorMask(s.colorMask[0], s.colorMask[1], s.colorMask[2], s.colorMask[3])
	s.wrapper.CullFace(s.cullFace)
	s.wrapper.FrontFace(s.frontFace)
	s.wrapper.StencilMask(s.stencilMask)
	s.wrapper.StencilFunc(s.stencilFunc.Func, s.stencilFunc.Ref, s.stencilFunc.Mask)
	s.wrapper.StencilOp(s.stencilOp[0], s.stencilOp[1], s.stencilOp[2])
	s.wrapper.Scissor(s.scissor[0], s.scissor[1], s.scissor[2], s.scissor[3])
	s.wrapper.Viewport(s.viewport[0], s.viewport[1], s.viewport[2], s.viewport[3])
}
