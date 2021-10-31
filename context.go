package gorge

type gorge = Gorge

// Context holds a gorge system and is mostly used as a Prop.
type Context struct {
	*gorge
}

// Gorge context interface helper.
func (c *Context) Gorge() *Context {
	return c
}
