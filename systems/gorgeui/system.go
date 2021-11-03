package gorgeui

import (
	"log"
	"math/rand"
	"sort"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/m32/ray"
	"github.com/stdiopt/gorge/systems/input"
	"github.com/stdiopt/gorge/systems/resource"
	"github.com/stdiopt/gorge/text"
)

// System initializes UI system in gorge.
func System(g *gorge.Context, rc *resource.Context) error {
	log.Println("Initializing system")
	dbg := newDebugLines()
	dbg.Queue = 200
	dbg.SetCullMask(1 << 17)
	g.Add(dbg)

	DefaultFont = &text.Font{}
	if err := rc.Load(DefaultFont, "_gorge/fonts/font.ttf"); err != nil {
		return err
	}

	s := &system{
		gorge: g,
		font:  DefaultFont,
		dbg:   dbg,
	}
	g.Handle(s)
	g.PutProp(func() *Context {
		return &Context{s}
	})
	return nil
}

type system struct {
	Debug DebugFlag

	gorge *gorge.Context
	font  *text.Font

	uis   map[*UI]struct{}
	elems []Entity // This could be upthere?

	// Should have these per pointer
	// meaning mouse button[0-5], touch{1...}
	pointOver    Entity
	pointDown    Entity
	pointDownPos m32.Vec2

	dragging Entity

	curMouse   m32.Vec2
	deltaMouse m32.Vec2
	dbg        *debugLines
}

func (s *system) HandleEvent(v event.Event) {
	switch e := v.(type) {
	// TODO: Warning experimental code here
	case input.EventPointer:
		s.deltaMouse = e.Pointers[0].Pos.Sub(s.curMouse)
		s.curMouse = e.Pointers[0].Pos
		hit, r := s.rayTest(s.curMouse)

		curDown := s.pointDown

		if e.Type == input.MouseDown && e.Button == 0 {
			s.pointDown = hit
			s.pointDownPos = s.curMouse
			if hit != nil {
				p := &PointerData{
					RayResult: r,
					Delta:     s.deltaMouse,
					Position:  s.curMouse,
					Target:    hit,
				}
				EachParent(hit, func(e Entity) bool {
					triggerOn(e, EventPointerDown{p})
					return !p.stopPropagation
				})
			}
		}
		if e.Type == input.MouseUp && e.Button == 0 {
			if s.pointDown != nil {
				p := &PointerData{
					RayResult: r,
					Delta:     s.deltaMouse,
					Position:  s.curMouse,
					Target:    hit,
				}
				EachParent(s.pointDown, func(e Entity) bool {
					triggerOn(e, EventPointerUp{p})
					return !p.stopPropagation
				})
				s.pointDown = nil
			}
			if s.dragging != nil {
				p := &PointerData{
					RayResult: r,
					Delta:     s.deltaMouse,
					Position:  s.curMouse,
					Target:    hit,
				}
				triggerOn(s.dragging, EventDragEnd{p})
				s.dragging = nil
			}
		}

		// Drag detection
		// I mouse is still down and not nil and pointerDown is still and dragging is nil
		if s.pointDown != nil && curDown == s.pointDown && s.dragging == nil {
			p := &PointerData{
				RayResult: r,
				Delta:     s.deltaMouse,
				Position:  s.curMouse,
				Target:    hit,
			}
			ui := RootUI(s.pointDown)
			d := s.curMouse.Sub(s.pointDownPos).Abs()
			if d[0] > ui.DragThreshold || d[1] > ui.DragThreshold {
				EachParent(hit, func(e Entity) bool {
					if !e.Element().DragEvents {
						return true
					}
					s.dragging = e
					triggerOn(s.dragging, EventDragBegin{p})
					return false
				})
			}
		} else if s.dragging != nil {
			p := &PointerData{
				RayResult: r,
				Delta:     s.deltaMouse,
				Position:  s.curMouse,
				Target:    hit,
			}
			triggerOn(s.dragging, EventDrag{p})
		}

	case gorge.EventPreUpdate:
		if s.Debug != 0 {
			s.dbg.Clear()
		}

	case gorge.EventPostUpdate:
		for _, el := range s.elems {
			triggerOn(el, EventUpdate(e.DeltaTime()))
		}

		if s.Debug&DebugRects != 0 {
			s.debugRects()
		}

	case gorge.EventAddEntity:
		if v, ok := e.Entity.(Entity); ok {
			s.addEntity(v)
		}
	case gorge.EventRemoveEntity:
		if v, ok := e.Entity.(Entity); ok {
			s.removeEntity(v)
		}
	}
}

func (s *system) addEntity(e Entity) {
	if ui, ok := e.(*UI); ok {
		if s.uis == nil {
			s.uis = map[*UI]struct{}{}
		}
		s.uis[ui] = struct{}{}
		return
	}
	s.elems = append(s.elems, e)

	el := e.Element()
	el.Attached = true
	if v, ok := e.(Attacher); ok {
		v.Attached(e)
	}
}

func (s *system) removeEntity(e Entity) {
	if ui, ok := e.(*UI); ok {
		if s.uis == nil {
			return
		}
		delete(s.uis, ui)
		return
	}

	wc := e.Element()
	for i, el := range s.elems {
		if e == el {
			t := s.elems
			s.elems = append(s.elems[:i], s.elems[i+1:]...)
			t[len(t)-1] = nil // remove last element since it was copied
		}
	}
	wc.Attached = false
	if d, ok := e.(Detacher); ok {
		d.Detached(e)
	}
}

func (s *system) rayPick(pointerPos m32.Vec2) (Entity, ray.Result) {
	uiMap := map[*UI][]Entity{}
	uis := []*UI{}
	// Heavy'ish?
	for _, el := range s.elems {
		u := RootUI(el)
		if u == nil {
			continue
		}
		if _, ok := uiMap[u]; !ok {
			uis = append(uis, u)
		}
		uiMap[u] = append(uiMap[u], el)
	}
	sort.Sort(uiSorter(uis))
	for _, u := range uis {
		elems := uiMap[u]

		r := ray.FromScreen(s.gorge.ScreenSize(), u.Camera, pointerPos)

		for i := len(elems) - 1; i >= 0; i-- {
			k := elems[i]
			bw := k.Element()
			if bw.DisableRaycast {
				continue
			}
			t := k.RectTransform()
			rect := t.Rect()

			m := t.Mat4()
			v0 := m.MulV4(m32.Vec4{rect[0], rect[1], 0, 1}).Vec3()
			v1 := m.MulV4(m32.Vec4{rect[2], rect[1], 0, 1}).Vec3() // right
			v2 := m.MulV4(m32.Vec4{rect[0], rect[3], 0, 1}).Vec3() // up)

			if res := ray.IntersectRect(r, v0, v1, v2); res.Hit {
				return k, res
			}
		}
	}
	return nil, ray.Result{}
}

func (s *system) rayTest(pointerPos m32.Vec2) (Entity, ray.Result) {
	hit, r := s.rayPick(pointerPos)
	// debug
	if s.Debug&DebugRays != 0 {
		s.debugRayIntersection(pointerPos)
	}

	if hit == s.pointOver {
		return hit, r
	}
	enterOn := hit
	exitFrom := s.pointOver
	s.pointOver = hit

	var commonParent Entity

	EachParent(exitFrom, func(exitEl Entity) bool {
		EachParent(enterOn, func(enterEl Entity) bool {
			if exitEl == enterEl {
				commonParent = exitEl
				return false
			}
			return true
		})
		return commonParent == nil
	})

	{
		p := &PointerData{RayResult: r, Target: exitFrom}
		EachParent(exitFrom, func(exitEl Entity) bool {
			if exitEl == commonParent {
				return false
			}
			triggerOn(exitEl, EventPointerLeave{p})
			return !p.stopPropagation
		})
	}
	{
		p := &PointerData{RayResult: r, Target: enterOn}
		EachParent(enterOn, func(enterEl Entity) bool {
			if enterEl == commonParent {
				return false
			}
			triggerOn(enterEl, EventPointerEnter{p})
			return !p.stopPropagation
		})
	}
	return hit, r
}

func (s *system) debugRects() {
	// Delete everytime for each update
	rs := rand.New(rand.NewSource(1))
	for _, el := range s.elems {
		t := el.RectTransform()
		rect := t.Rect()

		// Get Rect after Transform
		m := t.Mat4()
		v0 := m.MulV4(m32.Vec4{rect[0], rect[1], 0, 1}).Vec3() // left bottom
		v1 := m.MulV4(m32.Vec4{rect[2], rect[1], 0, 1}).Vec3() // right
		v2 := m.MulV4(m32.Vec4{rect[0], rect[3], 0, 1}).Vec3() // up)

		// Clockwise
		L0 := v2.Sub(v0)
		L1 := v1.Sub(v0)

		// Get the Plane "IntersectPlane"
		planePos := v0 // use v0
		planeNorm := L1.Cross(L0).Normalize()

		s.dbg.SetColor(1, 0, 0, 1)
		s.dbg.AddLine(planePos, planePos.Add(planeNorm))

		s.dbg.SetColor(.5+rs.Float32()*.5, rs.Float32(), rs.Float32(), 1)
		v3 := m.MulV4(m32.Vec4{rect[2], rect[3], 0, 1}).Vec3()
		s.dbg.AddRect(v0, v1, v3, v2)

		dot := m32.Vec4{
			rect[0] + (t.Dim[0])*t.Pivot[0],
			rect[1] + (t.Dim[1])*t.Pivot[1],
			0,
			1,
		}
		s.dbg.AddPoint(m.MulV4(dot).Vec3())
	}
}

func (s *system) debugRayIntersection(pointerPos m32.Vec2) {
	uiMap := map[*UI][]Entity{}
	uis := []*UI{}
	// Heavy'ish?
	for _, el := range s.elems {
		u := RootUI(el)
		if u == nil {
			continue
		}
		if _, ok := uiMap[u]; !ok {
			uis = append(uis, u)
		}
		uiMap[u] = append(uiMap[u], el)
	}

	screenSize := s.gorge.ScreenSize()
	for _, u := range uis {
		elems := uiMap[u]
		cam := u.Camera.Camera()

		vp := m32.Vec4{
			cam.Viewport[0] * screenSize[0],
			cam.Viewport[1] * screenSize[1],
			cam.Viewport[2] * screenSize[0],
			cam.Viewport[3] * screenSize[1],
		}
		width := vp[2] - vp[0]
		height := vp[3] - vp[1]
		ndc := m32.Vec4{
			2*pointerPos[0]/width - 1,
			1 - 2*pointerPos[1]/height,
			1, 1,
		}

		camTransform := u.Camera.Transform()
		m := cam.ProjectionWithAspect(width / height)
		m = m.Mul(camTransform.Inv()).Inv()
		dir := m.MulV4(ndc).Vec3()

		// Ray from camera Entity func somewhere
		r := ray.Ray{
			Position:  camTransform.WorldPosition(),
			Direction: dir,
		}
		if cam.ProjectionType == gorge.ProjectionOrtho {
			r.Position = dir
			r.Direction = camTransform.Forward()
		}
		// Ray
		s.dbg.SetColor(1, 0, 0, 1)
		s.dbg.AddLine(
			r.Position,
			r.Direction.Mul(100),
		)
		s.dbg.SetColor(0, 0, 0, 1)

		for i := len(elems) - 1; i >= 0; i-- {
			k := elems[i]
			bw := k.Element()
			if bw.DisableRaycast {
				continue
			}
			t := k.RectTransform()
			rect := t.Rect() // left bottom right top

			m := t.Mat4()
			v0 := m.MulV4(m32.Vec4{rect[0], rect[1], 0, 1}).Vec3()
			v1 := m.MulV4(m32.Vec4{rect[2], rect[1], 0, 1}).Vec3() // right
			v2 := m.MulV4(m32.Vec4{rect[0], rect[3], 0, 1}).Vec3() // up)
			// Clock wise
			L0 := v2.Sub(v0)
			L1 := v1.Sub(v0)

			// Get the Plane "IntersectPlane"
			planePos := v0 // use v0
			planeNorm := L1.Cross(L0).Normalize()

			res := ray.IntersectPlane(r, planeNorm, planePos)
			ipos := res.Position

			s.dbg.SetColor(0, 1, 0, 1)
			s.dbg.AddPoint(ipos)

			vlen1 := v1.Sub(v0)
			plen1 := ipos.Sub(v0)
			edge1 := plen1.Dot(vlen1) / vlen1.Dot(vlen1)

			vlen2 := v2.Sub(v0)
			plen2 := ipos.Sub(v0)
			edge2 := plen2.Dot(vlen2) / vlen2.Dot(vlen2)

			s.dbg.SetColor(0, 1, 0, 1)
			if edge1 < 0 || edge1 > 1 {
				s.dbg.SetColor(1, 0, 0, 1)
			}
			s.dbg.AddLine(v0, v1)
			s.dbg.AddLine(v1, ipos)
			s.dbg.AddLine(ipos, vlen1.Mul(edge1).Add(v0))
			s.dbg.AddPoint(vlen1.Mul(edge1).Add(v0))
			s.dbg.SetColor(0, 1, 0, 1)
			if edge2 < 0 || edge2 > 1 {
				s.dbg.SetColor(1, 0, 0, 1)
			}
			s.dbg.AddLine(v0, v2)
			s.dbg.AddLine(v0, ipos)
			s.dbg.AddLine(ipos, vlen2.Mul(edge2).Add(v0))
			s.dbg.AddPoint(vlen2.Mul(edge2).Add(v0))
		}
	}
}

func (s *system) Dragging() Entity {
	return s.dragging
}
