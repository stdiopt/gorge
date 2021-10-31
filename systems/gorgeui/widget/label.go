package widget

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/static"
	"github.com/stdiopt/gorge/systems/gorgeui"
	"github.com/stdiopt/gorge/text"
)

// AlignType label text alignment type.
type AlignType = text.AlignType

// Text alignments.
const (
	AlignStart  = text.AlignStart
	AlignCenter = text.AlignCenter
	AlignEnd    = text.AlignEnd
)

// Label widget state.
type Label struct {
	Component
	Text      string
	Color     m32.Vec4
	Overflow  text.Overflow
	Mode      text.Mode
	Size      float32
	Font      *text.Font
	Alignment [2]AlignType

	entity *text.Entity
}

// HandleEvent handles events for label.
func (l *Label) HandleEvent(e event.Event) {
	_, ok := e.(gorgeui.EventUpdate) // Change to PreUpdate?
	if !ok {
		return
	}
	r := l.Rect()
	ent := l.entity
	bounds := m32.Vec2{r[2] - r[0], r[3] - r[1]}
	ent.Position[0] = r[0] // left

	update := false
	if ent.Size != l.Size {
		ent.Size = l.Size
		update = true
	}
	if ent.Overflow != l.Overflow {
		ent.Overflow = l.Overflow
		update = true
	}
	if ent.Mode != l.Mode {
		ent.Mode = l.Mode
		update = true
	}
	if ent.Boundary != bounds {
		ent.Boundary = bounds
		update = true
	}
	if ent.Alignment != l.Alignment[0] {
		ent.Alignment = l.Alignment[0]
		update = true
	}
	if ent.Text != l.Text {
		ent.Text = l.Text
		update = true
	}
	if ent.Font != l.Font {
		ent.Material.SetTexture("albedoMap", l.Font.Texture)
		ent.Font = l.Font
		update = true
	}
	if update {
		ent.Update()
	}

	// This is executed regardless the text change
	textHeight := float32(ent.Lines) * l.Size
	switch l.Alignment[1] {
	case AlignStart:
		ent.Position[1] = r[1] + l.Size*0.25
	case AlignCenter:
		ent.Position[1] = r[3] - (bounds[1]/2 + textHeight/2) + l.Size*0.25 // top, center
	case AlignEnd:
		ent.Position[1] = r[3] - textHeight + l.Size*0.25
	}
	ent.SetColorv(l.Color)
}

// NewLabel returns a label widget
func NewLabel(s ...string) *Label {
	mat := gorge.NewShaderMaterial(static.Shaders.Unlit)
	mat.SetQueue(100)
	mat.SetDepth(gorge.DepthNone)
	mat.SetTexture("albedoMap", gorgeui.DefaultFont.Texture)

	textEnt := text.New(gorgeui.DefaultFont)
	textEnt.SetMaterial(mat)
	textEnt.SetScale(1, -1, 1)

	var defText string
	if len(s) > 0 {
		defText = s[0]
	}
	l := &Label{
		Component: *NewComponent(),
		Color:     v4Color(0, 1),
		Overflow:  text.OverflowWordWrap,
		Alignment: [2]AlignType{AlignCenter, AlignCenter},
		Size:      2,
		Font:      gorgeui.DefaultFont,
		Text:      defText,
		entity:    textEnt,
	}
	gorgeui.AddChildrenTo(l, l.entity)

	return l
}

// SetFont sets the label font.
func (l *Label) SetFont(font *text.Font) {
	l.Font = font
}

// SetFontScale sets the font scale.
func (l *Label) SetFontScale(s float32) {
	l.Size = s
}

// SetText sets the label text.
func (l *Label) SetText(s string) {
	l.Text = s
}

// SetColor sets the label text color.
func (l *Label) SetColor(c ...float32) {
	l.Color = v4Color(c...)
}

// SetOverflow sets text overflow behaviour.
func (l *Label) SetOverflow(o text.Overflow) {
	l.Overflow = o
}

// SetAlign sets the label text alignment.
func (l *Label) SetAlign(a ...AlignType) {
	switch len(a) {
	case 0:
		l.Alignment = [2]AlignType{
			AlignStart,
			AlignStart,
		}
	case 1:
		l.Alignment = [2]AlignType{a[0], a[0]}
	default:
		l.Alignment = [2]AlignType{a[0], a[1]}
	}
}
