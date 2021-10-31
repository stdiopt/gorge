package debug

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/gorgeutil"
)

// MakeGrid creates a generic 3D grid.
func MakeGrid(gridSize float32) *gorgeutil.Renderable {
	shdr := gorge.ShaderData{
		Name: "_native/grid.glsl",
		Src: []byte(`#version 300 es
#ifdef VERT_SRC
	layout (location = 0) in vec3  a_Position;
	layout (location = 1) in float a_Color;
	layout (location = 4) in mat4  a_Transform;

	out vec4 Color;
	
	layout(std140) uniform Camera {
		mat4 VP;
		vec3 ambient;
		vec3 viewPos;
	};

	const vec4 colors[] = vec4[] (
		vec4(1, 0, 0, .7),
		vec4(0, 1, 0, .7),
		vec4(0, 0, 1, .7),
		vec4(1, 1, 1, .2)
	);

	void main() {
		Color = colors[ int(a_Color) ];
		vec3 FragPos = vec3(a_Transform * vec4(a_Position, 1.0));

		gl_Position = VP * vec4(FragPos, 1.0);
	}
#endif
#ifdef FRAG_SRC
	precision mediump float;

	out vec4 FragColor;
	in vec4 Color;

	void main() {
		FragColor = vec4(Color.rgb * Color.a, Color.a);
	}
#endif`),
	}

	mat := gorge.NewShaderMaterial(&shdr)
	if gridSize == 0 {
		gridSize = .5
	}
	sz := float32(70)
	t := sz * gridSize

	meshData := &gorge.MeshData{}

	c := float32(3)
	for i := 1; i <= int(sz); i++ {
		p := float32(i) * gridSize
		meshData.Vertices = append(meshData.Vertices,
			p, 0, t, c, p, 0, -t, c,
			-p, 0, t, c, -p, 0, -t, c,
			-t, 0, p, c, t, 0, p, c,
			-t, 0, -p, c, t, 0, -p, c,
		)
	}
	meshData.Vertices = append(meshData.Vertices,
		-t, 0, 0, 0, t, 0, 0, 0,
		0, 0, t, 2, 0, 0, -t, 2,
	)

	meshData.Format = gorge.VertexFormat{
		gorge.VertexAttrib(3, "a_Position", "HAS_POSITION"),
		gorge.VertexAttrib(1, "a_Color", "HAS_COLOR"),
	}
	mesh := gorge.NewMesh(meshData)
	mesh.DrawMode = gorge.DrawLines
	p := gorgeutil.NewRenderable(mesh, mat)
	p.Renderable().DisableShadow = true
	return p
}
