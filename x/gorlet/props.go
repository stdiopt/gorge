package gorlet

type props map[string]interface{}

func (p props) Set(k string, v interface{}) {
	p[k] = v
}

func (p props) Clone() props {
	cp := props{}
	for k, v := range p {
		cp[k] = v
	}
	return cp
}

type propStack struct {
	stack []props
	root  props
}

func (p *propStack) cur() props {
	if len(p.stack) > 0 {
		return p.stack[len(p.stack)-1]
	}
	if p.root == nil {
		p.root = props{}
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
