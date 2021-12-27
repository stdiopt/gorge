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

// CastShadow flags
type CastShadow int

// CastShadow defaults
const (
	CastShadowEnabled = CastShadow(iota)
	CastShadowDisabled
)

// RenderableComponent contains info for renderer
// material and mesh
type RenderableComponent struct {
	Name string
	*Material
	*Mesh

	Order      int
	CullMask   CullMaskFlags
	CastShadow CastShadow
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

// SetCastShadow sets the castshadow if CastShadowDisabled it will be disabled
// for all lights won't cast a shadow, CastShadowEnabled will enable it.
func (r *RenderableComponent) SetCastShadow(s CastShadow) {
	r.CastShadow = s
}

// SetOrder sets the render order lower will render first.
func (r *RenderableComponent) SetOrder(o int) {
	r.Order = o
}
