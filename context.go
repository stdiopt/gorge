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

// AddSystem to gorge context.
func AddSystem(g *Context, k, v interface{}) {
	g.addSystem(k, v)
}

// GetSystem gets a system from a gorge context.
func GetSystem(g *Context, k interface{}) interface{} {
	return g.getSystem(k)
}
