package gorge

// Contexter interface to return a context containing Gorge instance.
type Contexter interface {
	Gorge() *Context
}

type gorge = Gorge

// Context holds a gorge system and is mostly used as a Prop.
type Context struct {
	*gorge
}

// Gorge context interface helper.
func (c *Context) Gorge() *Context {
	return c
}

// AddContext to gorge context.
func AddContext(g *Context, k, v interface{}) {
	g.addSystem(k, v)
}

// GetContext gets a system from a gorge context.
func GetContext(g *Context, k interface{}) interface{} {
	return g.getSystem(k)
}
