package main

import (
	"fmt"
	"io"
)

// CountingWriter wraps an io.Writer and tracks how many bytes have been written.
// A common pattern for progress reporting or quota enforcement.
type CountingWriter struct {
	w     io.Writer
	total int64
}

func NewCountingWriter(w io.Writer) *CountingWriter {
	return &CountingWriter{w: w}
}

func (c *CountingWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.total += int64(n)
	return n, err
}

func (c *CountingWriter) BytesWritten() int64 { return c.total }

// DiscardWriter silently drops all bytes — useful in tests or benchmarks.
// (The stdlib provides io.Discard for this; shown here for illustration.)
type DiscardWriter struct{ total int64 }

func (d *DiscardWriter) Write(p []byte) (int, error) {
	d.total += int64(len(p))
	return len(p), nil
}

func DemoCustomWriter() {
	fmt.Println("=== Custom io.Writer ===")

	// CountingWriter: track bytes through Fprintf/write calls.
	cw := NewCountingWriter(io.Discard) // write destination doesn't matter here
	fmt.Fprintf(cw, "hello, %s!\n", "world")
	fmt.Fprintf(cw, "line two\n")
	fmt.Printf("bytes written: %d\n", cw.BytesWritten())

	// DiscardWriter: measure how many bytes would be written.
	dw := &DiscardWriter{}
	for range 4 {
		fmt.Fprintf(dw, "chunk of data\n")
	}
	fmt.Printf("discarded bytes: %d\n", dw.total)

	// io.MultiWriter fans out writes to multiple destinations simultaneously.
	buf1 := &DiscardWriter{}
	buf2 := &DiscardWriter{}
	mw := io.MultiWriter(buf1, buf2)
	fmt.Fprintf(mw, "written to both\n")
	fmt.Printf("buf1: %d, buf2: %d\n", buf1.total, buf2.total)
}
