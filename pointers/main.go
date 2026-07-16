package main

import "fmt"

func main() {
	// Nil safety: nil receivers, optional pointer fields, nil as sentinel
	DemoNilSafety()
	fmt.Println()

	// Value vs pointer receivers: mutation, interface satisfaction, auto-deref
	DemoReceivers()
	fmt.Println()

	// Aliasing: shared backing arrays in slices, map reference semantics, defensive copy
	DemoAliasing()
	fmt.Println()

	// Escape analysis: stack vs heap, interface boxing, when to prefer values over pointers
	DemoEscape()
	fmt.Println()

	// Loop variable capture: the classic goroutine-in-range gotcha and its fixes
	DemoLoopCapture()
}
