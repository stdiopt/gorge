package gorlet

import "strings"

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

func (p Props) Select(prefix string) Props {
	prefix += "."
	solve := map[string]string{}
	for k := range p {
		pre := ""
		key := k
		if n := strings.LastIndex(k, "."); n > 0 {
			pre = k[:n]
			key = k[n+1:]
		}
		// log.Println("Has prefix:", prefix, pre, strings.HasPrefix(prefix, pre))
		if !strings.HasPrefix(prefix, pre) { // discard non prefix ones
			continue
		}
		if full, ok := solve[key]; !ok || len(full) < len(k) {
			solve[key] = k
		}
	}
	r := Props{}
	for k, fullK := range solve {
		r[k] = p[fullK]
	}
	return r
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
