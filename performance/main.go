package main

import "fmt"

func main() {
	// sync.Pool: reuse temporary objects to reduce allocations and GC pressure
	DemoPool()
	fmt.Println()

	// Preallocation: eliminate slice/map growth reallocations
	DemoPrealloc()
	fmt.Println()

	// strings.Builder: build strings without O(N²) intermediate allocations
	DemoStringBuilder()
	fmt.Println()

	// Struct alignment: reorder fields to eliminate compiler padding
	DemoStructAlign()
	fmt.Println()

	// Escape analysis: understand stack vs heap allocation decisions
	DemoEscape()
}
