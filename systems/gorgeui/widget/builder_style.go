package widget

import (
	"github.com/stdiopt/gorge/m32"
)

// cursorData is stacked data
// This could probably contain a map with custom props
type cursorStyle struct {
	// could be Stacked options
	dim        m32.Vec2
	background m32.Vec4
	color      m32.Vec4
	fontSize   float32
	textAlign  [2]AlignType
	// Could be called margin somehow
	spacing m32.Vec4
}

// BuilderStyle manages the styles of the widget builder.
type BuilderStyle struct {
	def   cursorStyle
	stack []*cursorStyle
	once  *cursorStyle
}

func (b *BuilderStyle) cur() *cursorStyle {
	cur := b.edit()
	if b.once != nil {
		b.once = nil
	}
	return cur
}

// edit same as cur but doesn't remove the once.
func (b *BuilderStyle) edit() *cursorStyle {
	if b.once != nil {
		return b.once
	}
	if len(b.stack) == 0 {
		return &b.def
	}
	return b.stack[len(b.stack)-1]
}

// Once returns BuilderStyle with once set.
func (b *BuilderStyle) Once() *BuilderStyle {
	if b.once != nil {
		return b
	}
	s := *b.cur() // copy
	b.once = &s
	return b
}

// Save saves the current style, previous style can be restored with Restore().
func (b *BuilderStyle) Save() {
	s := *b.cur() // copy
	b.stack = append(b.stack, &s)
}

// Restore restores the previous style.
func (b *BuilderStyle) Restore() {
	if len(b.stack) == 0 {
		return
	}
	b.stack = b.stack[:len(b.stack)-1]
}

// Reset resets the style to the root style.
func (b *BuilderStyle) Reset() {
	b.stack = b.stack[:0]
}

// SetColor sets the next widget colors.
func (b *BuilderStyle) SetColor(c ...float32) {
	b.edit().color = v4Color(c...)
}

// SetBackground sets next widgets background color.
func (b *BuilderStyle) SetBackground(c ...float32) {
	b.edit().background = v4Color(c...)
}

// SetWidth sets next widgets width.
func (b *BuilderStyle) SetWidth(w float32) {
	b.edit().dim[0] = w
}

// SetHeight sets next widgets height.
func (b *BuilderStyle) SetHeight(h float32) {
	b.edit().dim[1] = h
}

// SetSpacing of the next widget.
func (b *BuilderStyle) SetSpacing(s ...float32) {
	switch len(s) {
	case 0:
		b.edit().spacing = m32.Vec4{}
	case 1:
		b.edit().spacing = m32.Vec4{s[0], s[0], s[0], s[0]}
	case 2:
		b.edit().spacing = m32.Vec4{s[0], s[1], s[0], s[1]}
	case 3:
		b.edit().spacing = m32.Vec4{s[0], s[1], s[2], s[1]}
	default:
		b.edit().spacing = m32.Vec4{s[0], s[1], s[2], s[3]}

	}
}

// SetFontSize sets the font size for next wiget.
func (b *BuilderStyle) SetFontSize(s float32) {
	b.edit().fontSize = s
}

// SetTextAlign aligns next labels.
func (b *BuilderStyle) SetTextAlign(a ...AlignType) {
	align := b.edit().textAlign
	switch len(a) {
	case 0:
		align = [2]AlignType{
			AlignStart,
			AlignStart,
		}
	case 1:
		align[0] = a[0]
	default:
		align = [2]AlignType{a[0], a[1]}
	}
	b.edit().textAlign = align
}

func layoutVertical() layoutFunc {
	var pos m32.Vec2
	return func(w *Component, s *cursorStyle) {
		w.SetAnchor(0, 0, 1, 0)
		w.SetRect(s.spacing[0], s.spacing[1]+pos[1], s.spacing[2], s.dim[1])
		pos[1] += s.dim[1] + s.spacing[3]
	}
}

func layoutFlex(dir Direction, sz ...float32) layoutFunc {
	var list []*Component
	var sum float32
	for _, f := range sz {
		sum += f
	}
	var start float32
	return func(w *Component, s *cursorStyle) {
		i := len(list)
		if i >= len(sz) {
			panic("out of bounds")
		}
		end := start + sz[i]/sum
		switch dir {
		case DirectionHorizontal:
			w.SetAnchor(start, 0, end, 1)
		case DirectionVertical:
			w.SetAnchor(0, start, 1, end)
		}
		w.SetRect(s.spacing[0], s.spacing[1], s.spacing[2], s.spacing[3])
		start = end
		list = append(list, w)
	}
}
