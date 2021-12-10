package gorlet

import (
	"fmt"
	"log"
	"reflect"

	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// ForwardProp to be used to forward properties.
type ForwardProp struct {
	prop string
	def  interface{}
}

// PlacementFunc will be used in a container and will define clients rect.
type PlacementFunc func(w *Entity) // OnAdd in the Entity

type curEntity struct {
	placement PlacementFunc
	elem      *Entity
}

// BuildFunc to build a guilet
type BuildFunc func(b *Builder)

// Create creates builds and prepares a guilet
func Create(fn BuildFunc) *Entity {
	root := &Entity{
		RectComponent: *gorgeui.NewRectComponent(),
	}
	b := Builder{root: &curEntity{
		elem: root,
	}}

	fn(&b)
	return root
}

// AddMode builder add mode.
type AddMode int

// AddMode constants.
const (
	ChildrenAdd = AddMode(iota)
	ElementAdd
)

// Builder used to build a guilet.
type Builder struct {
	placement PlacementFunc
	layout    gorgeui.LayoutFunc

	mode AddMode

	stack []*curEntity
	root  *curEntity
	props propStack
}

// SetAddMode set Entity add mode.
func (b *Builder) SetAddMode(mode AddMode) {
	b.mode = mode
}

// Placement sets the placement func.
func (b *Builder) Placement(fn PlacementFunc) {
	b.placement = fn
}

// Layout set next widget layout.
func (b *Builder) Layout(fns ...gorgeui.LayoutFunc) {
	b.layout = gorgeui.MultiLayout(fns...)
}

// Root returns root guilet.
func (b *Builder) Root() *Entity {
	return b.root.elem
}

// Set a property for the next widget.
func (b *Builder) Set(k string, v interface{}) {
	b.props.cur().Set(k, v)
}

// SetProps a property for the next widget.
func (b *Builder) SetProps(p Props) {
	cur := b.props.cur()
	for k, v := range p {
		cur.Set(k, v)
	}
}

// Prop returns a property forwarded by k with optional default value.
func (b *Builder) Prop(k string, v ...interface{}) ForwardProp {
	var def interface{}
	if len(v) > 0 {
		def = v[0]
	}
	return ForwardProp{prop: k, def: def}
}

// Observe adds a function to observe a property.
func (b Builder) Observe(k string, fn interface{}) {
	b.root.elem.observe(k, fn)
}

// Create creates an Entity with builder properties
// NOTE: it does not add to the current container.
func (b *Builder) Create(fn BuildFunc) *Entity {
	w := Create(fn)
	b.setupProps(w)
	return w
}

// Add an Entity to the current container.
func (b *Builder) Add(fn BuildFunc) *Entity {
	e := b.Create(fn)

	cur := b.cur()
	switch b.mode {
	case ChildrenAdd:
		if cur.placement != nil {
			cur.placement(e)
		}
		cur.elem.Add(e)
	case ElementAdd:
		e.SetPivot(.5)
		e.SetAnchor(0, 0, 1, 1)
		e.SetRect(0)
		cur.elem.AddElement(e)
	}
	return e
}

// AddEntity adds a prebuilt entity.
func (b *Builder) AddEntity(e *Entity) *Entity {
	cur := b.cur()
	switch b.mode {
	case ChildrenAdd:
		if cur.placement != nil {
			cur.placement(e)
		}
		cur.elem.Add(e)
	case ElementAdd:
		e.SetPivot(.5)
		e.SetAnchor(0, 0, 1, 1)
		e.SetRect(0)
		cur.elem.AddElement(e)
	}
	return e
}

// Begin creates and pushes a guilet onto stack.
func (b *Builder) Begin(fn BuildFunc) *Entity {
	w := b.Add(fn)
	b.push(w)
	b.props.Save()
	return w
}

// End pops the current guilet from the stack.
func (b *Builder) End() {
	b.props.Restore()
	b.pop()
}

func (b *Builder) setupProps(w *Entity) {
	p := b.props.cur()
	for k, v := range p {
		pk, ok := v.(ForwardProp)
		if !ok {
			w.Set(k, v)
			continue
		}
		k := k // shadow
		b.Observe(pk.prop, func(v interface{}) {
			w.Set(k, v)
		})
		if pk.def != nil { // Set the default value
			w.Set(k, pk.def)
		}
		continue
	}
}

func (b *Builder) push(g *Entity) {
	cur := curEntity{elem: g}
	cur.placement = b.placement
	g.LayoutFunc = b.layout

	b.placement = nil
	b.layout = nil

	b.stack = append(b.stack, &cur)
}

func (b *Builder) pop() {
	if len(b.stack) == 0 {
		return
	}
	b.stack = b.stack[:len(b.stack)-1]
}

func (b *Builder) cur() *curEntity {
	if len(b.stack) == 0 {
		return b.root
	}
	return b.stack[len(b.stack)-1]
}

func makePropFunc(fn interface{}) func(interface{}) {
	if fn, ok := fn.(func(interface{})); ok {
		return fn
	}
	fnVal := reflect.ValueOf(fn)
	inTyp := fnVal.Type().In(0)

	return func(v interface{}) {
		arg := reflect.ValueOf(v)
		if aTyp := arg.Type(); aTyp != inTyp {
			if !arg.CanConvert(inTyp) {
				panic(fmt.Sprintf("Can't convert %v to %v", aTyp, inTyp))
			}
			arg = arg.Convert(inTyp)
		}
		// Type check somewhere
		fnVal.Call([]reflect.Value{arg})
	}
}

// Vertical placement
func Vertical(spacing m32.Vec4, dim m32.Vec2) PlacementFunc {
	var pos m32.Vec2
	return func(w *Entity) {
		w.SetAnchor(0, 0, 1, 0)
		w.SetRect(spacing[0], spacing[1]+pos[1], spacing[2], dim[1])
		w.SetPivot(0)
		pos[1] += dim[1] + spacing[3]
		log.Println("Setting:", w.Rect())
	}
}
