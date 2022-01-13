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

	CullMask      gorge.CullMaskFlags
	DragThreshold float32 // int perhaps
	Order         int
	Camera        cameraEntity

	entities []gorge.Entity

	gorge *gorge.Context
}

// New returns a new UI.
func New() *UI {
	e := &UI{
		RectComponent: RectIdent(),
		// gorge:         g.Gorge(),
		CullMask: gorge.CullMaskUI, // default
	}
	e.SetAnchor(0)
	// e.SetPivot(.5)
	// UI inverse, might be better elsewhere
	e.SetScale(1, -1, 1)
	return e
}

// SetCullMask for the UI.
func (w *UI) SetCullMask(a gorge.CullMaskFlags) {
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
	// The ui must have a base rect or use the camera one
	// so it should be probably positioned from 0,0 of the screen
	cam := w.Camera.Camera()
	if cam.ProjectionType != gorge.ProjectionOrtho {
		return m32.Vec4{0, 0, w.Dim[0], w.Dim[1]}
	}
	// ScreenSize
	vp := cam.CalcViewport(w.ScreenSize())
	aspectRatio := cam.AspectRatio
	if aspectRatio == 0 {
		// ss := gorge.ScreenSize()
		aspectRatio = vp[2] / vp[3]
	}
	halfH := cam.OrthoSize * aspectRatio / 2
	halfV := cam.OrthoSize / 2
	return m32.Vec4{
		-halfH + w.Position[0], -halfV + w.Position[1],
		halfH + w.Dim[0], halfV + w.Dim[1],
	}
}

func (w *UI) CalcDim() m32.Vec2 {
	cam := w.Camera.Camera()
	if cam.ProjectionType != gorge.ProjectionOrtho {
		return w.Dim
	}
	// ScreenSize
	vp := cam.CalcViewport(w.ScreenSize())
	aspectRatio := cam.AspectRatio
	if aspectRatio == 0 {
		// ss := gorge.ScreenSize()
		aspectRatio = vp[2] / vp[3]
	}
	return m32.Vec2{
		cam.OrthoSize * aspectRatio / 2,
		cam.OrthoSize,
	}
}

// ScreenSize returns the screensize.
func (w *UI) ScreenSize() m32.Vec2 { return w.gorge.ScreenSize() }

func (w *UI) GetEntities() []gorge.Entity {
	return w.entities
}

type parenter interface {
	SetParent(gorge.Matrixer)
	Parent() gorge.Matrixer
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
	w.entities = append(w.entities, ents...)
	if w.Attached && w.gorge != nil {
		w.gorge.Add(ents...)
	}
}

// Remove alias to remove entities from gorge.
func (w *UI) Remove(ents ...gorge.Entity) {
	for _, e := range ents {
		if r, ok := e.(parenter); ok {
			if r.Parent() == w {
				r.SetParent(nil)
			}
		}
		for i, ee := range w.entities {
			if ee == e {
				t := w.entities
				w.entities = append(w.entities[:i], w.entities[i+1:]...)
				t[len(t)-1] = nil
				break
			}
		}
	}
	if w.Attached && w.gorge != nil {
		w.gorge.Remove(ents...)
	}
}

type uiSorter []*UI

func (s uiSorter) Len() int { return len(s) }
func (s uiSorter) Less(i, j int) bool {
	return s[i].Order < s[j].Order
}

func (s uiSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
