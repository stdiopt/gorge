package gorgeui

import (
	"math/rand"
	"sort"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/core/event"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/math/ray"
	"github.com/stdiopt/gorge/systems/input"
	"github.com/stdiopt/gorge/text"
)

type system struct {
	Debug DebugFlag

	gorge *gorge.Context
	font  *text.Font

	uis map[*UI]struct{}
	// We might not need this and just run through the UI children elements
	elems []Entity

	// Should have these per pointer
	// meaning mouse button[0-5], touch{1...}
	pointOver    Entity
	pointDown    Entity
	pointDownPos gm.Vec2

	dragging Entity

	curMouse   gm.Vec2
	deltaMouse gm.Vec2
	dbg        *debugLines
}

func (s *system) setupEvents(g *gorge.Context) {
	event.Handle(g, s.handlePointer)

	event.Handle(g, func(gorge.EventPreUpdate) {
		if s.Debug != 0 {
			s.dbg.Clear()
		}
	})
	event.Handle(g, func(e gorge.EventPostUpdate) {
		for ui := range s.uis {
			ui.update(e.DeltaTime())
		}
		for _, el := range s.elems {
			triggerOn(el, EventUpdate(e.DeltaTime()))
		}

		if s.Debug&DebugRects != 0 {
			s.debugRects()
		}
	})
	event.Handle(g, func(e gorge.EventAddEntity) {
		if v, ok := e.Entity.(Entity); ok {
			s.addEntity(v)
		}
	})
	event.Handle(g, func(e gorge.EventRemoveEntity) {
		if v, ok := e.Entity.(Entity); ok {
			s.removeEntity(v)
		}
	})
}

func (s *system) addEntity(e Entity) {
	el := e.Element()
	el.Attached = true

	if ui, ok := e.(*UI); ok {
		if s.uis == nil {
			s.uis = map[*UI]struct{}{}
		}
		s.uis[ui] = struct{}{}
		ui.gorge = s.gorge
		return
	}
	s.elems = append(s.elems, e)

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

func (s *system) handlePointer(e input.EventPointer) {
	s.deltaMouse = e.Pointers[0].Pos.Sub(s.curMouse)
	s.curMouse = e.Pointers[0].Pos

	hit, res := s.rayTest(s.curMouse)

	pd := &PointerData{
		RayResult: res,
		Delta:     s.deltaMouse,
		Position:  s.curMouse,
		Wheel:     e.Pointers[0].ScrollDelta,
		Target:    hit,
	}
	curDown := s.pointDown

	if e.Type == input.MouseWheel {
		if hit != nil {
			EachParent(hit, func(e Entity) bool {
				triggerOn(e, EventPointerWheel{pd})
				return !pd.stopPropagation
			})
		}
	}

	if e.Type == input.MouseDown && e.Button == 0 && s.dragging == nil {
		s.pointDown = hit
		s.pointDownPos = s.curMouse
		if hit != nil {
			EachParent(hit, func(e Entity) bool {
				triggerOn(e, EventPointerDown{pd})
				return !pd.stopPropagation
			})
		}
	}
	if e.Type == input.MouseUp && e.Button == 0 {
		if s.pointDown != nil {
			EachParent(s.pointDown, func(e Entity) bool {
				triggerOn(e, EventPointerUp{pd})
				return !pd.stopPropagation
			})
			s.pointDown = nil
		}
		if s.dragging != nil {
			triggerOn(s.dragging, EventDragEnd{pd})
			event.Trigger(s.gorge, EventDragging{nil})
			s.dragging = nil
		}
	}

	// Drag detection
	// I mouse is still down and not nil and pointerDown is still and dragging is nil
	if s.pointDown != nil && curDown == s.pointDown && s.dragging == nil {
		ui := RootUI(s.pointDown)
		d := s.curMouse.Sub(s.pointDownPos).Abs()
		if d[0] > ui.DragThreshold || d[1] > ui.DragThreshold {
			EachParent(s.pointDown, func(e Entity) bool {
				if !e.Element().DragEvents {
					return true
				}
				s.dragging = e
				triggerOn(s.dragging, EventDragBegin{pd})
				event.Trigger(s.gorge, EventDragging{e})
				return false
			})
		}
	} else if s.dragging != nil {
		triggerOn(s.dragging, EventDrag{pd})
	}
}

// NewPick, it might be slower but can overcome masked entities
func (s *system) rayPick(pointerPos gm.Vec2) (Entity, ray.Result) {
	type masker interface{ IsMasked() bool }
	// Organize UI's
	uis := []*UI{}
	for ui := range s.uis {
		uis = append(uis, ui)
	}
	sort.Sort(uiSorter(uis))

	var pick func(r ray.Ray, e Entity) (Entity, ray.Result)

	pick = func(r ray.Ray, e Entity) (Entity, ray.Result) {
		res := rayRect(r, e)
		// if entity is interface we check against it, else we continue checking
		if !res.Hit {
			if m, ok := e.(masker); ok && m.IsMasked() {
				return e, res
			}
		}

		c, ok := e.(gorge.EntityContainer)
		if !ok {
			return e, res
		}
		children := c.GetEntities()
		for i := len(children) - 1; i >= 0; i-- {
			ce := children[i]

			ge, ok := ce.(Entity)
			if !ok { // skip non entity
				continue
			}
			// slip ray cast disabled
			if ge.Element().DisableRaycast {
				continue
			}
			if ret, res := pick(r, ge); res.Hit {
				return ret, res
			}
		}
		return e, res
	}

	for _, u := range uis {
		r := ray.FromScreen(s.gorge.ScreenSize(), u.Camera, pointerPos)
		if ent, res := pick(r, u); res.Hit {
			return ent, res
		}
	}
	return nil, ray.Result{}
}

/*
func (s *system) rayPick(pointerPos gm.Vec2) (Entity, ray.Result) {
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
			if res := rayRect(r, k); res.Hit {
				return k, res
			}
		}
	}
	return nil, ray.Result{}
}

/**/

// Ray test and return an entity
// if the pick is the same to the pointOver it will return
// if not it will trigger Pointer leave on the existing pointOver
// and PointerEnter on the new pick
func (s *system) rayTest(pointerPos gm.Vec2) (Entity, ray.Result) {
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

	// Move this to handleEvents
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

	for ui := range s.uis {
		s.dbg.SetColor(1, 0, 0, 1)
		s.dbg.AddCross(ui.Position, 5)
	}

	for _, el := range s.elems {
		t := el.RectTransform()
		rect := t.Rect()

		// Get Rect after Transform
		m := t.Mat4()
		v0 := m.MulV4(gm.Vec4{rect[0], rect[1], 0, 1}).Vec3() // left bottom
		v1 := m.MulV4(gm.Vec4{rect[2], rect[1], 0, 1}).Vec3() // right
		v2 := m.MulV4(gm.Vec4{rect[0], rect[3], 0, 1}).Vec3() // up)

		// Clockwise
		L0 := v2.Sub(v0)
		L1 := v1.Sub(v0)

		// Get the Plane "IntersectPlane"
		planePos := v0 // use v0
		planeNorm := L1.Cross(L0).Normalize()

		s.dbg.SetColor(1, 0, 0, 1)
		s.dbg.AddLine(planePos, planePos.Add(planeNorm))

		s.dbg.SetColor(
			.5+float32(rs.Float64())*.5,
			float32(rs.Float32()),
			float32(rs.Float32()), 1,
		)
		v3 := m.MulV4(gm.Vec4{rect[2], rect[3], 0, 1}).Vec3()
		s.dbg.AddRect(v0, v1, v3, v2)

		dot := gm.Vec4{
			rect[0] + (t.Dim[0])*t.Pivot[0],
			rect[1] + (t.Dim[1])*t.Pivot[1],
			0,
			1,
		}
		s.dbg.AddCross(m.MulV4(dot).Vec3(), .5)
	}
}

func (s *system) debugRayIntersection(pointerPos gm.Vec2) {
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

		vp := gm.Vec4{
			cam.Viewport[0] * screenSize[0],
			cam.Viewport[1] * screenSize[1],
			cam.Viewport[2] * screenSize[0],
			cam.Viewport[3] * screenSize[1],
		}
		width := vp[2] - vp[0]
		height := vp[3] - vp[1]
		ndc := gm.Vec4{
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
			v0 := m.MulV4(gm.Vec4{rect[0], rect[1], 0, 1}).Vec3()
			v1 := m.MulV4(gm.Vec4{rect[2], rect[1], 0, 1}).Vec3() // right
			v2 := m.MulV4(gm.Vec4{rect[0], rect[3], 0, 1}).Vec3() // up)
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
