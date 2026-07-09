package main

import "fmt"

func main() {
	// Sentinel errors: package-level var errors compared with errors.Is
	DemoSentinel()
	fmt.Println()

	// Wrapping / unwrapping: add context at each layer, preserve the root cause
	DemoWrapping()
	fmt.Println()

	// Custom error types: structured errors with Is / As / Unwrap methods
	DemoCustom()
	fmt.Println()

	// panic / recover: convert panics to errors at a library boundary
	DemoPanicRecover()
	fmt.Println()

	// errors.Join: aggregate multiple errors from batch or parallel operations
	DemoJoin()
}
