package gorgeui

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
)

type gorger interface {
	Gorge() *gorge.Context
}

// UI root of UI if we get to use the UI tree
// like canvas in unity
// This is a entity
type UI struct {
	RectComponent
	ElementComponent

	CullMask      uint32
	DragThreshold float32 // int perhaps
	Order         int
	Camera        cameraEntity

	gorge *gorge.Context
}

// New returns a new UI.
func New(g gorger) *UI {
	e := &UI{
		RectComponent: RectIdent(),
		gorge:         g.Gorge(),
		CullMask:      gorge.MaskUI, // default
	}
	e.SetRect(0, 0, 100, 100)
	e.SetAnchor(0)
	// UI inverse
	e.SetScale(1, -1, 1)
	return e
}

// SetCullMask for the UI.
func (w *UI) SetCullMask(a uint32) {
	w.CullMask = a
}

// SetOrder for the UI.
func (w *UI) SetOrder(a int) {
	w.Order = a
}

// SetCamera sets this UI Camera.
func (w *UI) SetCamera(c cameraEntity) {
	w.Camera = c
}

// SetDragThreshold sets drag threshold for when pressing a pointer and allows
// to move v times before dragging.
func (w *UI) SetDragThreshold(v float32) {
	w.DragThreshold = v
}

// Rect return the rect for UI if camera is Ortho it will use camera as parent
// rect.
/*func (w *UI) Rect() m32.Vec4 {
	cam := w.Camera.Camera()
	if cam.ProjectionType != gorge.ProjectionOrtho {
		return w.RectComponent.Rect()
	}
	// ScreenSize
	vp := cam.CalcViewport(w.ScreenSize())
	aspectRatio := cam.AspectRatio
	if aspectRatio == 0 {
		// ss := gorge.ScreenSize()
		aspectRatio = vp[2] / vp[3]
	}

	halfV := cam.OrthoSize * .5
	halfH := cam.OrthoSize * .5 * aspectRatio

	camRect := m32.Vec4{-halfH, -halfV, halfH, halfV}

	return w.RelativeRect(camRect)
	// World rect
}*/

// Rect return the rect for UI if camera is Ortho it will use camera as parent
// rect.
func (w *UI) Rect() m32.Vec4 {
	// The ui must have a base rect or use the camera one
	// so it should be probably positioned from 0,0 of the screen
	cam := w.Camera.Camera()
	if cam.ProjectionType != gorge.ProjectionOrtho {
		return w.RectComponent.Rect()
	}
	// ScreenSize
	vp := cam.CalcViewport(w.ScreenSize())
	aspectRatio := cam.AspectRatio
	if aspectRatio == 0 {
		// ss := gorge.ScreenSize()
		aspectRatio = vp[2] / vp[3]
	}

	// halfV := cam.OrthoSize * .5
	// halfH := cam.OrthoSize * .5 * aspectRatio

	// camRect := m32.Vec4{-halfH, -halfV, halfH, halfV}
	camRect := m32.Vec4{0, 0, cam.OrthoSize * aspectRatio, cam.OrthoSize}

	// Default maybe
	w.Pivot[0] = .5
	w.Pivot[1] = .5
	w.Dim[0] = cam.OrthoSize * aspectRatio
	w.Dim[1] = cam.OrthoSize

	// camRect := m32.Vec4{0, 0, cam.OrthoSize, cam.OrthoSize * aspectRatio}
	// camRect := m32.Vec4{0, 0, halfH, halfV}
	// Should not change these since user can write these?
	// w.Position[0] = -halfH
	// w.Position[1] = halfV

	// camRect := m32.Vec4{0, 0, cam.OrthoSize, cam.OrthoSize * aspectRatio}
	// log.Println("returning camRect:", camRect)
	// w.Dim[0] = cam.OrthoSize
	// w.Dim[1] = cam.OrthoSize * aspectRatio

	// return camRect
	return w.RelativeRect(camRect)
	// World rect
}

// ScreenSize returns the screensize.
func (w *UI) ScreenSize() m32.Vec2 { return w.gorge.ScreenSize() }

type parenter interface {
	SetParent(gorge.Transformer)
	Parent() gorge.Transformer
}

// Add alias to add entities to gorge.
func (w *UI) Add(ents ...gorge.Entity) {
	for _, e := range ents {
		if r, ok := e.(parenter); ok {
			if r.Parent() == nil {
				r.SetParent(w)
			}
		}
	}
	w.gorge.Add(ents...)
}

// Remove alias to remove entities from gorge.
func (w *UI) Remove(ents ...gorge.Entity) { w.gorge.Remove(ents...) }

type uiSorter []*UI

func (s uiSorter) Len() int { return len(s) }
func (s uiSorter) Less(i, j int) bool {
	return s[i].Order < s[j].Order
}

func (s uiSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
