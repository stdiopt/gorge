package anim

// Ptr helper that returns a func that when called manipulates the pointer.
// ex NewChannel(anim.Vec3).On(anim.Ptr(&t.Position))
func Ptr[T any](p *T) func(T) {
	return func(v T) {
		*p = v
	}
}

// InterpolatorFunc type of func to interpolate a channel.
type InterpolatorFunc[T any] func(a, b T, dt float32) T

// Channel provides a way to interpolate between two values.
type Channel[T any] struct {
	intp  InterpolatorFunc[T]
	keys  []*Key[T]
	value T
	on    func(T)
}

// NewChannel creaates a new channel with type T.
func NewChannel[T any](intp InterpolatorFunc[T]) *Channel[T] {
	return &Channel[T]{intp: intp}
}

// AddChannel creates and adds the channel to Animation.
func AddChannel[T any](a *Animation, intp InterpolatorFunc[T]) *Channel[T] {
	c := &Channel[T]{intp: intp}
	a.AddChannel(c)
	return c
}

func (c *Channel[T]) On(fn func(T)) {
	c.on = fn
}

// EndTime returns the end time for the channel.
func (c *Channel[T]) EndTime() float32 {
	if len(c.keys) == 0 {
		return 0
	}
	return c.keys[len(c.keys)-1].time
}

func (c *Channel[T]) update(v T) T {
	c.value = v
	if c.on != nil {
		c.on(v)
	}
	return v
}

func (c *Channel[T]) UpdateValue(curTime float32) T {
	if len(c.keys) == 0 {
		var z T
		return c.update(z)
	}
	curKey := c.keys[0]
	nextKey := c.keys[0]
	for _, k := range c.keys {
		curKey = nextKey
		nextKey = k
		if k.time >= curTime {
			break
		}
	}
	if curTime < curKey.time { // clamp
		return c.update(c.intp(curKey.val, curKey.val, 0))
	}
	if curTime > nextKey.time { // clamp
		return c.update(c.intp(nextKey.val, nextKey.val, 1))
	}
	normTime := float32(0)
	keyDur := nextKey.time - curKey.time
	if keyDur > 0 {
		normTime = (curTime - curKey.time) / keyDur
	}
	if nextKey.easeFn != nil {
		normTime = nextKey.easeFn(normTime)
	}
	return c.update(c.intp(curKey.val, nextKey.val, normTime))
}

// Update triggers the update and calls the key interpolators for the channel.
func (c *Channel[T]) Update(curTime float32) {
	c.UpdateValue(curTime)
}

// SetKey sets the channel key with the specific value v.
func (c *Channel[T]) SetKey(ct float32, v T) *Key[T] {
	kf := &Key[T]{time: ct, val: v}
	for i, k := range c.keys {
		kt := k.time
		if ct == kt {
			k.val = v
			return k
		}
		if ct < kt {
			c.keys = append(c.keys[:i+1], c.keys[i:]...)
			c.keys[i] = kf
			return kf
		}
	}

	c.keys = append(c.keys, kf)
	return kf
}

func (C *Channel[T]) SetKeys(m map[float32]T) {
	for k, v := range m {
		C.SetKey(k, v)
	}
}

// Key is the animation key on a animation channel.
type Key[T any] struct {
	val    T
	time   float32
	easeFn func(float32) float32
}

// SetEase will set the key easing, the ease will work based on next Key
func (k *Key[T]) SetEase(fn func(float32) float32) {
	k.easeFn = fn
}
