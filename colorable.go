package gorge

import "github.com/stdiopt/gorge/math/gm"

// To avoid locking instancing we can use this to set main object color

// ColorableComponent sets A Main color for geometry
type ColorableComponent struct {
	Color gm.Vec4
}

// Colorable returns the colorable component
func (c *ColorableComponent) Colorable() *ColorableComponent { return c }

// NewColorableComponent returns a new colorable
func NewColorableComponent(r, g, b, a float32) *ColorableComponent {
	return &ColorableComponent{gm.Vec4{r, g, b, a}}
}

// SetColor sets the Color.
func (c *ColorableComponent) SetColor(r, g, b, a float32) {
	c.Color = gm.Vec4{r, g, b, a}
}

// SetColorv sets the Color from a vec4.
func (c *ColorableComponent) SetColorv(v gm.Vec4) {
	c.Color = v
}

// GetColor returns the Color
func (c *ColorableComponent) GetColor() gm.Vec4 {
	if c == nil {
		return gm.Vec4{1, 1, 1, 1}
	}
	return c.Color
}
