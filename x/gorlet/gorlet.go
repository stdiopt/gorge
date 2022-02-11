package gorlet

import (
	"fmt"
	"log"
	"runtime"
	"sync/atomic"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/math/ray"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// Debug prints reference debug
var Debug = false

// DebugCounter debug release stuff
var debugCounter int64

// DebugCounter returns the number of gorlets instantiated.
func DebugCounter() int64 {
	return atomic.LoadInt64(&debugCounter)
}

type Entity interface {
	gorgeui.Entity
	event.Buser

	ID() string
	SetID(string)
	Find(string) Entity

	GetEntities() []gorge.Entity

	add(...gorge.Entity)
	remove(...gorge.Entity)
}

type renderabler interface {
	Renderable() *gorge.RenderableComponent
}

type widgeter interface {
	CalcMax() gm.Vec2
	gorgeui.RectTransformer
}

type parenter interface {
	gorge.ParentGetter
	gorge.ParentSetter
}

type addRemover interface {
	add(...gorge.Entity)
	remove(...gorge.Entity)
}

type masker interface {
	setMaskDepth(n int)
}

type entityRemover interface {
	remove(...gorge.Entity)
}

type gcref struct{ n int }

type entityConstraint[T any] interface {
	*T
	Entity
}

type Widget[T any, Tp entityConstraint[T]] struct {
	id string
	gorgeui.RectComponent
	gorgeui.ElementComponent
	base Tp

	clientArea Entity
	masked     bool
	layout     Layouter

	gorge.Container

	gcref *gcref
}

func (w *Widget[T, Tp]) String() string {
	return fmt.Sprintf("%T %s", w.base, w.id)
}

func (b *Widget[T, Tp]) setBase(b2 Widget[T, Tp]) {
	*b = b2
}

type baseSetter[T any, Tp entityConstraint[T]] interface {
	setBase(Widget[T, Tp])
}

func Build[T any, Tp entityConstraint[T]](s Tp) Tp {
	ref := &gcref{1}

	obj := fmt.Sprintf("%T", s)
	if Debug {
		log.Println("Allocating:", obj)
	}
	runtime.SetFinalizer(ref, func(any) {
		if Debug {
			log.Printf("Finalizing Gorlet OBJ: %v", obj)
		}
		atomic.AddInt64(&debugCounter, -1)
	})

	atomic.AddInt64(&debugCounter, 1)

	bw := Widget[T, Tp]{
		gcref: ref,
		RectComponent: gorgeui.RectComponent{
			Anchor:   gm.Vec4{0, 0, 1, 1},
			Size:     gm.Vec2{0, 0},
			Scale:    gm.Vec3{1, 1, 1},
			Rotation: gm.QIdent(),
		},
		ElementComponent: gorgeui.ElementComponent{},
		base:             s,
	}

	any(s).(baseSetter[T, Tp]).setBase(bw)

	if builder, ok := any(s).(interface{ Build(b *B) }); ok {
		b := &B{root: &curEntity{entity: s}}
		b.Do(builder.Build)
		entityUpdate(s)
	}
	return s
}

func entityUpdate(ent gorge.Entity) {
	if v, ok := ent.(event.Handler); ok {
		v.HandleEvent(gorgeui.EventUpdate(0))
	}
	if v, ok := ent.(event.Buser); ok {
		event.Trigger(v, gorgeui.EventUpdate(0))
	}
	if v, ok := ent.(gorge.EntityContainer); ok {
		for _, c := range v.GetEntities() {
			entityUpdate(c)
		}
	}
}

func (w *Widget[T, Tp]) ID() string {
	return w.id
}

func (w *Widget[T, Tp]) SetID(id string) {
	w.id = id
}

func (w *Widget[T, Tp]) Attached(e gorgeui.Entity) {
	// Sounds bad
	for _, c := range w.GetEntities() {
		if r, ok := c.(renderabler); ok {
			if ui := gorgeui.RootUI(e); ui != nil {
				r.Renderable().CullMask = ui.CullMask
			}
		}
	}
}

func (w *Widget[T, Tp]) GetClientArea() Entity {
	return w.clientArea
}

func (w *Widget[T, Tp]) SetClientArea(a Entity) Tp {
	w.clientArea = a
	return w.base
}

func (w *Widget[T, Tp]) SetDragEvents(b bool) Tp {
	w.DragEvents = b
	return w.base
}

func (w *Widget[T, Tp]) SetMargin(vs ...float32) Tp {
	w.RectComponent.SetMargin(vs...)
	return w.base
}

func (w *Widget[T, Tp]) SetBorder(vs ...float32) Tp {
	w.RectComponent.SetBorder(vs...)
	return w.base
}

func (w *Widget[T, Tp]) SetDisableRaycast(b bool) Tp {
	w.ElementComponent.SetDisableRaycast(b)
	return w.base
}

func (w *Widget[T, Tp]) SetAnchor(vs ...float32) Tp {
	w.RectComponent.SetAnchor(vs...)
	return w.base
}

func (w *Widget[T, Tp]) SetRect(vs ...float32) Tp {
	w.RectComponent.SetRect(vs...)
	return w.base
}

func (w *Widget[T, Tp]) FillParent() Tp {
	w.RectComponent.SetAnchor(0, 0, 1, 1)
	w.RectComponent.SetSize(0)
	return w.base
}

func (w *Widget[T, Tp]) SetPivot(vs ...float32) Tp {
	w.RectComponent.SetPivot(vs...)
	return w.base
}

func (w *Widget[T, Tp]) SetSize(v ...float32) Tp {
	w.RectComponent.SetSize(v...)
	return w.base
}

func (w *Widget[T, Tp]) SetWidth(v float32) Tp {
	w.RectComponent.Size[0] = v
	return w.base
}

func (w *Widget[T, Tp]) SetHeight(v float32) Tp {
	w.RectComponent.Size[1] = v
	return w.base
}

func (w *Widget[T, Tp]) SetPosition(x, y, z float32) Tp {
	w.Position = gm.Vec3{x, y, z}
	return w.base
}

// IntersectFromScreen intersects the entity rect from screen coordinates.
func (w *Widget[T, Tp]) IntersectFromScreen(pos gm.Vec2) ray.Result {
	sz := w.ContentSize()
	m := w.Mat4()
	v0 := m.MulV4(gm.Vec4{0, 0, 0, 1}).Vec3()     // 0
	v1 := m.MulV4(gm.Vec4{sz[0], 0, 0, 1}).Vec3() // right
	v2 := m.MulV4(gm.Vec4{0, sz[1], 0, 1}).Vec3() // up)

	ui := gorgeui.RootUI(w.base)
	r := ray.FromScreen(ui.ScreenSize(), ui.Camera, pos)
	return ray.IntersectRect(r, v0, v1, v2)
}

// CalcMax calculates children bounds and positions and return min max
// Calc maximum of the children
func (w *Widget[T, Tp]) CalcMax() gm.Vec2 { // CalcMax
	sz := w.ContentSize()
	sz[0] += w.Margin[0] + w.Margin[2]
	sz[1] += w.Margin[1] + w.Margin[3]
	if w.masked {
		return sz
	}

	for _, c := range w.GetEntities() {
		e, ok := c.(widgeter)
		if !ok {
			continue
		}

		b := e.CalcMax()
		p := e.RectTransform().Position.XY()
		p = p.Sub(e.RectTransform().Pivot.MulVec2(b))
		sz[0] = gm.Max(sz[0], p[0]+b[0])
		sz[1] = gm.Max(sz[1], p[1]+b[1])
	}
	return sz
}

// IsMasked returns true if the widget is masked
// so it won't receive further events outside of its bounds.
func (w *Widget[T, Tp]) IsMasked() bool {
	return w.masked
}

func (w *Widget[T, Tp]) Add(cs ...gorge.Entity) Tp {
	w.add(cs...)
	return w.base
}

func (w *Widget[T, Tp]) Remove(cs ...gorge.Entity) Tp {
	w.remove(cs...)
	return w.base
}

/*
func (w *Widget[T, Tp]) SetLayout(l ...Layouter) Tp {
	w.setLayout(l...)
	return w.base
}
*/

// XXX:
// we could use layouter func here instead of gorgeui
// this way we could just run through *Entity since this is
// an higher level entity.
//
// Other way would be use the eventBus on the Create func.
// Since then we wouldn't have HandleEvent exposed

// HandleEvent handles events.
func (w *Widget[T, Tp]) HandleEvent(evt event.Event) {
	switch evt.(type) {
	case gorgeui.EventUpdate:
		if w.layout != nil {
			w.layout.Layout(w.base)
		}
	default:
	}
}

func (w *Widget[T, Tp]) add(cs ...gorge.Entity) {
	if w.clientArea != nil {
		w.clientArea.add(cs...)
		return
	}
	for _, c := range cs {
		if p, ok := c.(parenter); ok {
			if r, ok := p.Parent().(entityRemover); ok {
				r.remove(c)
			}
			p.SetParent(w.base)
		}
	}
	w.Container.Add(cs...)

	if ui := gorgeui.RootUI(w.base); ui != nil {
		ui.GAdd(cs...)
	}
}

func (w *Widget[T, Tp]) remove(cs ...gorge.Entity) {
	if w.clientArea != nil {
		w.clientArea.remove(cs...)
		return
	}

	for _, c := range cs {
		if p, ok := c.(parenter); ok {
			p.SetParent(nil)
		}
		if c == w.clientArea {
			w.clientArea = nil
		}
	}
	w.Container.Remove(cs...)
	if ui := gorgeui.RootUI(w.base); ui != nil {
		ui.GRemove(cs...)
	}
}

// Find returns the first entity in the graph that matches the predicate.
func (w *Widget[T, Tp]) Find(id string) Entity {
	if w.ID() == id {
		return w.base
	}
	children := w.GetEntities()
	for _, c := range children {
		e, ok := c.(Entity)
		if !ok {
			continue
		}
		if r := e.Find(id); r != nil {
			return r
		}
	}
	return nil
}

/*
func (w *Widget[T, Tp]) setLayout(l ...Layouter) {
	if w.clientArea != nil {
		w.clientArea.setLayout(l...)
	}
	w.layout = MultiLayout(l...)
}
*/

func (w *Widget[T, Tp]) setClientArea(a Entity) {
	w.clientArea = a
}

func (w *Widget[T, Tp]) setMaskDepth(n int) {
	for _, c := range w.GetEntities() {
		if s, ok := c.(masker); ok {
			s.setMaskDepth(n)
		}
	}
}
