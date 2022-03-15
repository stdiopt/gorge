package gorlet

import (
	"fmt"
	"log"
	"runtime"
	"strings"
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

type entityConstraint interface {
	Entity
}

type Widget[T entityConstraint] struct {
	id   string
	name string
	gorgeui.RectComponent
	gorgeui.ElementComponent
	base T

	clientArea Entity
	masked     bool
	layout     Layouter

	gorge.Container

	gcref *gcref
}

func (w *Widget[T]) String() string {
	b := &strings.Builder{}

	fmt.Fprintf(b, "%T", w.base)
	if w.name != "" {
		fmt.Fprintf(b, "[%q]", w.name)
	}
	if w.id != "" {
		fmt.Fprintf(b, "#%q", w.id)
	}
	return b.String()
}

func (b *Widget[T]) setBase(b2 Widget[T]) {
	*b = b2
}

type baseSetter[T entityConstraint] interface {
	setBase(Widget[T])
}

func Build[T entityConstraint](s T) T {
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

	bw := Widget[T]{
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

	any(s).(baseSetter[T]).setBase(bw)

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

func (w *Widget[T]) ID() string {
	return w.id
}

func (w *Widget[T]) SetID(id string) {
	w.id = id
}

func (w *Widget[T]) SetName(name string) T {
	w.name = name
	return w.base
}

func (w *Widget[T]) Attached(e gorgeui.Entity) {
	// Sounds bad
	for _, c := range w.GetEntities() {
		if r, ok := c.(renderabler); ok {
			if ui := gorgeui.RootUI(e); ui != nil {
				r.Renderable().CullMask = ui.CullMask
			}
		}
	}
}

func (w *Widget[T]) GetClientArea() Entity {
	return w.clientArea
}

func (w *Widget[T]) SetClientArea(a Entity) T {
	w.clientArea = a
	return w.base
}

func (w *Widget[T]) SetDragEvents(b bool) T {
	w.DragEvents = b
	return w.base
}

func (w *Widget[T]) SetMargin(vs ...float32) T {
	w.RectComponent.SetMargin(vs...)
	return w.base
}

func (w *Widget[T]) SetBorder(vs ...float32) T {
	w.RectComponent.SetBorder(vs...)
	return w.base
}

func (w *Widget[T]) SetDisableRaycast(b bool) T {
	w.ElementComponent.SetDisableRaycast(b)
	return w.base
}

func (w *Widget[T]) SetAnchor(vs ...float32) T {
	w.RectComponent.SetAnchor(vs...)
	return w.base
}

func (w *Widget[T]) SetRect(vs ...float32) T {
	w.RectComponent.SetRect(vs...)
	return w.base
}

func (w *Widget[T]) FillParent() T {
	w.RectComponent.SetAnchor(0, 0, 1, 1)
	w.RectComponent.SetSize(0)
	return w.base
}

func (w *Widget[T]) SetPivot(vs ...float32) T {
	w.RectComponent.SetPivot(vs...)
	return w.base
}

func (w *Widget[T]) SetSize(v ...float32) T {
	w.RectComponent.SetSize(v...)
	return w.base
}

func (w *Widget[T]) SetWidth(v float32) T {
	w.Size[0] = v
	return w.base
}

func (w *Widget[T]) SetHeight(v float32) T {
	w.Size[1] = v
	return w.base
}

func (w *Widget[T]) SetPosition(x, y, z float32) T {
	w.Position = gm.Vec3{x, y, z}
	return w.base
}

// IntersectFromScreen intersects the entity rect from screen coordinates.
func (w *Widget[T]) IntersectFromScreen(pos gm.Vec2) ray.Result {
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
func (w *Widget[T]) CalcMax() gm.Vec2 { // CalcMax
	r := w.Rect()
	r[2] += w.Margin[0] + w.Margin[2]
	r[3] += w.Margin[1] + w.Margin[3]
	if w.masked {
		return r.ZW()
	}

	for _, c := range w.GetEntities() {
		e, ok := c.(widgeter)
		if !ok {
			continue
		}
		rt := e.RectTransform()
		p := gm.Vec2{
			rt.Position[0] + r[2]*rt.Anchor[0],
			rt.Position[1] + r[3]*rt.Anchor[1],
		}

		b := e.CalcMax()

		p = p.Sub(e.RectTransform().Pivot.MulVec2(b))
		r[2] = gm.Max(r[2], p[0]+b[0])
		r[3] = gm.Max(r[3], p[1]+b[1])
	}
	return r.ZW()
}

// IsMasked returns true if the widget is masked
// so it won't receive further events outside of its bounds.
func (w *Widget[T]) IsMasked() bool {
	return w.masked
}

func (w *Widget[T]) Add(cs ...gorge.Entity) T {
	w.add(cs...)
	return w.base
}

func (w *Widget[T]) Remove(cs ...gorge.Entity) T {
	w.remove(cs...)
	return w.base
}

// XXX:
// we could use layouter func here instead of gorgeui
// this way we could just run through *Entity since this is
// an higher level entity.
//
// Other way would be use the eventBus on the Create func.
// Since then we wouldn't have HandleEvent exposed

// HandleEvent handles events.
func (w *Widget[T]) HandleEvent(evt event.Event) {
	switch evt.(type) {
	case gorgeui.EventUpdate:
		if w.layout != nil {
			w.layout.Layout(w.base)
		}
	default:
	}
}

func (w *Widget[T]) add(cs ...gorge.Entity) {
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

func (w *Widget[T]) remove(cs ...gorge.Entity) {
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
func (w *Widget[T]) Find(id string) Entity {
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

func (w *Widget[T]) setClientArea(a Entity) {
	w.clientArea = a
}

func (w *Widget[T]) setMaskDepth(n int) {
	for _, c := range w.GetEntities() {
		if s, ok := c.(masker); ok {
			s.setMaskDepth(n)
		}
	}
}
