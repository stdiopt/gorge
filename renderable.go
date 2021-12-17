package gorge

type (
	// Materialer interface for material controllers.
	Materialer interface{ Material() *Material }
	// Mesher interface for Mesh controller.
	Mesher interface{ Mesh() *Mesh }
)

// TODO: Add masking constants and camera masking stuff here
// All, None, UI, Debug etc
const (
	CullMaskDefault = 0xFF
	CullMaskUI      = 1 << 8
	CullMaskUIDebug = 1 << 9
)

// RenderableComponent contains info for renderer
// material and mesh
type RenderableComponent struct {
	Name string
	*Material
	*Mesh

	CullMask      uint32
	DisableShadow bool
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
func (r *RenderableComponent) SetCullMask(m uint32) {
	r.CullMask = m
}

// SetDisableShadow sets the disableShadow property which if it is true it
// won't cast a shadow.
func (r *RenderableComponent) SetDisableShadow(b bool) {
	r.DisableShadow = b
}
