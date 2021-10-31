package render

// New: 04-05-2021 Manual render component

// Component component to be injected in entities for manual rendering.
// TODO: {lpf} useless for now.
type Component struct {
	Draw func()
}

// Render to satisfy an interface
func (v *Component) Render() *Component { return v }
