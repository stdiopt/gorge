package gorge

// Contexter interface to return a context containing Gorge instance.
type Contexter interface {
	Add(...Entity)
	G() *Context
}

type gorge = Gorge

// Context holds a gorge system and is mostly used as a Prop.
type Context struct {
	*gorge
}

// Gorge context interface helper.
func (c *Context) G() *Context {
	return c
}

// SetContext to gorge context.
func SetContext[T comparable](g *Context, c T) T {
	if _, ok := GetContext[T](g); ok {
		panic("context already exists")
	}
	g.contexts = append(g.contexts, c)
	return c
}

// GetContext gets a system from a gorge context.
func GetContext[T comparable](g *Context) (T, bool) {
	for _, c := range g.contexts {
		if c, ok := c.(T); ok {
			return c, true
		}
	}
	var z T
	return z, false
}
