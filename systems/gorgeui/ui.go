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

	gorge *gorge.Context

	CullMask      uint32
	DragThreshold float32 // int perhaps
	Order         int
	Camera        cameraEntity
}

// New returns a new UI.
func New(g gorger) *UI {
	e := &UI{
		gorge:         g.Gorge(),
		RectComponent: RectIdent(),
		CullMask:      gorge.MaskUI, // default
	}
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
func (w *UI) Rect() m32.Vec4 {
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
}

// ScreenSize returns the screensize.
func (w *UI) ScreenSize() m32.Vec2 { return w.gorge.ScreenSize() }

// Add alias to add entities to gorge.
func (w *UI) Add(ents ...gorge.Entity) {
	for _, e := range ents {
		if r, ok := e.(rectTransformer); ok {
			if r.RectTransform().Parent() == nil {
				r.RectTransform().SetParent(w)
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
