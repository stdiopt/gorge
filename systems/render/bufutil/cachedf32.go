package bufutil

// Although the native benchmarks shows equal performance on copying bytes and
// floats, this has more performance on wasm

// Cachedf32 holds a float32 buffer.
type Cachedf32 struct {
	buffer
	data []float32
}

// NewCachedf32 creates a []float32 transfer buf.
func NewCachedf32(b buffer) *Cachedf32 {
	return &Cachedf32{buffer: b}
}

// Init is represented in floats (4 bytes).
func (b *Cachedf32) Init(sz int) {
	b.data = make([]float32, sz)
	b.buffer.Init(sz * 4)
}

// WriteAt write floats at offset in buffer.
func (b *Cachedf32) WriteAt(floats []float32, off int) {
	copy(b.data[off:], floats)
}

// Flush will send data to underlying buffer.
func (b *Cachedf32) Flush() {
	b.buffer.WriteAt(b.data, 0)
	b.buffer.Flush()
}

// Size returns the size in 32 bit floats (4 bytes per unit).
func (b *Cachedf32) Size() int {
	return len(b.data)
}
