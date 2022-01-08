package bufutil

import "unsafe"

// Although the native benchmarks shows equal performance on copying bytes and
// floats, having a specific Type has more performance on wasm

// Cachedf32 holds a float32 buffer.
type Cached[T number] struct {
	buffer
	data []T
}

func NewCached[T number](b buffer) *Cached[T] {
	return &Cached[T]{
		buffer: b,
	}
}

func (c *Cached[T]) Init(sz int) {
	var t T
	bsz := int(unsafe.Sizeof(t))

	c.data = make([]T, sz)
	c.buffer.Init(sz * bsz)
}

func (c *Cached[T]) WriteAt(data []T, offs int) {
	copy(c.data[offs:], data)
}

func (c *Cached[T]) Flush() {
	c.buffer.WriteAt(AsBytes(c.data), 0)
	c.buffer.Flush()
}

func (c *Cached[T]) Size() int {
	return len(c.data)
}
