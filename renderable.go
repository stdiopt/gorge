package gorge

type (
	// Materialer interface for material controllers.
	Materialer interface{ Material() *Material }
	// Mesher interface for Mesh controller.
	Mesher interface{ Mesh() *Mesh }
)

// TODO: Add masking constants and camera masking stuff here
// All, None, UI, Debug etc

// CullMaskFlags flag type for culling masks.
type CullMaskFlags uint32

// CullMask defaults
const (
	CullMaskDefault = CullMaskFlags(0xFF)
	CullMaskUI      = CullMaskFlags(1 << 8)
	CullMaskUIDebug = CullMaskFlags(1 << 9)
)

// RenderableComponent contains info for renderer
// material and mesh
type RenderableComponent struct {
	GPU
	Name string
	*Material
	*Mesh

	// Maybe move to material?!
	// Layering and Culling are usually per GameObject in unity.
	Order    int
	CullMask CullMaskFlags
}

// NewRenderableComponent returns a new renderable component
func NewRenderableComponent(mesh Mesher, mat Materialer) *RenderableComponent {
	return &RenderableComponent{
		Mesh:     mesh.Mesh(),
		Material: mat.Material(),
	}
}

// Renderable returns the renderable component
func (r *RenderableComponent) Renderable() *RenderableComponent { return r }

// SetMaterial sets the material.
func (r *RenderableComponent) SetMaterial(m Materialer) {
	r.Material = m.Material()
}

// SetMesh sets the mesh.
func (r *RenderableComponent) SetMesh(m Mesher) {
	r.Mesh = m.Mesh()
}

// SetCullMask will set the cull mask which is used in conjunction with camera
// mask to filter which renderables will render in each camera.
func (r *RenderableComponent) SetCullMask(m CullMaskFlags) {
	r.CullMask = m
}

// SetOrder sets the render order lower will render first.
func (r *RenderableComponent) SetOrder(o int) {
	r.Order = o
}
