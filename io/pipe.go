package main

// --- io.Pipe ---
//
// io.Pipe creates a synchronous, in-memory connection between a writer and a
// reader with no intermediate buffer. Every Write blocks until a corresponding
// Read has consumed the data — they rendezvous directly.
//
// Pattern:
//
//	pr, pw := io.Pipe()
//	go func() {
//	    defer pw.Close()       // signals EOF to the reader
//	    pw.Write(data)
//	}()
//	io.ReadAll(pr)             // blocks until goroutine closes pw
//
// When to use: connecting a function that writes (json.Encoder, gzip.Writer,
// tar.Writer) to a function that reads (http.Request.Body, io.ReadAll) without
// allocating an intermediate bytes.Buffer. The pipe avoids buffering the entire
// payload in memory.
//
// Key points:
//   - Writer and reader MUST run in separate goroutines; same-goroutine use deadlocks.
//   - pw.CloseWithError(err) propagates the error to the reader's next Read call.
//   - pr.CloseWithError(err) propagates the error to the writer's next Write call.

import (
	"fmt"
	"io"
	"strings"
)

// DemoPipe shows io.Pipe: a synchronous, in-memory pipe that connects a writer
// to a reader without buffering. The writer blocks until the reader consumes.
func DemoPipe() {
	fmt.Println("=== io.Pipe ===")

	pr, pw := io.Pipe()

	// Writer goroutine: produce data and close the pipe when done.
	go func() {
		defer pw.Close()
		for _, line := range []string{"line 1\n", "line 2\n", "line 3\n"} {
			fmt.Fprint(pw, line)
		}
	}()

	// Reader side runs in the current goroutine.
	data, err := io.ReadAll(pr)
	fmt.Printf("received: %q err: %v\n", data, err)

	// Error propagation: CloseWithError passes an error to the reader.
	pr2, pw2 := io.Pipe()
	go func() {
		fmt.Fprint(pw2, "partial data")
		pw2.CloseWithError(fmt.Errorf("upstream failure"))
	}()
	buf := new(strings.Builder)
	_, err = io.Copy(buf, pr2)
	fmt.Printf("after error: %q err: %v\n", buf.String(), err)
}
