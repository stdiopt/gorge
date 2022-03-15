package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
)

type Overflow int

const (
	OverflowVisible = Overflow(iota)
	OverflowHidden
	OverflowScroll
)

func (o Overflow) String() string {
	switch o {
	case OverflowHidden:
		return "OverflowHidden"
	case OverflowScroll:
		return "OverflowScroll"
	default:
		return "OverflowVisible"
	}
}

type WPanel struct {
	Widget[*WPanel]

	main       *WCustom
	background *WPane
	wrapper    Entity
	content    *WContainer

	overflow   Overflow
	scrollSize gm.Vec2
}

// Panel creates a new panel widget.
func Panel(cs ...gorge.Entity) *WPanel {
	w := Build(&WPanel{
		overflow:   OverflowVisible,
		scrollSize: gm.Vec2{1, 1},
	})
	return w.Add(cs...)
}

func (w *WPanel) Build(b *B) {
	w.background = b.Pane()
	w.wrapper = b.BeginContainer()
	w.content = b.Container()
	b.EndContainer()
	w.SetClientArea(w.content)
}

func (w *WPanel) SetLayout(l ...Layouter) *WPanel {
	w.content.SetLayout(l...)
	return w
}

func (w *WPanel) SetColor(vs ...float32) *WPanel {
	w.background.SetColor(vs...)
	return w
}

func (w *WPanel) SetBorderColor(vs ...float32) *WPanel {
	w.background.SetBorderColor(vs...)
	return w
}

func (w *WPanel) SetOverflow(o Overflow) *WPanel {
	/*
		Remove wrapper and replace with one that prevents overflowing
	*/
	if w.overflow == o {
		return w
	}
	w.overflow = o
	w.SetClientArea(nil)
	w.Remove(w.wrapper)
	switch o {
	case OverflowHidden:
		w.wrapper = Mask(w.content)
	case OverflowScroll:
		w.wrapper = Scroll(w.content).
			SetScrollSize(w.scrollSize[:]...)
	default:
		w.content.setMaskDepth(-1)
		w.wrapper = Container(w.content)
	}
	w.Add(w.wrapper)
	w.SetClientArea(w.content)
	return w
}

func (w *WPanel) SetScrollSize(s ...float32) *WPanel {
	switch len(s) {
	case 0:
		w.scrollSize = gm.Vec2{1, 1}
	case 1:
		w.scrollSize = gm.Vec2{s[0], s[0]}
	default:
		w.scrollSize = gm.Vec2{s[0], s[1]}
	}
	if w.overflow == OverflowScroll {
		w.wrapper.(*WScroll).SetScrollSize(w.scrollSize[:]...)
	}
	return w
}

func (b *B) BeginPanel() *WPanel {
	p := Panel()
	b.Begin(p)
	return p
}

func (b *B) EndPanel() {
	b.End()
}
