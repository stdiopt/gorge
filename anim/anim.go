// Package anim provides basic animation functions for gorge
package anim

import (
	"time"

	"github.com/stdiopt/gorge/m32"
)

type channeler interface {
	Channel() *Channel
}

// Interpolator interface for accepting interpolator structs.
type Interpolator interface {
	Interpolate(a, b interface{}, dt float32)
}

// InterpolatorFunc a func that implements the Interpolator interface.
type InterpolatorFunc func(p, n interface{}, dt float32)

// Interpolate implements the interpolator func.
func (fn InterpolatorFunc) Interpolate(a, b interface{}, dt float32) { fn(a, b, dt) }

// State represents the current anim state.
type State int

// Animation states.
const (
	StateStopped = State(iota)
	StateRunning
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
	channels  []channeler

	state State
}

// NewAnimation returns a new animation.
func NewAnimation() *Animation {
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
func (a *Animation) Channel(intp Interpolator) *Channel {
	c := &Channel{intp: intp}
	a.channels = append(a.channels, c)
	return c
}

// AddChannel adds a channel to the animation.
func (a *Animation) AddChannel(c channeler) {
	a.channels = append(a.channels, c)
}

func (a *Animation) update() {
	// Go through all channels and check key Times
	// the latest key will mandate where we are in the delta
	var lastTime float32
	for _, cc := range a.channels {
		c := cc.Channel()
		if len(c.keys) == 0 {
			continue
		}
		if t := c.keys[len(c.keys)-1].time; t > lastTime {
			lastTime = t
		}
	}
	curTime := a.loopTime(a.curTime, lastTime)
	// Check the loop property and the last channel time
	for _, c := range a.channels {
		c.Channel().Update(curTime)
	}
}

func (a *Animation) loopTime(ct, last float32) float32 {
	switch a.loop {
	case LoopAlways:
		ct = m32.Mod(ct, last)
	case LoopMirror:
		ct = m32.Mod(ct, 2*last)
		if ct > last {
			ct = last - (ct - last)
			return ct
		}
	}
	return ct
}
