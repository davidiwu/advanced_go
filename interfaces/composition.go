package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// --- Interface Composition ---
//
// Go interfaces can embed other interfaces, building larger contracts from
// smaller ones. This mirrors how io.ReadWriter embeds io.Reader and io.Writer:
// each piece is useful on its own, and you can compose them when you need both.
//
// Prefer small, focused interfaces. A function that only needs to read should
// accept io.Reader, not io.ReadWriter — this makes it easier to test and reuse.

// Logger and Flusher are narrow single-method interfaces.
// Each can be used independently wherever only that behavior is needed.
type Logger interface {
	Log(msg string)
}

type Flusher interface {
	Flush() error
}

// BufferedLogger composes both. Functions that need both behaviors accept this.
type BufferedLogger interface {
	Logger
	Flusher
}

// InMemoryLogger satisfies BufferedLogger without any declaration.
type InMemoryLogger struct {
	buf []string
}

func (l *InMemoryLogger) Log(msg string)  { l.buf = append(l.buf, msg) }
func (l *InMemoryLogger) Flush() error   {
	fmt.Println("  flushing", len(l.buf), "log entries:")
	for _, m := range l.buf {
		fmt.Println("   ", m)
	}
	l.buf = l.buf[:0]
	return nil
}

// logSomething only needs to write — it accepts the narrow Logger interface
// so it works with any Logger, not just a BufferedLogger.
func logSomething(l Logger, msg string) {
	l.Log(msg)
}

// flushAndReport needs both behaviors, so it accepts BufferedLogger.
func flushAndReport(l BufferedLogger) {
	if err := l.Flush(); err != nil {
		fmt.Println("  flush error:", err)
	}
}

// DemoComposition demonstrates io.Reader/Writer composition alongside
// the custom Logger/Flusher example.
func DemoComposition() {
	fmt.Println("=== Interface Composition ===")

	l := &InMemoryLogger{}
	logSomething(l, "request started")
	logSomething(l, "processing item 1")
	logSomething(l, "processing item 2")
	flushAndReport(l)

	// The standard library's io interfaces are the canonical example.
	// strings.NewReader satisfies io.Reader; bytes.Buffer satisfies io.ReadWriter.
	fmt.Println("  io.Reader from string:")
	r := strings.NewReader("hello, interface composition")
	// io.ReadAll only needs io.Reader — it doesn't care about write capability.
	data, _ := io.ReadAll(r)
	fmt.Println(" ", string(data))

	fmt.Println("  bytes.Buffer satisfies io.ReadWriter:")
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "written via io.Writer")
	fmt.Println(" ", buf.String())
}
