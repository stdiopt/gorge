package gorlet

import (
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

func Mask() Func {
	return func(b *Builder) {
		root := b.Root()
		// This will Clip input
		root.Masked = true
		maskOn := Create(Quad())
		maskOn.Set("colorMask", &[4]bool{false, false, false, false})
		maskOn.SetDisableRaycast(true)

		root.AddElement(maskOn)

		container := b.BeginContainer()
		b.ClientArea()
		b.EndContainer()

		maskOff := b.Create(Quad())
		maskOff.Set("colorMask", &[4]bool{false, false, false, false})
		maskOff.SetDisableRaycast(true)
		root.AddElement(maskOff)

		depthMask := 0

		Observe(b, "_maskDepth", func(n int) {
			depthMask = n
		})

		event.Handle(root, func(gorgeui.EventUpdate) {
			maskOn.Set("stencil", calcMaskOn(depthMask))
			maskOff.Set("stencil", calcMaskOff(depthMask))
			for _, c := range container.Children() {
				c.Set("_maskDepth", depthMask+1)
			}
		})
	}
}

// BeginMask pushes a Mask entity into builder.
func (b *Builder) BeginMask() *Entity {
	return b.Begin(Mask())
}

// EndMask convinient alias to End.
func (b *Builder) EndMask() {
	b.End()
}
