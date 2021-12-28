package gorlet

// Props represents a k.v with some methods
type Props map[string]interface{}

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
func (p Props) Set(k string, v interface{}) {
	if v == nil {
		delete(p, k)
		return
	}
	p[k] = v
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

/*
type propSetter interface {
	Set(string, interface{})
}

type PropsGroup map[string]interface{}

// Apply props to thing.
func (p PropsGroup) Apply(prefix string, s propSetter) {
	if p == nil {
		return
	}
	solve := map[string]string{}
	prefix += "."
	for k := range p {
		pre := ""
		key := k
		if n := strings.LastIndex(k, "."); n > 0 {
			pre = k[:n]
			key = k[n+1:]
		}
		if !strings.HasPrefix(prefix, pre) { // discard non prefix ones
			continue
		}
		if full, ok := solve[key]; !ok || len(full) < len(k) {
			solve[key] = k
		}
	}

	for k, fullK := range solve {
		s.Set(k, p[fullK])
	}
}*/
