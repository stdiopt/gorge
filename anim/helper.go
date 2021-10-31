package anim

type channel = Channel

// Channelf32 channel handler for f32.
type Channelf32 struct {
	channel
}

// NewChannelf32 returns a New Channel32
func NewChannelf32(f *float32) *Channelf32 {
	return &Channelf32{*NewChannel(Float32(f))}
}

// SetKey sets the key for the channel.
func (c *Channelf32) SetKey(t float32, v float32) *Key {
	return c.channel.SetKey(t, v)
}

// ChannelFuncf32 channel handler for func(float32).
type ChannelFuncf32 struct {
	channel
}

// NewChannelFuncf32 helper to create a channel with underlying Func32 interpolator.
func NewChannelFuncf32(fn func(v float32)) *ChannelFuncf32 {
	return &ChannelFuncf32{
		*NewChannel(Funcf32(fn)),
	}
}

// SetKey sets the key for the Channel.
func (c *ChannelFuncf32) SetKey(t float32, v float32) *Key {
	return c.channel.SetKey(t, v)
}
