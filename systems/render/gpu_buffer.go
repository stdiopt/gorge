package render

import (
	"runtime"

	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/systems/render/gl"
)

// Bufferer interface to handle renderer buffer
type Bufferer interface {
	ID() gl.Buffer
	Init(sz int)
	Destroy()
	Bind()
	WriteAt(data interface{}, offs int)
	Flush()
	Size() int
}

type bufferManager struct {
	gorge *gorge.Context
	count int
}

func newBufferManager(g *gorge.Context) *bufferManager {
	return &bufferManager{gorge: g}
}

func (m *bufferManager) New(target, usage gl.Enum) *buffer {
	b := newBuffer(m, target, usage)

	runtime.SetFinalizer(b, func(b *buffer) {
		m.gorge.RunInMain(func() {
			b.Destroy()
		})
	})

	return b
}

type buffer struct {
	manager *bufferManager
	buf     gl.Buffer

	// Should we have upload target?
	// gl.UNIFORM_BUFFER
	// gl.ARRAY_ARRAY
	// gl.ELEMENT_ARRAY_BUFFER
	target gl.Enum
	usage  gl.Enum

	size int
}

func newBuffer(manager *bufferManager, target, usage gl.Enum) *buffer {
	return &buffer{
		manager: manager,
		usage:   usage,
		target:  target,
	}
}

func (b *buffer) Init(size int) {
	if !gl.IsValid(b.buf) {
		b.buf = gl.CreateBuffer()
		b.manager.count++
	}

	if b.size != size {
		b.size = size
		gl.BindBuffer(b.target, b.buf)
		gl.BufferInit(b.target, size, b.usage)
		gl.BindBuffer(b.target, gl.Null)
	}
}

func (b *buffer) Destroy() {
	if gl.IsValid(b.buf) {
		gl.DeleteBuffer(b.buf)
		b.manager.count--
	}
	runtime.SetFinalizer(b, nil)
}

func (b *buffer) Size() int { return b.size }

func (b *buffer) ID() gl.Buffer { return b.buf }

func (b *buffer) Bind() { gl.BindBuffer(b.target, b.buf) }

func (b *buffer) WriteAt(data interface{}, offs int) {
	if b.size == 0 {
		return
	}
	if !gl.IsValid(b.buf) {
		panic("writeAt on unitialized buffer")
	}
	gl.BindBuffer(b.target, b.buf)
	gl.BufferSubData(b.target, offs, data)
	gl.BindBuffer(b.target, gl.Null)
}

// Does nothing since it Writes directly
func (b *buffer) Flush() {}
