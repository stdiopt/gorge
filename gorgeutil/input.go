package gorgeutil

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/systems/input"
)

// KeyDirection returns a vec2 based on the direction of the key
// - Boost
// - Top, Right, Bottom, Left
func KeyDirection(
	g *gorge.Context,
	boost, top, right, bottom, left input.Key,
) func(gorge.EventUpdate) gm.Vec2 {
	ic := input.FromContext(g)
	f := gm.Vec2{}
	return func(e gorge.EventUpdate) gm.Vec2 {
		power := float32(1)
		s := float32(.4)
		if ic.KeyDown(boost) {
			power *= 10
		}
		if ic.KeyDown(top) {
			f = f.Add(gm.Vec2{0, s})
		}
		if ic.KeyDown(right) {
			f = f.Add(gm.Vec2{s, 0})
		}
		if ic.KeyDown(bottom) {
			f = f.Add(gm.Vec2{0, -s})
		}
		if ic.KeyDown(left) {
			f = f.Add(gm.Vec2{-s, 0})
		}
		f = f.Clamp(gm.Vec2{-1, -1}, gm.Vec2{1, 1})
		f = f.Lerp(gm.Vec2{}, e.DeltaTime()*10)
		return f.Mul(power)
	}
}
