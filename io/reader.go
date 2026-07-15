package main

// --- Custom io.Reader ---
//
// io.Reader is the most fundamental input interface in Go's stdlib:
//
//	type Reader interface {
//	    Read(p []byte) (n int, err error)
//	}
//
// Any type that implements Read can be passed to io.ReadAll, io.Copy,
// bufio.NewScanner, json.NewDecoder, and virtually every other I/O function.
//
// Contract:
//   - Fill p with up to len(p) bytes; return the count written.
//   - Return io.EOF (alongside the last bytes, or with n=0) when exhausted.
//   - Never return (0, nil) at the end — that signals "no data yet, try again"
//     and will cause callers to loop forever.
//
// When to use: when your data source is not a file, network connection, or
// string — e.g. generated content, in-memory structures, protocol framing.

import (
	"fmt"
	"io"
	"strings"
)

// RepeatReader is an io.Reader that repeats a string a fixed number of times.
// Demonstrates the minimal contract: implement Read(p []byte) (n int, err error).
type RepeatReader struct {
	data   string
	repeat int
	pos    int // byte offset into the fully-expanded string
}

func NewRepeatReader(s string, n int) *RepeatReader {
	return &RepeatReader{data: s, repeat: n}
}

func (r *RepeatReader) Read(p []byte) (int, error) {
	total := len(r.data) * r.repeat
	if r.pos >= total {
		return 0, io.EOF
	}
	// Fill p from the virtual repeated string.
	n := 0
	for n < len(p) && r.pos < total {
		// Index into the single copy via modulo.
		p[n] = r.data[r.pos%len(r.data)]
		r.pos++
		n++
	}
	if r.pos >= total {
		return n, io.EOF
	}
	return n, nil
}

func DemoCustomReader() {
	fmt.Println("=== Custom io.Reader ===")

	// Read in small chunks to show the contract works across multiple calls.
	r := NewRepeatReader("Go! ", 3)
	buf := make([]byte, 5)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			fmt.Printf("read %d bytes: %q\n", n, buf[:n])
		}
		if err == io.EOF {
			break
		}
	}

	// io.ReadAll consumes any io.Reader — our custom type works with stdlib.
	r2 := NewRepeatReader("hello", 2)
	all, _ := io.ReadAll(r2)
	fmt.Printf("ReadAll: %q\n", all)

	// strings.NewReader is the stdlib's simplest Reader — good baseline to compare.
	sr := strings.NewReader("stdlib reader")
	content, _ := io.ReadAll(sr)
	fmt.Printf("strings.Reader: %q\n", content)
}
