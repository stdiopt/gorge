package gorlet

type BuildFunc func(*B)

type curEntity struct {
	entity Entity
}

type B struct {
	stack []*curEntity
	root  *curEntity

	clientArea Entity
}

func (b *B) Do(fn BuildFunc) {
	fn(b)
	if b.clientArea != nil {
		if ca, ok := b.root.entity.(interface{ setClientArea(Entity) }); ok {
			ca.setClientArea(b.clientArea)
		}
	}
}

func (b *B) ClientArea() {
	b.clientArea = b.cur().entity
}

func (b *B) Begin(e Entity) {
	b.Add(e)
	b.push(e)
}

// End pops the current guilet from the stack.
func (b *B) End() {
	b.pop()
}

func (b *B) Add(e Entity) {
	cur := b.cur()
	cur.entity.add(e)
}

// Might be wrong if it's not custom widget
// but we only need too in custom.
func (b *B) Root() Entity {
	return b.root.entity
}

func (b *B) SetRoot(e Entity) {
	if len(b.stack) > 0 {
		panic("Builder.Start() called while already in a container")
	}
	b.root.entity = e
}

func (b *B) push(e Entity) {
	cur := curEntity{entity: e}

	b.stack = append(b.stack, &cur)
}

func (b *B) pop() {
	if len(b.stack) == 0 {
		return
	}
	t := b.stack
	b.stack = b.stack[:len(b.stack)-1]
	t[len(t)-1] = nil
}

func (b *B) cur() *curEntity {
	if len(b.stack) == 0 {
		return b.root
	}
	return b.stack[len(b.stack)-1]
}
