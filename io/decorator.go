package main

import (
	"fmt"
	"io"
	"strings"
)

// UpperReader wraps an io.Reader and uppercases every ASCII letter it returns.
// This is the decorator pattern: same Read interface, different behavior.
type UpperReader struct {
	r io.Reader
}

func NewUpperReader(r io.Reader) *UpperReader { return &UpperReader{r: r} }

func (u *UpperReader) Read(p []byte) (int, error) {
	n, err := u.r.Read(p)
	for i := range n {
		if p[i] >= 'a' && p[i] <= 'z' {
			p[i] -= 32
		}
	}
	return n, err
}

// LogWriter wraps an io.Writer and logs each Write call's size before forwarding.
type LogWriter struct {
	w   io.Writer
	tag string
}

func NewLogWriter(w io.Writer, tag string) *LogWriter { return &LogWriter{w: w, tag: tag} }

func (l *LogWriter) Write(p []byte) (int, error) {
	fmt.Printf("[%s] writing %d bytes\n", l.tag, len(p))
	return l.w.Write(p)
}

func DemoDecorator() {
	fmt.Println("=== Decorator pattern ===")

	// Stack decorators: UpperReader wraps a strings.Reader.
	base := strings.NewReader("hello from the decorator pattern")
	upper := NewUpperReader(base)
	result, _ := io.ReadAll(upper)
	fmt.Printf("uppercased: %q\n", result)

	// io.LimitReader truncates a reader to N bytes — stdlib decorator.
	limited := io.LimitReader(strings.NewReader("only the first ten characters matter"), 10)
	data, _ := io.ReadAll(limited)
	fmt.Printf("limited:    %q\n", data)

	// io.TeeReader: read from src while simultaneously writing to dst.
	// Useful for hashing or logging data as it streams through.
	var captured strings.Builder
	src := strings.NewReader("tee me")
	tee := io.TeeReader(src, &captured)
	io.ReadAll(tee) // consume
	fmt.Printf("tee captured: %q\n", captured.String())

	// Stacked LogWriter over io.Discard.
	lw := NewLogWriter(io.Discard, "demo")
	fmt.Fprintf(lw, "first write")
	fmt.Fprintf(lw, "second write — longer")
}
