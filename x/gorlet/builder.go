package gorlet

import (
	"fmt"
	"reflect"
)

// ForwardProp to be used to forward properties.
type ForwardProp struct {
	prop string
	def  interface{}
}

type curEntity struct {
	// will be called on each child
	placement EntityFunc
	entity    *Entity
}

// Func to build a guilet
type Func func(b *Builder)

type nextData struct {
	// Placement will be set on
	// placement EntityFunc

	layout     Layouter
	margin     []float32
	rect       []float32
	anchor     []float32
	pivot      []float32
	dragEvents *bool

	props Props
}

// Builder used to build a guilet.
type Builder struct {
	next nextData

	onAddFn func(e *Entity)

	stack []*curEntity
	root  *curEntity

	// Save Restore SetProp stuff and all that.
	propStack propStack
}

// Root returns root guilet.
func (b *Builder) Root() *Entity {
	return b.root.entity
}

// ClientArea sets the root entity client area, when adding entities using Add
// those will be added to the current container
// calling this twice will override the previous call.
func (b Builder) ClientArea() {
	cur := b.cur()
	if cur == b.root {
		return
	}
	b.root.entity.SetClientArea(cur.entity)
}

///////////////////////////////////////////////////////////////////////////////
// Properties related to the next entity
///////////////////////////////////////////////////////////////////////////////

// UsePlacement sets the placement func.
func (b *Builder) SetPlacement(fn EntityFunc) {
	b.cur().placement = fn
}

// UseLayout set next widget layout.
func (b *Builder) UseLayout(fns ...Layouter) {
	if len(fns) == 0 {
		return
	}
	b.next.layout = MultiLayout(fns...)
}

// UseDimRect uses dimension rect sets anchor to 0
func (b *Builder) UseDimRect(v ...float32) {
	b.UseAnchor(0)
	b.UsePivot(0)
	b.UseRect(v...)
}

// UseRelRect uses relative from parent rect.
func (b *Builder) UseRelRect(v ...float32) {
	b.UseAnchor(0, 0, 1, 1)
	b.UsePivot(0)
	b.UseRect(v...)
}

// UseMargin sets the next entity padding.
func (b *Builder) UseMargin(v ...float32) {
	b.next.margin = v
}

// UseRect sets next Entity Rect.
func (b *Builder) UseRect(v ...float32) {
	b.next.rect = v
}

// UseAnchor sets next Entity Anchor.
func (b *Builder) UseAnchor(v ...float32) {
	b.next.anchor = v
}

// UsePivot sets next Entity Pivot.
func (b *Builder) UsePivot(v ...float32) {
	b.next.pivot = v
}

func (b *Builder) UseDragEvents(v bool) {
	b.next.dragEvents = &v
}

// Use a property for the next widget.
func (b *Builder) Use(k string, v interface{}) {
	if b.next.props == nil {
		b.next.props = Props{}
	}
	b.next.props.Set(k, v)
}

// UseProps a property for the next widget.
func (b *Builder) UseProps(p Props) {
	if b.next.props == nil {
		b.next.props = Props{}
	}
	for k, v := range p {
		b.next.props.Set(k, v)
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

// ForwardProps will forward any entity props in this builder to the entity
// using a prefix
func (b *Builder) ForwardProps(pre string, e *Entity) {
	for k, v := range e.observers {
		key := k
		if pre != "" {
			key = pre + "." + k
		}
		for _, fn := range v {
			b.Observe(key, fn)
		}
	}
}

// Observe adds a function to observe a property in the root Entity.
func (b Builder) Observe(k string, fn func(interface{})) {
	b.root.entity.Observe(k, fn)
}

// Push will set the prop to any added entity.
func (b *Builder) Push(k string, v interface{}) {
	b.propStack.cur().Set(k, v)
}

// PushProps will set the props to any added entity.
func (b *Builder) PushProps(p Props) {
	cur := b.propStack.cur()
	for k, v := range p {
		cur.Set(k, v)
	}
}

// Save save props state onto stack.
func (b *Builder) Save() {
	b.propStack.Save()
}

// Restore restores props from stack.
func (b *Builder) Restore() {
	b.propStack.Restore()
}

// Create creates an Entity with builder properties
// NOTE: it does not add to the current container.
func (b *Builder) Create(fn Func) *Entity {
	e := Create(fn)

	if pfn := b.cur().placement; pfn != nil {
		pfn(e)
	}
	// Different thing
	// e.OnAdd(b.next.placement)

	if b.next.dragEvents != nil {
		e.SetDragEvents(*b.next.dragEvents)
	}
	if b.next.margin != nil {
		e.SetMargin(b.next.margin...)
	}
	if b.next.layout != nil {
		e.SetLayout(b.next.layout)
	}
	if b.next.rect != nil {
		e.SetRect(b.next.rect...)
	}
	if b.next.anchor != nil {
		e.SetAnchor(b.next.anchor...)
	}
	if b.next.pivot != nil {
		e.SetPivot(b.next.pivot...)
	}

	// Merge props
	props := b.propStack.cur().Merge(b.next.props)
	b.setupProps(e, props)

	b.next = nextData{}

	return e
}

// Add creates and add an Entity to the current container.
func (b *Builder) Add(fn Func) *Entity {
	return b.AddEntity(b.Create(fn))
}

// AddEntity adds a prebuilt entity.
func (b *Builder) AddEntity(e *Entity) *Entity {
	cur := b.cur()
	cur.entity.Add(e)
	if b.onAddFn != nil {
		b.onAddFn(e)
	}

	return e
}

// SetRoot will set the root container.
func (b *Builder) SetRoot(fn Func) *Entity {
	if len(b.stack) > 0 {
		panic("Builder.Start() called while already in a container")
	}
	if len(b.root.entity.observers) > 0 {
		panic("Builder.Start() called while root already has observers")
	}
	e := b.Create(fn)

	b.root.entity = e
	return e
}

// Begin creates and pushes an Entity onto stack it will save
// property state and restore on end.
func (b *Builder) Begin(fn Func) *Entity {
	return b.BeginEntity(b.Create(fn))
}

func (b *Builder) BeginEntity(e *Entity) *Entity {
	b.AddEntity(e)
	b.push(e)
	b.propStack.Save()
	return e
}

// End pops the current guilet from the stack.
func (b *Builder) End() {
	b.propStack.Restore()
	b.pop()
}

func (b *Builder) setupProps(e *Entity, props Props) {
	for k, v := range props {
		k, v := k, v

		// If we don't have the target observer, don't bother setting it.
		if _, ok := e.observers[k]; !ok {
			continue
		}

		pk, ok := v.(ForwardProp)
		if !ok {
			e.Set(k, v)
			continue
		}
		b.Observe(pk.prop, e.PropSetter(k))
		if pk.def != nil { // Set the default value
			e.Set(k, pk.def)
		}
	}
}

func (b *Builder) push(e *Entity) {
	cur := curEntity{entity: e}

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

// ObsFunc generics way
/*func ObsFunc[T any](fn func(T)) func(interface{}) {
	return func(vv interface{}) {
		v, ok := vv.(T)
		if !ok {
			var z T
			panic(fmt.Sprintf("Can't convert prop %T(%v) to func(%T)", vv, v, z))
		}
		fn(v)
	}
}

func Ptr[T any](p *T) func(interface{}) {
	return func(vv interface{}) {
		v, ok := vv.(T)
		if !ok {
			var z T
			panic(fmt.Sprintf("Can't convert prop %T(%v) to %T", vv, v, z))
		}
		*p = v
	}
}*/

// ObsFunc creates a typed observer func from reflection.
func ObsFunc(fn interface{}) func(interface{}) {
	if fn, ok := fn.(func(interface{})); ok {
		return fn
	}
	fnVal := reflect.ValueOf(fn)
	inTyp := fnVal.Type().In(0)

	return func(v interface{}) {
		arg := reflect.ValueOf(v)
		if aTyp := arg.Type(); aTyp != inTyp {
			if !arg.CanConvert(inTyp) {
				panic(fmt.Sprintf("Can't convert prop %v(%v) to %v", aTyp, v, inTyp))
			}
			arg = arg.Convert(inTyp)
		}
		// Type check somewhere
		fnVal.Call([]reflect.Value{arg})
	}
}

func Ptr(p interface{}) func(interface{}) {
	typ := reflect.TypeOf(p).Elem()
	return func(v interface{}) {
		arg := reflect.ValueOf(v)
		if aTyp := arg.Type(); aTyp != typ {
			if !arg.CanConvert(typ) {
				panic(fmt.Sprintf("Can't convert prop %v(%v) to %v", aTyp, v, typ))
			}
			arg = arg.Convert(typ)
		}
		reflect.ValueOf(p).Elem().Set(arg)
	}
}

/*
// ObsFunc creates a typed observer func from type switch it works on tinygo since
// tiny go doesn't support reflection NumIn.
func ObsFunc(fn interface{}) func(interface{}) {
	switch fn := fn.(type) {
	case func(interface{}):
		return fn
	case func(string):
		return func(v interface{}) {
			fn(v.(string))
		}
	case func(float32):
		return func(v interface{}) { fn(v.(float32)) }
	case func(float64):
		return func(v interface{}) { fn(v.(float64)) }
	case func(int):
		return func(v interface{}) { fn(v.(int)) }
	case func(m32.Vec4):
		return func(v interface{}) { fn(v.(m32.Vec4)) }
	case func(bool):
		return func(v interface{}) { fn(v.(bool)) }
	case func([]text.Align):
		return func(v interface{}) { fn(v.([]text.Align)) }
	case func(text.Overflow):
		return func(v interface{}) { fn(v.(text.Overflow)) }
	default:
		panic(fmt.Sprintf("unsupported observer: %T", fn))
	}
}
*/
