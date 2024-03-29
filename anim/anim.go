// Package anim provides basic animation functions for gorge
package anim

import (
	"time"

	"github.com/stdiopt/gorge/math/gm"
)

type Channeler interface {
	Update(float32)
	EndTime() float32
}

// State represents the current anim state.
type State int

// Animation states.
const (
	StateStopped = State(iota)
	StateRunning
	StateFinished
)

// LoopType defines anim or track loop type.
type LoopType int

// Loop type.
const (
	LoopNone = LoopType(iota)
	LoopAlways
	LoopMirror
)

// Animation will track time and sent time to channels.
type Animation struct {
	loop      LoopType
	scale     time.Duration
	startTime time.Time
	curTime   float32
	channels  []Channeler
	state     State

	endfn func()
}

// New returns a new animation.
func New() *Animation {
	return &Animation{}
}

// SetScale set the time scale, defaults to 1 second.
func (a *Animation) SetScale(d time.Duration) {
	a.scale = d
}

// SetLoop sets the looping mode for this track.
func (a *Animation) SetLoop(l LoopType) {
	a.loop = l
}

// SetEnd callback when the animation finishes it will call fn
func (a *Animation) SetEnd(fn func()) {
	a.endfn = fn
}

// State returns the current state of the animation.
func (a *Animation) State() State {
	return a.state
}

// Start animation.
func (a *Animation) Start() {
	if a.scale == 0 {
		a.scale = time.Second
	}
	// Recalc totalTime from tracks regardless duration
	a.startTime = time.Now()
	a.curTime = 0
	a.state = StateRunning
}

// Update using internal timing to update.
func (a *Animation) Update() {
	curDur := time.Since(a.startTime)
	a.curTime = float32(curDur) / float32(a.scale)
	a.update()
}

// UpdateDelta updates with delta time, time is in seconds.
func (a *Animation) UpdateDelta(dt float32) {
	a.curTime += dt * float32(a.scale) / float32(time.Second)
	a.update()
}

// Channel adds a channel to animation.
/*func (a *Animation) Channel(intp Interpolator) *Channel {
	c := &Channel{intp: intp}
	a.channels = append(a.channels, c)
	return c
}*/

// AddChannel adds a channel to the animation.
func (a *Animation) AddChannel(c Channeler) {
	a.channels = append(a.channels, c)
}

func (a *Animation) update() {
	// Go through all channels and check key Times
	// the latest key will mandate where we are in the delta
	var lastTime float32
	for _, cc := range a.channels {
		if t := cc.EndTime(); t > lastTime {
			lastTime = t
		}
	}
	curTime := a.loopTime(a.curTime, lastTime)
	if curTime > lastTime {
		a.state = StateFinished
		if a.endfn != nil {
			a.endfn()
		}
	}
	// Check the loop property and the last channel time
	for _, c := range a.channels {
		c.Update(curTime)
	}
}

func (a *Animation) loopTime(ct, last float32) float32 {
	switch a.loop {
	case LoopAlways:
		ct = gm.Mod(ct, last)
	case LoopMirror:
		ct = gm.Mod(ct, 2*last)
		if ct > last {
			ct = last - (ct - last)
			return ct
		}
	}
	return ct
}
