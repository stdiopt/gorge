package bufutil

// Cached buffer allow to write data locally and call flush to sync.
type Cached struct {
	buffer
	data []byte
}

// NewCached returns a new Cached buffer that writes to b on flush.
func NewCached(b buffer) *Cached {
	return &Cached{buffer: b}
}

// Init initializes a buffer with specific byte size.
func (b *Cached) Init(sz int) {
	b.data = make([]byte, sz)
	b.buffer.Init(sz)
}

// WriteAt writes data at offset.
func (b *Cached) WriteAt(data interface{}, offs int) {
	copy(b.data[offs:], asBytes(data))
}

// Flush writes the buffer to the underlying buffer controller.
func (b *Cached) Flush() {
	b.buffer.WriteAt(b.data, 0)
	b.buffer.Flush()
}
