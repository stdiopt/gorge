package gorlet

// Props to set multiple properties at once.
type Props map[string]interface{}

func (p Props) Set(k string, v interface{}) {
	p[k] = v
}

func (p Props) Clone() Props {
	cp := Props{}
	for k, v := range p {
		cp[k] = v
	}
	return cp
}

type propStack struct {
	stack []Props
	root  Props
}

func (p *propStack) cur() Props {
	if len(p.stack) > 0 {
		return p.stack[len(p.stack)-1]
	}
	if p.root == nil {
		p.root = Props{}
	}
	return p.root
}

func (p *propStack) Save() {
	p.stack = append(p.stack, p.cur().Clone())
}

func (p *propStack) Restore() {
	if len(p.stack) == 0 {
		return
	}
	p.stack = p.stack[:len(p.stack)-1]
}

// PropGroup used to apply props to an entity
type PropGroup map[string]Props

// Apply named group to entity.
func (p PropGroup) Apply(k string, e *Entity) *Entity {
	props, ok := p[k]
	if !ok {
		return e
	}
	for k, v := range props {
		e.Set(k, v)
	}
	return e
}
