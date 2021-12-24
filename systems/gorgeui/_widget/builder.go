package widget

// TODO: prepare horizontal mode or fix heights by calculating heights on END

import (
	"fmt"

	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/systems/gorgeui"
)

type builder interface {
	W
	Build(b *Builder)
}

// Build builds a widget.
func Build(b builder) {
	// *b.RectTransform() = gorgeui.RectIdent()
	ctx := Builder{
		root: curEntity{widget: b},
		style: BuilderStyle{
			def: cursorStyle{
				background: m32.Vec4{0, 0, 0, 0.2},
				color:      m32.Vec4{1, 1, 1, 1},
				dim:        m32.Vec2{20, 5},
				spacing:    m32.Vec4{1, 1, 1, 1},
			},
		},
	}
	// b.Widget().SetAnchor(0)
	// b.Widget().SetRect(0, 0, 20, 0)
	// b.Widget().SetPivot(0)
	// Defaults to full screen?
	// root.SetAnchor(0, 0, 0, 0)
	// root.SetRect(0, 0, 20, 0)
	// root.SetPivot(0)
	b.Build(&ctx)
}

// BuildFunc calls fn to build a widget and returns the widget.
func BuildFunc(fn func(b *Builder)) *Widget {
	root := New()
	root.SetAnchor(0, 0, 0, 0)
	root.SetRect(0, 0, 20, 0)
	root.SetPivot(0)
	b := Builder{
		root: curEntity{
			layout: layoutVertical(),
			widget: root,
		},
		style: BuilderStyle{
			def: cursorStyle{
				background: m32.Vec4{0, 0, 0, 0.2},
				color:      m32.Vec4{1, 1, 1, 1},
				dim:        m32.Vec2{20, 5},
				spacing:    m32.Vec4{1, 1, 1, 1},
				fontSize:   2,
				textAlign:  [2]AlignType{AlignCenter, AlignCenter},
			},
		},
	}
	// Defaults to full screen?
	fn(&b)
	return root
}

// Direction widget builder flow
type Direction int

// Directions
const (
	// Vertical widgets will be anchored to parent vertically.
	DirectionVertical Direction = iota
	// Horizontal widgets will be anchored to parent horizontally.
	DirectionHorizontal
	// Free widgets will not be anchored.
	DirectionFree
)

type layoutFunc func(w *Component, s *cursorStyle)

type curEntity struct {
	layout layoutFunc
	// cursor position
	// pos    m32.Vec2
	// dir    Direction
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
	w.Widget().SetPivot(0)

	// Used with anchor 0,0,1,0, with will set content to the built height
	// w.Widget().SetRect(1+cur.pos[0], 1+cur.pos[1], 1, cur.height)

	c.layout(w.Widget(), s)
	// TODO: This should be controlled by the parent With a default
	// although it couldn't be the widget it self controlling this but rather
	// having a position controller or something
	/*switch c.dir {
	case DirectionFree:
		w.Widget().SetAnchor(0)
		w.Widget().SetRect(c.pos[0], c.pos[1], s.dim[0], s.dim[1])
		c.pos[0] = 0
		c.pos[1] += s.dim[1] + s.spacing
	case DirectionVertical:
		w.Widget().SetAnchor(0, 0, 1, 0)
		w.Widget().SetRect(s.spacing, s.spacing+c.pos[1], s.spacing, s.dim[1])
		c.pos[1] += s.dim[1] + s.spacing
	case DirectionHorizontal:
		w.Widget().SetAnchor(0, 0, 0, 1)
		w.Widget().SetRect(s.spacing+c.pos[0], s.spacing, s.dim[0], s.spacing)
		c.pos[0] += s.dim[0] + s.spacing
	}*/

	gorgeui.AddChildrenTo(c.widget, w)

	if b.addfn != nil {
		b.addfn(w)
	}
}

func (b *Builder) push(w W, lfn layoutFunc) *curEntity {
	// cur := b.cur()
	if lfn == nil {
		lfn = layoutVertical()
	}
	e := &curEntity{
		layout: lfn,
		// cursorData: cur.cursorData,
		// pos:    m32.Vec2{0, 0},
		// dir:    b.cur().dir,
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

func (b *Builder) begin(w W, lfn layoutFunc, s *cursorStyle) {
	b.style.Save()
	b.add(w, s)
	b.push(w, lfn)
}

func (b *Builder) end() {
	b.style.Restore()
	b.pop()
}

// SetDirection container direction.
// func (b *Builder) SetDirection(d Direction) {
//	c := b.cur()
//	c.dir = d
// }

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
/*func (b *Builder) Begin(w W) {
	b.add(w, b.style.cur())
	b.push(w)
}*/

// End pops a widget from stack.
func (b *Builder) End() {
	b.end()
	// Change cursor to whatever it is
	// cur := b.cur()

	// e.widget.RectTransform().Dim[1] = e.pos[1] + 1
	// cur.pos[1] = e.widget.RectTransform().Dim[1] + 1
}

// BeginPanel starts a panel, the next widgets will be added as childs to panel.
func (b *Builder) BeginPanel(ls ...gorgeui.Layouter) *Panel {
	// cur := b.cur()
	s := b.style.cur() // this won't pop 'once'
	panel := NewPanel()
	if len(ls) > 0 {
		panel.Layouter = gorgeui.MultiLayout(ls...)
	}
	panel.SetColor(s.background[:]...)

	/*dir := ListVertical
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
	// panel.Handle(ResizeToContent(panel, 1))
	// panel.HandleFunc(AutoHeight(panel, 1))
	// } else {
	// panel.HandleFunc(AutoWidth(panel, 1))
	// }
	b.begin(panel, layoutVertical(), s)

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
	label.SetFontScale(s.fontSize)
	label.SetAlign(s.textAlign[:]...)

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
	case func() string:
		label.HandleFunc(func(e event.Event) {
			if _, ok := e.(gorgeui.EventUpdate); !ok {
				return
			}
			txt := v()
			if label.Text != txt {
				label.SetText(txt)
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
	btnLabel.SetFontScale(s.fontSize)
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
	lbl.SetFontScale(s.fontSize)
	gorgeui.AddChildrenTo(btn, lbl)

	var v interface{}
	if len(args) == 0 || args[0] == nil {
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
		t := min + float32(e)*(rng)
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
			lbl.SetText(fmt.Sprintf("err %T", v))
		}
	})
	b.add(slider, s)
	return slider
}

// SpinnerVec3 adds and returns a vec3 spinner.
func (b *Builder) SpinnerVec3(v *m32.Vec3) *SpinnerVec3 {
	s := b.style.cur()
	if v == nil {
		v = &m32.Vec3{}
	}
	w := NewSpinnerVec3(*v)
	w.x.label.SetFontScale(s.fontSize)
	w.x.valueLabel.SetFontScale(s.fontSize)
	w.y.label.SetFontScale(s.fontSize)
	w.y.valueLabel.SetFontScale(s.fontSize)
	w.z.label.SetFontScale(s.fontSize)
	w.z.valueLabel.SetFontScale(s.fontSize)

	b.add(w, s)
	return w
}

// Spinner adds and returns a vec3 spinner.
func (b *Builder) Spinner(l string, f *float32) *Spinner {
	s := b.style.cur()
	if f == nil {
		f = new(float32)
	}
	w := NewSpinner(l, *f)
	w.label.SetFontScale(s.fontSize)
	w.valueLabel.SetFontScale(s.fontSize)
	w.labelBg.SetColor(s.color[:]...)
	b.add(w, s)
	return w
}

// BeginHFlex starts a flex layout where items will be placed horizontally.
func (b *Builder) BeginHFlex(sz ...float32) *Widget {
	w := New()
	b.begin(w, layoutFlex(DirectionHorizontal, sz...), b.style.cur())
	return w
}

// BeginVFlex starts a vertical flex layout.
func (b *Builder) BeginVFlex(sz ...float32) *Widget {
	w := New()
	b.begin(w, layoutFlex(DirectionVertical, sz...), b.style.cur())
	return w
}

// Begin starts an empty widget group.
func (b *Builder) Begin() *Widget {
	w := New()
	b.begin(w, layoutVertical(), b.style.cur())
	return w
}

// BeginList starts a list.
/*func (b *Builder) BeginList(dir ...ListDirection) W {
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

	b.begin(w, b.style.cur())
	return w
}

// EndList finishes a list previously called with BeginList().
func (b *Builder) EndList() {
	b.End()
}*/

// Root returns the root widget of the builder.
func (b *Builder) Root() *Widget {
	return b.root.widget.(*Widget)
}
