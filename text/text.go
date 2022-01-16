// Package text implements texture rendering from font texture map
package text

import (
	"fmt"
	"unicode"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/static"
)

// Overflow type for specifying text overflow.
type Overflow int

const (
	// OverflowOverlap let the text overflow the boundary
	OverflowOverlap = Overflow(iota)
	// OverflowWordWrap breaks lines on spaces and char by char if no space
	OverflowWordWrap
	// OverflowBreakWord breaks words
	OverflowBreakWord
)

func (o Overflow) String() string {
	switch o {
	case OverflowOverlap:
		return "OverflowOverlap"
	case OverflowWordWrap:
		return "OverflowWordWrap"
	case OverflowBreakWord:
		return "OverflowBreakWord"
	}
	return "Overflow#unknown"
}

// Align type for specifying text alignment.
type Align int

// Text Alignment types.
const (
	AlignStart = Align(iota)
	AlignCenter
	AlignEnd
)

func (a Align) String() string {
	switch a {
	case AlignStart:
		return "AlignStart"
	case AlignCenter:
		return "AlignCenter"
	case AlignEnd:
		return "AlignEnd"
	}
	return "Align#unknown"
}

// Mode is a type to indicate text flow mode.
type Mode int

// Flow modes
const (
	ModeFlow = Mode(iota)
	ModeRaw
)

func (m Mode) String() string {
	switch m {
	case ModeFlow:
		return "ModeFlow"
	case ModeRaw:
		return "ModeRaw"
	}
	return "Mode#unknown"
}

type mesh = gorge.Mesh

// Mesh to be used on renderable update must be called to update the mesh after
// changing.
type Mesh struct {
	mesh

	Text string
	// Boundary counting from starting Point
	Boundary m32.Vec2
	// Wrap just naive string wrap, should have better options
	Alignment Align
	Overflow  Overflow
	// Size 0 will return a default size (1)
	Size float32
	Mode Mode

	Font *Font

	// TODO: {lpf} Calculated values should be private and functions added
	Min, Max m32.Vec2
	Lines    int

	// Allow this to be swapped?

	// cached mesh and meshData
	meshData gorge.MeshData
}

// NewMesh returns a new Text Mesh.
func NewMesh(font *Font, opts ...MeshFunc) *Mesh {
	m := &Mesh{
		Font: font,
		Size: 1,
		meshData: gorge.MeshData{
			Format:  gorge.VertexFormatPT(),
			Updates: 1,
		},
	}
	m.mesh = *gorge.NewMesh(&m.meshData)
	for _, fn := range opts {
		fn(m)
	}

	return m
}

// SetFont sets the font used to draw the vertices, Material texture should
// change too to show the assigned font.
func (m *Mesh) SetFont(font *Font) {
	m.Font = font
	m.Update()
}

// SetMode sets text flow mode and updates underlying mesh.
func (m *Mesh) SetMode(mm Mode) {
	m.Mode = mm
	m.Update()
}

// SetText sets the text and updates underlying mesh.
func (m *Mesh) SetText(a ...any) {
	m.Text = fmt.Sprint(a...)
	m.Update()
}

// SetBoundary set text boundary, and updates mesh.
func (m *Mesh) SetBoundary(w, h float32) {
	m.Boundary = m32.Vec2{w, h}
	m.Update()
}

// SetOverflow sets text overflow based on boundary and updates underlying mesh.
func (m *Mesh) SetOverflow(o Overflow) {
	m.Overflow = o
	m.Update()
}

// SetSize sets the relative font size.
func (m *Mesh) SetSize(v float32) {
	m.Size = v
	m.Update()
}

// SetAlignment sets the horizontal alignment based on boundary and updates
// underlying mesh.
func (m *Mesh) SetAlignment(a Align) {
	m.Alignment = a
	m.Update()
}

// getSize and prevent returning a 0 to avoid division by 0
func (m *Mesh) getSize() float32 {
	if m.Size != 0 {
		return m.Size
	}
	return 1
}

// MeasureWidth measures a width of a text best usage without "\n"
func (m *Mesh) MeasureWidth(w []rune) float32 {
	size := m.getSize()
	var x, totalX float32
	for i, ch := range w {

		kern := float32(0)
		if i > 0 && m.Font.face != nil {
			k := m.Font.face.Kern(w[i-1], ch)
			kern = float32(k>>6) / 72
		}
		g := m.Font.getGlyph(ch)
		totalX = (x + g.BearingH + g.Size[0]) * size
		x += g.Advance + kern
	}
	return totalX
}

// Update Mesh
// Ideas:
//  - grab a line, for { measure, remove last word } until we get a proper size
func (m *Mesh) Update() {
	// Or raw
	if m.Mode == ModeRaw {
		m.updateRaw()
		return
	}
	m.updateFlow()
}

func (m *Mesh) updateRaw() {
	text := []rune(m.Text)
	size := m.getSize()
	m.Min = m32.Vec2{}
	m.Max = m32.Vec2{}
	m.Lines = 1
	m.meshData.Vertices = m.meshData.Vertices[:0]
	m.meshData.Updates++

	var x, y float32
	isFirstChar := true // for size calculation

	for i, ch := range text {
		switch ch {
		case ' ':
			x += m.Font.SpaceAdv
			continue
		case '\n':
			y--
			m.Lines++
			x = 0
			continue
		case '\t':
			x += m.Font.SpaceAdv * 4
			continue
		}

		kern := float32(0)
		if i > 0 && m.Font.face != nil {
			k := m.Font.face.Kern(text[i-1], ch)
			kern = float32(k>>6) / 72
		}

		g := m.Font.getGlyph(ch)

		xpos := (x + g.BearingH) + kern
		if (xpos+g.Size[0])*size > m.Boundary[0] {
			x = 0
			y--
			m.Lines++
			xpos = (x + g.BearingH) + kern
		}
		ypos := (y + g.Size[1] - g.BearingV) - 0.5 // chars will be centered as vertices 0.5?

		x1 := xpos * size
		y1 := ypos * size
		x2 := (xpos + g.Size[0]) * size
		y2 := (ypos - g.Size[1]) * size

		m.meshData.Vertices = append(m.meshData.Vertices,
			x1, y1, 0, g.Uv1[0], g.Uv1[1], // A
			x2, y1, 0, g.Uv2[0], g.Uv1[1], // B
			x2, y2, 0, g.Uv2[0], g.Uv2[1], // C
			x2, y2, 0, g.Uv2[0], g.Uv2[1], // C
			x1, y2, 0, g.Uv1[0], g.Uv2[1], // D
			x1, y1, 0, g.Uv1[0], g.Uv1[1], // A
		)
		if y == 0 && isFirstChar {
			m.Min = m32.Vec2{x1, y2}
			m.Max = m32.Vec2{x2, y1}
			isFirstChar = false
		} else {
			// update min max could be done getting mesh bounds too
			m.Min[0] = m32.Min(m.Min[0], x1)
			m.Max[0] = m32.Max(m.Max[0], x2)
			m.Min[1] = m32.Min(m.Min[1], y2)
			m.Max[1] = m32.Max(m.Max[1], y1)
		}
		x += g.Advance + kern
	}
}

func (m *Mesh) updateFlow() {
	text := []rune(m.Text)
	size := m.getSize()
	m.Min = m32.Vec2{}
	m.Max = m32.Vec2{}
	m.Lines = 0
	m.meshData.Vertices = m.meshData.Vertices[:0]
	m.meshData.Updates++

	var x, y float32
	isFirstChar := true // for size calculation
	// TODO: {lpf} find better way to reflow instead of a bunch of checks
	// the rules would be
	// - consider tabs (with or without alignment)
	// - if there is a wrap trim?, if not

	// Line treatment
	for len(text) > 0 {
		// Instead of trimming left we could Trim the Line if needed?
		for len(text) > 0 && unicode.IsSpace(text[0]) {
			if text[0] == '\n' {
				y--
				m.Lines++
			}
			text = text[1:]
		}
		// Something to grab a fit chunk
		line := grabLine(text)
		if len(line) == 0 {
			if len(text) > 0 {
				text = text[1:]
			}
			continue
		}
		w := m.MeasureWidth(line)
		cut := 1
		switch m.Overflow {
		case OverflowWordWrap: // And wrap
			for ; len(line) > 1 && w > m.Boundary[0]; w = m.MeasureWidth(line) {
				line = discardLastWord(line)
				cut = 0
				if len(line) <= 1 {
					break
				}
			}

		// TODO: {lpf} Not very efficient? as it is looking backward
		// might be better looking forward?
		case OverflowBreakWord:
			for ; len(line) > 1 && w > m.Boundary[0]; w = m.MeasureWidth(line) {
				line = line[:len(line)-1]
				cut = 0
				if len(line) <= 1 {
					break
				}
			}
		}
		n := len(line) + cut
		if n > len(text) {
			n = len(text)
		}
		text = text[n:]

		/*
			log.Printf("Untrim: %q Trimmed %q", string(line), string(TrimSpaces(line)))
			line = TrimSpaces(line)
			w = m.measureWidth(line)
		*/

		switch m.Alignment {
		case AlignStart:
			x = 0
		case AlignCenter:
			x = (m.Boundary[0]/2 - w/2) / size
		case AlignEnd:
			x = (m.Boundary[0] - w) / size
		}

		for i := 0; i < len(line); i++ {
			ch := line[i]
			if ch == ' ' {
				x += m.Font.SpaceAdv
				continue
			}

			kern := float32(0)
			if i > 0 && m.Font.face != nil {
				k := m.Font.face.Kern(line[i-1], ch)
				// kern = float32(k>>6) / (72 * 0.4)
				kern = float32(k>>6) / 72
			}

			g := m.Font.getGlyph(ch)

			xpos := (x + g.BearingH) + kern
			// chars will be centered as vertices 0.5?
			ypos := (y + g.Size[1] - g.BearingV) - 0.5

			x1 := xpos * size
			y1 := ypos * size
			x2 := (xpos + g.Size[0]) * size
			y2 := (ypos - g.Size[1]) * size

			m.meshData.Vertices = append(m.meshData.Vertices,
				x1, y1, 0, g.Uv1[0], g.Uv1[1], // A
				x2, y1, 0, g.Uv2[0], g.Uv1[1], // B
				x2, y2, 0, g.Uv2[0], g.Uv2[1], // C
				x2, y2, 0, g.Uv2[0], g.Uv2[1], // C
				x1, y2, 0, g.Uv1[0], g.Uv2[1], // D
				x1, y1, 0, g.Uv1[0], g.Uv1[1], // A
			)
			if y == 0 && isFirstChar {
				m.Min = m32.Vec2{x1, y2}
				m.Max = m32.Vec2{x2, y1}
				isFirstChar = false
			} else {
				// update min max could be done getting mesh bounds too
				m.Min[0] = m32.Min(m.Min[0], x1)
				m.Max[0] = m32.Max(m.Max[0], x2)
				m.Min[1] = m32.Min(m.Min[1], y2)
				m.Max[1] = m32.Max(m.Max[1], y1)
			}
			x += g.Advance + kern
		}
		m.Lines++
		y--
	}
}

type textMesh = Mesh

// Entity renderable entity
type Entity struct {
	gorge.TransformComponent
	gorge.RenderableComponent
	gorge.ColorableComponent

	*textMesh
}

// New setups the renderable with a font
func New(font *Font, opts ...MeshFunc) *Entity {
	mat := gorge.NewShaderMaterial(static.Shaders.Unlit)
	mat.SetTexture("albedoMap", font.Texture)

	textMesh := NewMesh(font, opts...)
	var _ gorge.Mesher = textMesh

	t := &Entity{
		TransformComponent:  *gorge.NewTransformComponent(),
		ColorableComponent:  *gorge.NewColorableComponent(1, 1, 1, 1),
		RenderableComponent: *gorge.NewRenderableComponent(textMesh, mat),
		textMesh:            textMesh,
	}
	return t
}

// SetTextf sets formated text
func (t *Entity) SetTextf(f string, args ...any) {
	t.SetText(fmt.Sprintf(f, args...))
}

// MeshFunc allows mesh function options to be added in initializator
type MeshFunc func(t *Mesh)

// Boundary returns a func that applies a boundary to mesh.
func Boundary(w, h float32) MeshFunc {
	return func(t *Mesh) {
		t.Boundary = m32.Vec2{w, h}
	}
}
