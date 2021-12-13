package gorlet

import (
	"fmt"
	"reflect"

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

	// Save Restore SetProp stuff and all that.
	propStack propStack
}

// Create creates builds and prepares a guilet
func Create(fn BuildFunc) *Entity {
	root := &Entity{
		RectComponent: *gorgeui.NewRectComponent(),
	}
	b := Builder{
		root: &curEntity{elem: root},
	}

	fn(&b)
	return root
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
	if len(fns) == 0 {
		return
	}
	b.layout = gorgeui.MultiLayout(fns...)
}

// Root returns root guilet.
func (b *Builder) Root() *Entity {
	return b.root.elem
}

// Set a property for the next widget.
func (b *Builder) Set(k string, v interface{}) {
	b.propStack.cur().Set(k, v)
}

// SetProps a property for the next widget.
func (b *Builder) SetProps(p Props) {
	cur := b.propStack.cur()
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

// Observe adds a function to observe a property in the root Entity.
func (b Builder) Observe(k string, fn interface{}) {
	b.root.elem.observe(k, fn)
}

// UseProps set props, if an entity is passed it will set only on entity and return the entity.
// else it will set on builder.
/*func (b *Builder) UseProps(k string, e ...*Entity) *Entity {
	if len(e) == 0 {
		if b.propGroup == nil {
			return nil
		}
		p := b.propGroup.Select(k)

		b.SetProps(p)
		return nil
	}

	if b.propGroup == nil {
		return e[0]
	}
	b.setupProps(b.propGroup.Select(k), e[0])
	return e[0]
}

// DefineProps to be used in as groups.
func (b *Builder) DefineProps(g PropsGroup) {
	b.SetProps(g.Select(""))
	b.propGroup = g
}*/

// Create creates an Entity with builder properties
// NOTE: it does not add to the current container.
func (b *Builder) Create(fn BuildFunc) *Entity {
	w := Create(fn)
	b.setupProps(b.propStack.cur(), w)
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

// Begin creates and pushes an Entity onto stack it will save
// property state and restore on end.
func (b *Builder) Begin(fn BuildFunc) *Entity {
	w := b.Add(fn)
	b.push(w)
	b.propStack.Save()
	return w
}

// End pops the current guilet from the stack.
func (b *Builder) End() {
	b.propStack.Restore()
	b.pop()
}

func (b *Builder) setupProps(props Props, e *Entity) {
	for k, v := range props {
		k, v := k, v

		// If we don't have any observer, don't bother setting it.
		if _, ok := e.props[k]; !ok {
			continue
		}

		pk, ok := v.(ForwardProp)
		if !ok {
			e.Set(k, v)
			continue
		}
		b.Observe(pk.prop, func(v interface{}) {
			e.Set(k, v)
		})
		if pk.def != nil { // Set the default value
			e.Set(k, pk.def)
		}
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

func makePropFunc(k string, fn interface{}) func(interface{}) {
	if fn, ok := fn.(func(interface{})); ok {
		return fn
	}
	fnVal := reflect.ValueOf(fn)
	inTyp := fnVal.Type().In(0)

	return func(v interface{}) {
		arg := reflect.ValueOf(v)
		if aTyp := arg.Type(); aTyp != inTyp {
			if !arg.CanConvert(inTyp) {
				panic(fmt.Sprintf("Can't convert prop [%q] %v(%v) to %v", k, aTyp, v, inTyp))
			}
			arg = arg.Convert(inTyp)
		}
		// Type check somewhere
		fnVal.Call([]reflect.Value{arg})
	}
}
