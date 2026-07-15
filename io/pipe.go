package main

import (
	"fmt"
	"io"
	"strings"
)

// DemoPipe shows io.Pipe: a synchronous, in-memory pipe that connects a writer
// to a reader without buffering. The writer blocks until the reader consumes.
//
// Primary use case: feeding data into a function that expects an io.Reader
// when your data source is a function that writes to an io.Writer (e.g.,
// encoding/json.Encoder → http.Request.Body).
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
