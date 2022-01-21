package gorge

import (
	"math"

	"github.com/stdiopt/gorge/math/gm"
)

// ProjectionType camera projection type.
type ProjectionType int

func (t ProjectionType) String() string {
	switch t {
	case ProjectionPerspective:
		return "Perspective"
	case ProjectionOrtho:
		return "Ortho"
	default:
		return "<invalid>"
	}
}

// Projection types for camera.
const (
	ProjectionPerspective = ProjectionType(iota)
	ProjectionOrtho
)

// ClearType method on how the camera clears.
type ClearType int

func (t ClearType) String() string {
	switch t {
	case ClearSkybox:
		return "ClearSkybox"
	case ClearColor:
		return "ClearColor"
	case ClearDepthOnly:
		return "ClearDepthOnly"
	case ClearNothing:
		return "ClearNothing"
	default:
		return "<invalid>"
	}
}

// Clear types
const (
	ClearColor = ClearType(iota)
	ClearDepthOnly
	ClearNothing
	ClearSkybox
)

// CameraComponent thing
type CameraComponent struct {
	Name           string
	ProjectionType ProjectionType
	CullMask       CullMaskFlags

	Fov       float32
	OrthoSize float32

	AspectRatio float32
	Near        float32
	Far         float32

	ClearFlag     ClearType
	ClearMaterial *Material
	Order         int
	Viewport      gm.Vec4
	ClearColor    gm.Vec3

	// Consider this to custom clear other buffers?!
	// ClearDepth   float32
	// ClearStencil byte

	// Or Forward, deferred, etc type of pass
	// RenderShadow bool
	// Or sky box
}

// NewCameraComponent returns a new default camera Component
func NewCameraComponent(name string) *CameraComponent {
	c := &CameraComponent{
		Name:     name,
		Viewport: gm.Vec4{0, 0, 1, 1},
	}
	c.SetPerspective(math.Pi/4, .1, 1000)
	return c
}

// Camera returns camera component
// We actually only need Projection :/
func (c *CameraComponent) Camera() *CameraComponent { return c }

// Projection returns the projection matrix with default aspect ratio based
// on registered size
func (c CameraComponent) Projection(screenSize gm.Vec2) gm.Mat4 {
	aspectRatio := c.AspectRatio
	if aspectRatio == 0 {
		vp := c.CalcViewport(screenSize)
		aspectRatio = vp[2] / vp[3]
	}

	return c.ProjectionWithAspect(aspectRatio)
}

// ProjectionWithAspect Sets the projection matrices with given aspect ratio
func (c CameraComponent) ProjectionWithAspect(aspect float32) gm.Mat4 {
	if c.ProjectionType == ProjectionPerspective {
		return gm.Perspective(c.Fov, aspect, c.Near, c.Far)
	}

	halfH := c.OrthoSize * aspect * .5
	halfV := c.OrthoSize * .5

	bottom := -halfV
	top := halfV

	left := -halfH
	right := halfH

	// Ortho
	return gm.Ortho(left, right, bottom, top, c.Near, c.Far)
}

// SetPerspective resets projection matrix to perspective
func (c *CameraComponent) SetPerspective(fov, near, far float32) {
	c.ProjectionType = ProjectionPerspective
	c.Fov = fov
	c.Near = near
	c.Far = far
}

// SetOrtho sets ortho matrix
func (c *CameraComponent) SetOrtho(size, near, far float32) {
	c.ProjectionType = ProjectionOrtho
	c.OrthoSize = size
	c.Near = near
	c.Far = far
}

// SetAspectRatio sets the camera aspect ratio
func (c *CameraComponent) SetAspectRatio(a float32) {
	c.AspectRatio = a
}

// SetClearFlag sets the clear flag and returns
func (c *CameraComponent) SetClearFlag(clr ClearType) {
	c.ClearFlag = clr
}

// SetClearColor for the camera.
func (c *CameraComponent) SetClearColor(r, g, b float32) {
	c.ClearColor = gm.Vec3{r, g, b}
}

// SetCullMask for camera, only specific renderables that masks this cullmask
// will render with this camera.
func (c *CameraComponent) SetCullMask(m CullMaskFlags) {
	c.CullMask = m
}

// SetOrder camera rendering order, higher will render later.
func (c *CameraComponent) SetOrder(n int) {
	c.Order = n
}

// SetViewport sets the viewport for camera, viewport is relative to screensize.
func (c *CameraComponent) SetViewport(x, y, w, h float32) {
	c.Viewport = gm.Vec4{x, y, w, h}
}

// CalcViewport gives the viewport in screen dimensions
// TODO: consider `ScreenViewport` name
func (c *CameraComponent) CalcViewport(screenSize gm.Vec2) gm.Vec4 {
	return gm.Vec4{
		c.Viewport[0] * screenSize[0],
		c.Viewport[1] * screenSize[1],
		c.Viewport[2] * screenSize[0],
		c.Viewport[3] * screenSize[1],
	}
}
