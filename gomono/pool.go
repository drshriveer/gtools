package gomono

import (
	"bytes"
	"sync"
)

var pool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

// GetBuffer returns a bytes.Buffer from the pool and a function to put it back.
func GetBuffer() (*bytes.Buffer, func(b *bytes.Buffer)) {
	return pool.Get().(*bytes.Buffer), PutBuffer
}

// PutBuffer returns a bytes.Buffer to the pool.
func PutBuffer(b *bytes.Buffer) {
	b.Reset()
	pool.Put(b)
}
