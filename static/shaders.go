package static

import (
	"fmt"

	"github.com/stdiopt/gorge"
)

// Might be good to use resource system in the other systems instead of these.

func genShaderData(name string) *gorge.ShaderData {
	return &gorge.ShaderData{
		Name: fmt.Sprintf("_gorge/%s", name),
		Src:  MustData(name),
	}
}

type shaders struct {
	Default          *gorge.ShaderData
	DefaultNew       *gorge.ShaderData
	PBRTexture       *gorge.ShaderData
	Unlit            *gorge.ShaderData
	UnlitDebug       *gorge.ShaderData
	UnlitAdditive    *gorge.ShaderData
	CubeEnv          *gorge.ShaderData
	Depth            *gorge.ShaderData
	DepthCube        *gorge.ShaderData
	ProceduralSkybox *gorge.ShaderData
	Quad             *gorge.ShaderData
	QuadDepth        *gorge.ShaderData
	UI               *gorge.ShaderData
}

type fonts struct {
	Default []byte
}

// Shaders default in binary shaders.
var Shaders = shaders{
	Default:          genShaderData("shaders/default.glsl"),
	DefaultNew:       genShaderData("shaders/default_new.glsl"),
	Unlit:            genShaderData("shaders/unlit.glsl"),
	UnlitDebug:       genShaderData("shaders/unlit-debug.glsl"),
	UnlitAdditive:    genShaderData("shaders/unlit_additive.glsl"),
	CubeEnv:          genShaderData("shaders/cube_env.glsl"),
	Depth:            genShaderData("shaders/depth.glsl"),
	DepthCube:        genShaderData("shaders/depth_cube.glsl"),
	ProceduralSkybox: genShaderData("shaders/skybox_proc.glsl"),
	Quad:             genShaderData("shaders/quad.glsl"),
	QuadDepth:        genShaderData("shaders/quad_depth.glsl"),
	UI:               genShaderData("shaders/ui.glsl"),
}

// Fonts default in binary fonts.
var Fonts = fonts{
	Default: MustData("fonts/font.ttf"),
}
