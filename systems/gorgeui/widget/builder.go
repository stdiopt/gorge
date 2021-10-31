package widget

// TODO: prepare horizontal mode or fix heights by calculating heights on END

import (
	"fmt"

	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

// Build calls fn to build a widget and returns the widget.
func Build(fn func(b *Builder)) *Widget {
	root := New()

	b := Builder{
		root: curEntity{
			widget: root,
		},
		style: BuilderStyle{
			def: cursorStyle{
				background: m32.Vec4{0, 0, 0, 0.2},
				color:      m32.Vec4{1, 1, 1, 1},
				dim:        m32.Vec2{20, 5},
			},
		},
	}
	fn(&b)
	return root
}

// cursorData is stacked data
// This could probably contain a map with custom props
type cursorStyle struct {
	// could be Stacked options
	dim        m32.Vec2
	background m32.Vec4
	color      m32.Vec4
}

// BuilderStyle manages the styles of the widget builder.
type BuilderStyle struct {
	def   cursorStyle
	stack []*cursorStyle
	once  *cursorStyle
}

func (b *BuilderStyle) cur() *cursorStyle {
	if b.once != nil {
		s := b.once
		b.once = nil
		return s
	}
	if len(b.stack) == 0 {
		return &b.def
	}
	return b.stack[len(b.stack)-1]
}

// edit same as cur but doesn't remove the once.
func (b *BuilderStyle) edit() *cursorStyle {
	if b.once != nil {
		s := b.once
		return s
	}
	if len(b.stack) == 0 {
		return &b.def
	}
	return b.stack[len(b.stack)-1]
}

// Once returns BuilderStyle with once set.
func (b *BuilderStyle) Once() *BuilderStyle {
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
	s := b.edit()
	s.color = v4Color(c...)
}

// SetBackground sets next widgets background color.
func (b *BuilderStyle) SetBackground(c ...float32) {
	s := b.edit()
	s.background = v4Color(c...)
}

// SetWidth sets next widgets width.
func (b *BuilderStyle) SetWidth(w float32) {
	s := b.edit()
	s.dim[0] = w
}

// SetHeight sets next widgets height.
func (b *BuilderStyle) SetHeight(h float32) {
	s := b.edit()
	s.dim[1] = h
}

type curEntity struct {
	// cursor position
	pos    m32.Vec2
	widget W
}

// Builder a stateful widget builder.
type Builder struct {
	stack []*curEntity
	root  curEntity

	style BuilderStyle

	addfn func(w W)
}

func (b *Builder) cur() *curEntity {
	if len(b.stack) == 0 {
		return &b.root
	}
	return b.stack[len(b.stack)-1]
}

// Add adds a widget.
func (b *Builder) add(w W, s *cursorStyle) {
	c := b.cur()
	w.Widget().SetAnchor(0, 0, 0, 0)
	w.Widget().SetPivot(0)
	// Used with anchor 0,0,1,0, with will set content to the built height
	// w.Widget().SetRect(1+cur.pos[0], 1+cur.pos[1], 1, cur.height)
	w.Widget().SetRect(1+c.pos[0], 1+c.pos[1], s.dim[0], s.dim[1])
	c.pos[1] += s.dim[1] + 1

	gorgeui.AddChildrenTo(c.widget, w)

	if b.addfn != nil {
		b.addfn(w)
	}
}

func (b *Builder) push(w W) *curEntity {
	// cur := b.cur()
	e := &curEntity{
		// cursorData: cur.cursorData,
		pos:    m32.Vec2{0, 0},
		widget: w,
	}
	b.stack = append(b.stack, e)
	return e
}

func (b *Builder) pop() *curEntity {
	e := b.stack[len(b.stack)-1]
	b.stack = b.stack[:len(b.stack)-1]
	return e
}

// Style returns the style manager.
func (b *Builder) Style() *BuilderStyle {
	return &b.style
}

// OnAdd sets a trigger fn that will be called whenever a widget is added.
func (b *Builder) OnAdd(fn func(W)) {
	b.addfn = fn
}

// Add adds a widget.
func (b *Builder) Add(w W) {
	b.add(w, b.style.cur())
}

// Begin Pushes widget onto the stack, the next widgets will be added as childs
// of w.
func (b *Builder) Begin(w W) {
	b.add(w, b.style.cur())
	b.push(w)
}

// End pops a widget from stack.
func (b *Builder) End() {
	e := b.pop()
	cur := b.cur()

	// e.widget.RectTransform().Dim[1] = e.pos[1] + 1
	cur.pos[1] = e.widget.RectTransform().Dim[1] + 1
}

// BeginPanel starts a panel, the next widgets will be added as childs to panel.
func (b *Builder) BeginPanel(d ...ListDirection) *Panel {
	// cur := b.cur()
	s := b.style.cur()
	panel := NewPanel()

	/*

		dir := ListVertical
		if len(d) > 0 {
			dir = d[0]
		}
		list := ListController(panel)
		list.SetDirection(dir)
		list.SetSizeMode(ListSizeOriginal)
		list.SetSpacing(1)
		list.SetPadding(1, 1)
	*/

	// resize := AutoHeight(1)
	// if dir == ListHorizontal {
	panel.HandleFunc(AutoHeight(panel, 1))
	// } else {
	panel.HandleFunc(AutoWidth(panel, 1))
	// }
	b.Begin(panel)
	panel.SetColor(s.background[:]...)

	return panel
}

// EndPanel ends the panel widget previous called with BeginPanel().
func (b *Builder) EndPanel() {
	b.End()
}

// Label adds a label widget.
func (b *Builder) Label(v interface{}) *Label {
	// cur := b.cur()
	s := b.style.cur()
	label := NewLabel()
	label.SetColor(s.color[:]...)

	switch v := v.(type) {
	case string:
		label.SetText(v)
	case *string:
		// Is this a good idea?
		label.HandleFunc(func(e event.Event) {
			if _, ok := e.(gorgeui.EventUpdate); !ok {
				return
			}
			if label.Text != *v {
				label.SetText(*v)
			}
		})
	}
	b.add(label, s)

	return label
}

// TextButton adds a text button widget.
func (b *Builder) TextButton(t string, click func()) *Button {
	s := b.style.cur()
	btnLabel := NewLabel()
	btnLabel.SetText(t)
	btnLabel.SetColor(s.color[:]...)

	button := NewButton()
	if click != nil {
		button.HandleFunc(func(e event.Event) {
			if _, ok := e.(EventClick); !ok {
				return
			}
			click()
		})
	}
	gorgeui.AddChildrenTo(button, btnLabel)

	b.add(button, s)

	return button
}

// Slider adds a slider with min and max values, the first argument of the
// variadic param should be either a:
// - func(float32):  which is called whenever the value changes
// - *float32: value is updated whenever the value changes
// - *int: value is updated as int whenever the value changes
// the second optional argument is the default start value for the slider.
func (b *Builder) Slider(min, max float32, args ...interface{}) *Slider {
	s := b.style.cur()

	rng := max - min

	slider := NewSlider()

	btn := NewButton()
	btn.SetAnchor(.5, 0, .5, 1)
	btn.SetRect(0, 0, 4, 0)
	slider.SetHandler(btn)

	lbl := NewLabel()
	gorgeui.AddChildrenTo(btn, lbl)

	var v interface{}
	if len(args) == 0 {
		var f float32
		v = &f
	} else {
		v = args[0]
	}

	def := v
	if len(args) > 1 {
		def = args[1]
	}

	// Initial value
	switch v := def.(type) {
	case float32:
		lbl.SetText(fmt.Sprintf("%.02f", def))
		slider.SetValue(v / rng)
	case int:
		lbl.SetText(fmt.Sprintf("%d", def))
		slider.SetValue(float32(v) / rng)
	default:
		def = min
		slider.SetValue(0)
		lbl.SetText(fmt.Sprintf("%.2f", def))
	}

	slider.HandleFunc(func(ee event.Event) {
		e, ok := ee.(EventValueChanged)
		if !ok {
			return
		}
		t := min + e.Value*(rng)
		switch v := v.(type) {
		case func(float32):
			v(t)
			lbl.SetText(fmt.Sprintf("%.02f", t))
		case *float32:
			*v = t
			lbl.SetText(fmt.Sprintf("%.02f", *v))
		case *int:
			*v = int(t)
			lbl.SetText(fmt.Sprintf("%d", *v))
		default:
			lbl.SetText("err")
		}
	})
	b.add(slider, s)
	return slider
}

// BeginList starts a list.
func (b *Builder) BeginList(dir ...ListDirection) W {
	d := ListVertical
	if len(dir) > 0 {
		d = dir[0]
	}
	w := New()
	l := ListController(w)
	l.SetDirection(d)
	l.SetSizeMode(ListSizeOriginal)
	l.SetSpacing(1)

	// Improve this
	w.HandleFunc(AutoHeight(w, 1))

	b.Begin(w)
	return w
}

// EndList finishes a list previously called with BeginList().
func (b *Builder) EndList() {
	b.End()
}
