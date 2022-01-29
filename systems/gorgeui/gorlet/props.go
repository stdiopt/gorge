package gorlet

// Props represents a k.v with some methods
type Props map[string]any

// Clone properties returns a new map with all keys from previous.
func (p Props) Clone() Props {
	cp := Props{}
	for k, v := range p {
		cp[k] = v
	}
	return cp
}

// Merge will merge properties p2 will override properties of receiver.
func (p Props) Merge(p2 Props) Props {
	ret := p.Clone()
	for k, v := range p2 {
		ret[k] = v
	}
	return ret
}

// Set a property to props.
func (p Props) Set(k string, v any) {
	if v == nil {
		delete(p, k)
		return
	}
	p[k] = v
}

func (p Props) SetProps(p2 Props) {
	for k, v := range p2 {
		p.Set(k, v)
	}
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
