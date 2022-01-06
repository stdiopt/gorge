package anim

type channel = Channel

// Channelf32 channel handler for f32.
type Channelf32 struct {
	channel
}

// AddChannelf32 returns a New Channel32
func AddChannelf32(a *Animation, f *float32) *Channelf32 {
	c := &Channelf32{*NewChannel(Float32(f))}
	a.AddChannel(c)
	return c
}

// SetKey sets the key for the channel.
func (c *Channelf32) SetKey(t float32, v float32) *Key {
	return c.channel.SetKey(t, v)
}

// ChannelFuncf32 channel handler for func(float32).
type ChannelFuncf32 struct {
	channel
}

// AddChannelFuncf32 helper to create a channel with underlying Func32 interpolator.
func AddChannelFuncf32(a *Animation, fn func(v float32)) *ChannelFuncf32 {
	c := &ChannelFuncf32{*NewChannel(Funcf32(fn))}
	a.AddChannel(c)
	return c
}

// SetKey sets the key for the Channel.
func (c *ChannelFuncf32) SetKey(t float32, v float32) *Key {
	return c.channel.SetKey(t, v)
}
