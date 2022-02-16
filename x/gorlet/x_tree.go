package gorlet

import (
	"log"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/text"
)

type TreeNode struct {
	Name     string
	Children []*TreeNode
}

type WTree struct {
	Widget[WTree]

	model *TreeNode

	entry *WPane
	icn   *WLabel
	lbl   *WLabel

	offset    float32 // padding left
	closed    bool
	color     gm.Vec4
	highlight gm.Vec4
	factor    float32
}

// Need to pass some kind of data like strings and string children
func Tree(n *TreeNode) *WTree {
	return Build(&WTree{
		color:     gm.Vec4{},
		highlight: gm.Vec4{.4, 0, 0, 1},
		factor:    20,
		model:     n,
	})
}

func (w *WTree) Build(b *B) {
	// var subtree *gorlet.WContainer
	log.Printf("Building tree for: %v with offset: %v", w.model, w.offset)

	prefix := ""
	if len(w.model.Children) != 0 {
		prefix = "-"
	}
	w.entry = b.BeginPane().
		SetColor(w.color[:]...).
		SetAnchor(0, 0, 1, 0).
		SetSize(0, 3)
	b.BeginContainer().SetRect(w.offset*2, 0, 0, 0)

	w.icn = b.Label(prefix).
		SetAnchor(0, 0, .1, 1).
		SetSize(0)
	w.lbl = b.Label(w.model.Name).
		SetAnchor(.1, 0, .1, 1).
		SetAutoSize(true).
		SetOverflow(text.OverflowOverlap).
		SetSize(0).
		// SetAnchor(.1, 0, 1, 1).
		SetTextAlign(text.AlignStart, text.AlignCenter)
	b.EndContainer()
	b.EndPane()

	off := w.offset + 1
	subtree := b.BeginContainer().
		SetAnchor(0, 0, 1, 0).
		SetRect(0, 3, 0, 3).
		SetLayout(LayoutList(0), AutoHeight(0)) // Children
	{
		for _, c := range w.model.Children {
			t := Build(&WTree{
				model:     c,
				color:     w.color,
				highlight: w.highlight,
				factor:    w.factor,
				offset:    off,
			})
			b.Add(t)
		}
	}
	b.EndContainer()

	targetColor := w.color
	event.Handle(w.entry, func(gorgeui.EventPointerEnter) {
		targetColor = w.highlight

		ui := gorgeui.RootUI(w)
		event.Trigger(ui.Gorge(), gorge.EventCursor(gorge.CursorHand))
	})
	event.Handle(w.entry, func(gorgeui.EventPointerLeave) {
		targetColor = w.color
		ui := gorgeui.RootUI(w)
		event.Trigger(ui.Gorge(), gorge.EventCursor(gorge.CursorArrow))
	})

	event.Handle(w.entry, func(gorgeui.EventPointerUp) {
		if len(w.model.Children) == 0 {
			return
		}
		if w.closed {
			w.icn.SetText("-")
			w.Add(subtree)
		} else {
			w.icn.SetText("+")
			w.Remove(subtree)
		}
		w.closed = !w.closed
	})

	c := targetColor
	height := AutoHeight(0)
	event.Handle(w, func(e gorgeui.EventUpdate) {
		height.Layout(w)
		c = c.Lerp(targetColor, e.DeltaTime()*w.factor)
		w.entry.SetColor(c[:]...)
	})
}

func (w *WTree) SetOffset(v float32) *WTree {
	w.offset = v
	return w
}

func (w *WTree) SetClosed(v bool) *WTree {
	w.closed = v
	return w
}

func (b *B) Tree(n *TreeNode) *WTree {
	w := Tree(n)
	b.Add(w)
	return w
}
