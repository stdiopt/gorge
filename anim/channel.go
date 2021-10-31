package anim

// Channel contains keys
type Channel struct {
	intp Interpolator
	keys []*Key
}

// NewChannel returns a new channel with the specific interpolator.
func NewChannel(intp Interpolator) *Channel {
	return &Channel{intp: intp}
}

// Channel returns the channel (to be embed in other structs)
func (c *Channel) Channel() *Channel { return c }

// Update triggers the update and calls the key interpolators for the channel.
func (c *Channel) Update(curTime float32) {
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
		c.intp.Interpolate(curKey.val, curKey.val, 0)
		return
	}
	if curTime > nextKey.time { // clamp
		c.intp.Interpolate(nextKey.val, nextKey.val, 1)
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
	c.intp.Interpolate(curKey.val, nextKey.val, normTime)
}

// SetKey sets the channel key with the specific value v.
func (c *Channel) SetKey(ct float32, v interface{}) *Key {
	kf := &Key{time: ct, val: v}
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
type Key struct {
	val    interface{}           // holds any value
	time   float32               // the current time
	easeFn func(float32) float32 // optional key specific easing function
}

// SetEase will set the key easing, the ease will work based on next Key
func (k *Key) SetEase(fn func(float32) float32) *Key {
	k.easeFn = fn
	return k
}
