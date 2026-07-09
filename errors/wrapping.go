package main

import (
	"errors"
	"fmt"
)

// --- Error Wrapping and Unwrapping Pattern ---
//
// fmt.Errorf("context: %w", err) attaches context to an error while
// preserving the original so errors.Is / errors.As can still reach it.
//
// Each layer in a call stack adds its own context string; the result is
// a chain where the outermost message is most specific and the innermost
// is the root cause.
//
// Design rule: wrap to add context, not to re-package. If you have nothing
// useful to add, return the original error unchanged.

var ErrDiskFull = errors.New("disk full")

// writeChunk is the innermost operation — it returns a raw sentinel.
func writeChunk(path string, size int) error {
	if size > 100 {
		return ErrDiskFull // no wrapping: this IS the root cause
	}
	return nil
}

// saveFile adds file-level context.
func saveFile(path string, data []byte) error {
	if err := writeChunk(path, len(data)); err != nil {
		return fmt.Errorf("saveFile %q: %w", path, err)
	}
	return nil
}

// handleRequest adds request-level context.
func handleRequest(userID int, path string, data []byte) error {
	if err := saveFile(path, data); err != nil {
		return fmt.Errorf("handleRequest user=%d: %w", userID, err)
	}
	return nil
}

func DemoWrapping() {
	fmt.Println("=== Error Wrapping / Unwrapping ===")

	err := handleRequest(42, "/tmp/log", make([]byte, 200))
	if err == nil {
		fmt.Println("  no error")
		return
	}

	fmt.Printf("full message:  %v\n\n", err)

	// errors.Is walks the Unwrap chain to find ErrDiskFull.
	fmt.Printf("errors.Is(err, ErrDiskFull): %v\n", errors.Is(err, ErrDiskFull))

	// Walk the chain manually with errors.Unwrap to show each layer.
	fmt.Println("\nchain (outermost → root):")
	for e := err; e != nil; e = errors.Unwrap(e) {
		fmt.Printf("  %q\n", e)
	}
}
