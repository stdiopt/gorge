package gorlet

// ForwardProp to be used to forward properties.
type ForwardProp struct {
	prop string
	def  any
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

	apply []func(e *Entity)

	props Props
}

func (n *nextData) add(fn func(e *Entity)) {
	n.apply = append(n.apply, fn)
}

// Builder used to build a guilet.
type Builder struct {
	next       nextData
	clientArea *Entity
	onAddFn    func(e *Entity)

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
func (b *Builder) ClientArea() {
	b.clientArea = b.cur().entity
}

///////////////////////////////////////////////////////////////////////////////
// Properties related to the next entity
///////////////////////////////////////////////////////////////////////////////

// UsePlacement sets the placement func.
func (b *Builder) SetPlacement(fn EntityFunc) {
	b.cur().placement = fn
}

// Next pushes a func to the next created entity.
func (b *Builder) Next(fn func(e *Entity)) {
	b.next.add(fn)
}

// UseLayout set next widget layout.
func (b *Builder) UseLayout(fns ...Layouter) {
	if len(fns) == 0 {
		return
	}
	b.next.add(func(e *Entity) {
		e.SetLayout(MultiLayout(fns...))
	})
}

// UseDimRect uses dimension rect sets anchor to 0
func (b *Builder) UseDimRect(v ...float32) {
	b.next.add(func(e *Entity) {
		e.SetAnchor(0)
		e.SetPivot(0)
		e.SetRect(v...)
	})
}

// UseRelRect uses relative from parent rect.
func (b *Builder) UseRelRect(v ...float32) {
	b.next.add(func(e *Entity) {
		e.SetAnchor(0, 0, 1, 1)
		e.SetPivot(0)
		e.SetRect(v...)
	})
}

// UseRect sets next Entity Rect.
func (b *Builder) UseRect(v ...float32) {
	b.next.add(func(e *Entity) {
		e.SetRect(v...)
	})
}

// UseWidth sets the next entity width.
func (b *Builder) UseWidth(v float32) {
	b.next.add(func(e *Entity) {
		e.SetWidth(v)
	})
}

// UseWidth sets the next entity Height.
func (b *Builder) UseHeight(v float32) {
	b.next.add(func(e *Entity) {
		e.SetHeight(v)
	})
}

// UseMargin sets the next entity padding.
func (b *Builder) UseMargin(v ...float32) {
	b.next.add(func(e *Entity) {
		e.SetMargin(v...)
	})
}

// UseAnchor sets next Entity Anchor.
func (b *Builder) UseAnchor(v ...float32) {
	b.next.add(func(e *Entity) {
		e.SetAnchor(v...)
	})
}

// UsePivot sets next Entity Pivot.
func (b *Builder) UsePivot(v ...float32) {
	b.next.add(func(e *Entity) {
		e.SetPivot(v...)
	})
}

func (b *Builder) UseDragEvents(v bool) {
	b.next.add(func(e *Entity) {
		e.SetDragEvents(v)
	})
}

// Use a property for the next widget.
func (b *Builder) Use(k string, v any) {
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
func (b *Builder) Prop(k string, v ...any) ForwardProp {
	var def any
	if len(v) > 0 {
		def = v[0]
	}
	return ForwardProp{prop: k, def: def}
}

// ForwardProps will forward any entity props in this builder to the entity
// using a prefix
func (b *Builder) ForwardProps(pre string, e *Entity) {
	for k, v := range e.observers.observers {
		key := k
		if pre != "" {
			key = pre + "." + k
		}
		for _, fn := range v.Funcs {
			b.root.entity.observeWithType(key, v.Type, fn)
		}
	}
}

// Observe adds a function to observe a property in the root Entity.
func (b Builder) Observe(k string, fn any) {
	b.root.entity.Observe(k, fn)
}

// Push will set the prop to any added entity.
func (b *Builder) Push(k string, v any) {
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

	// Should be on All in this container
	if pfn := b.cur().placement; pfn != nil {
		pfn(e)
	}

	// On next one
	for _, fn := range b.next.apply {
		fn(e)
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
	if len(b.root.entity.observers.observers) > 0 {
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
		// shadow
		k, v := k, v

		// If we don't have the target observer, don't bother setting it.
		o := e.observer(k)
		if o == nil {
			continue
		}

		pk, ok := v.(ForwardProp)
		if !ok {
			e.Set(k, v)
			continue
		}

		b.root.entity.ObserveTo(pk.prop, e, k)
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
