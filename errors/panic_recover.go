package main

import (
	"errors"
	"fmt"
)

// --- panic / recover Pattern ---
//
// panic is for truly unrecoverable states: violated invariants, programmer
// errors (nil dereference, out-of-bounds index), or conditions where
// continuing would corrupt data. It is NOT a substitute for error returns.
//
// recover() stops a panic mid-flight and returns the panic value.
// It only works inside a deferred function called in the same goroutine.
//
// The canonical library pattern: wrap user-supplied callbacks in a
// safeCall helper that converts any panic into an error, so internal
// panics never propagate to the caller as a crash.

// safeCall calls fn and returns any panic it raises as an error.
// This is the standard pattern for libraries that accept callbacks:
// callers should never have to handle panics from library code.
func safeCall(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("panic: %v", v)
			}
		}
	}()
	fn()
	return nil
}

// divide panics on zero denominator — an invariant violation.
func divide(a, b int) int {
	if b == 0 {
		panic(fmt.Errorf("divide: cannot divide %d by zero", a))
	}
	return a / b
}

// parsePositive panics with a plain string when the input violates a precondition.
func parsePositive(n int) int {
	if n <= 0 {
		panic("parsePositive: argument must be > 0")
	}
	return n * 10
}

func DemoPanicRecover() {
	fmt.Println("=== panic / recover ===")

	// Normal case: no panic.
	err := safeCall(func() {
		result := divide(10, 2)
		fmt.Printf("  divide(10, 2) = %d\n", result)
	})
	fmt.Printf("  error: %v\n\n", err)

	// Panic from an error value is recovered and returned as an error.
	err = safeCall(func() {
		_ = divide(10, 0)
	})
	fmt.Printf("  divide(10, 0) error: %v\n", err)
	fmt.Printf("  is error type: %v\n\n", errors.As(err, new(error)))

	// Panic from a plain string is wrapped into an error.
	err = safeCall(func() {
		_ = parsePositive(-5)
	})
	fmt.Printf("  parsePositive(-5) error: %v\n\n", err)

	// Demonstrate that recover without defer does nothing.
	// (Not shown as runnable — it would crash the process.)
	fmt.Println("  note: recover() only works inside a deferred function;")
	fmt.Println("        calling it elsewhere always returns nil.")
}
