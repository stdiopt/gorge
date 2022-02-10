package gorgeui

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
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

	// Also add focus and stuff
	gorge *gorge.Context
}

func (u *UI) String() string {
	return "UI"
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

func Add(g *gorge.Context) *UI {
	defCam := &defaultCamera{
		TransformComponent: *gorge.NewTransformComponent(),
		CameraComponent: gorge.CameraComponent{
			Name:           "ui camera",
			ProjectionType: gorge.ProjectionOrtho,
			CullMask:       gorge.CullMaskUI | gorge.CullMaskUIDebug,
			OrthoSize:      100,
			Near:           -100,
			Far:            100,
			ClearFlag:      gorge.ClearDepthOnly,
			Order:          3000,
			Viewport:       gm.Vec4{0, 0, 1, 1},
			ClearColor:     gm.Vec3{.3, .3, .3},
		},
	}
	g.Add(defCam)
	ui := New()
	ui.SetCamera(defCam)

	g.Add(ui)

	return ui
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
	if w.Camera == c {
		return
	}
	if c, ok := w.Camera.(*defaultCamera); ok && w.gorge != nil {
		w.gorge.Remove(c)
	}
	w.Camera = c
}

// SetDragThreshold sets drag threshold for when pressing a pointer and allows
// to move v times before dragging.
func (w *UI) SetDragThreshold(v float32) {
	w.DragThreshold = v
}

// Rect return the rect for UI if camera is Ortho it will use camera as parent
// rect.
func (w *UI) Rect() gm.Vec4 {
	// The ui must have a base rect or use the camera one
	// so it should be probably positioned from 0,0 of the screen
	cam := w.Camera.Camera()
	if cam.ProjectionType != gorge.ProjectionOrtho {
		return gm.Vec4{0, 0, w.Size[0], w.Size[1]}
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
	return gm.Vec4{
		-halfH + w.Position[0], -halfV + w.Position[1],
		halfH + w.Size[0], halfV + w.Size[1],
	}
}

func (w *UI) CalcSize() gm.Vec2 {
	cam := w.Camera.Camera()
	if cam.ProjectionType != gorge.ProjectionOrtho {
		return w.Size
	}
	// ScreenSize
	vp := cam.CalcViewport(w.ScreenSize())
	aspectRatio := cam.AspectRatio
	if aspectRatio == 0 {
		// ss := gorge.ScreenSize()
		aspectRatio = vp[2] / vp[3]
	}
	return gm.Vec2{
		cam.OrthoSize * aspectRatio / 2,
		cam.OrthoSize,
	}
}

// ScreenSize returns the screensize.
func (w *UI) ScreenSize() gm.Vec2 {
	return w.gorge.ScreenSize()
}

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
	w.GAdd(ents...)
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
	w.GRemove(ents...)
}

// GAdd adds to gorge if attached.
func (w *UI) GAdd(ents ...gorge.Entity) {
	if w.Attached && w.gorge != nil {
		w.gorge.Add(ents...)
	}
	w.reorder()
}

// GRemove removes from gorge if attached.
func (w *UI) GRemove(ents ...gorge.Entity) {
	if w.Attached && w.gorge != nil {
		w.gorge.Remove(ents...)
	}
	w.reorder()
}

func (w *UI) update(dt float32) {
}

func (w *UI) reorder() {
	type renderable interface {
		Renderable() *gorge.RenderableComponent
	}

	order := 0
	count := 0
	var walk func(e gorge.Entity)
	walk = func(e gorge.Entity) {
		count++
		if e, ok := e.(renderable); ok {
			e.Renderable().SetOrder(order)
			order++
		}
		if e, ok := e.(gorge.EntityContainer); ok {
			for _, ee := range e.GetEntities() {
				walk(ee)
			}
		}
	}

	for _, e := range w.entities {
		walk(e)
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
