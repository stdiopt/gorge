package anim

// InterpolatorFunc type of func to interpolate a channel.
type InterpolatorFunc[T any] func(a, b T, dt float32)

// Channel provides a way to interpolate between two values.
type Channel[T any] struct {
	intp InterpolatorFunc[T]
	keys []*Key[T]
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

// EndTime returns the end time for the channel.
func (c *Channel[T]) EndTime() float32 {
	if len(c.keys) == 0 {
		return 0
	}
	return c.keys[len(c.keys)-1].time
}

// Update triggers the update and calls the key interpolators for the channel.
func (c *Channel[T]) Update(curTime float32) {
	if len(c.keys) == 0 {
		return
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
		c.intp(curKey.val, curKey.val, 0)
		return
	}
	if curTime > nextKey.time { // clamp
		c.intp(nextKey.val, nextKey.val, 1)
		return
	}
	normTime := float32(0)
	keyDur := nextKey.time - curKey.time
	if keyDur > 0 {
		normTime = (curTime - curKey.time) / keyDur
	}
	if nextKey.easeFn != nil {
		normTime = nextKey.easeFn(normTime)
	}
	c.intp(curKey.val, nextKey.val, normTime)
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
