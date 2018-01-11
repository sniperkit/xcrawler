package scraper

import (
	"bytes"
	"sync"
)

const depthLimit = 5

// bufferPool maintains byte buffers used to read html content
type bufferPool struct {
	pool sync.Pool
}

// newbufferPool creates a new bufferPool bounded to the given size.
func newbufferPool(size int) *bufferPool {
	var bp bufferPool
	bp.pool.New = func() interface{} {
		return new(bytes.Buffer)
	}
	return &bp
}

// Get gets a Buffer from the bufferPool, or creates a new one if none are
// available in the pool.
func (bp *bufferPool) Get() *bytes.Buffer {
	return bp.pool.Get().(*bytes.Buffer)
}

// Put returns the given Buffer to the bufferPool.
func (bp *bufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	bp.pool.Put(b)
}
