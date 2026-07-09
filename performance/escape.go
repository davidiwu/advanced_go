package main

import "fmt"

// --- Escape Analysis Pattern ---
//
// The Go compiler decides at compile time whether a variable lives on the
// stack (cheap, no GC) or the heap (requires allocation, GC pressure).
// A variable "escapes" to the heap when its lifetime may outlive the
// function call — typically when its address is returned, stored in an
// interface, or captured by a closure that outlives the frame.
//
// You can observe escape decisions with:
//   go build -gcflags="-m" ./performance/
//
// The goal is not to avoid the heap entirely, but to understand when
// unnecessary escapes occur and whether eliminating them matters.

// stackAlloc returns an int, not a pointer. The local variable stays on
// the stack — no heap allocation, no GC involvement.
func stackAlloc(n int) int {
	x := n * 2 // x lives on the stack
	return x   // we return the value, not the address
}

// heapAlloc returns a pointer, so x must escape to the heap — the caller
// holds a reference that outlives this stack frame.
func heapAlloc(n int) *int {
	x := n * 2 // x escapes: address outlives the function
	return &x
}

// noEscape takes a pointer parameter but never lets it escape — it only
// reads through the pointer within the function. The compiler can prove
// the pointed-to value doesn't escape from here.
func noEscape(p *int) int {
	return *p + 1
}

// interfaceEscape passes a concrete value through an interface, which forces
// the value onto the heap because interfaces hold a pointer to a copy.
func interfaceEscape(n int) any {
	return n // n is boxed into an interface → heap allocation
}

// closureEscape shows that a closed-over variable escapes when the closure
// outlives the enclosing function.
func closureEscape(n int) func() int {
	x := n // x escapes: the returned closure captures it
	return func() int { return x }
}

// sum avoids escaping by operating directly on the slice rather than
// creating an interface or returning a pointer.
func sum(nums []int) int {
	total := 0
	for _, v := range nums {
		total += v
	}
	return total
}

func DemoEscape() {
	fmt.Println("=== Escape Analysis ===")

	// Stack allocation: fast, zero GC cost.
	v := stackAlloc(21)
	fmt.Printf("stackAlloc(21) = %d  (stack, no heap)\n", v)

	// Heap allocation: *int means the int must live on the heap.
	p := heapAlloc(21)
	fmt.Printf("heapAlloc(21)  = %d  (heap, via pointer return)\n", *p)

	// noEscape: pointer passed in, value returned — no new heap object.
	local := 10
	fmt.Printf("noEscape(&10)  = %d  (pointer param, no escape)\n", noEscape(&local))

	// Interface boxing: the int gets heap-allocated to fit in an any slot.
	iface := interfaceEscape(42)
	fmt.Printf("interfaceEscape(42) = %v  (heap: value boxed in interface)\n", iface)

	// Closure: captured variable lives on the heap.
	fn := closureEscape(7)
	fmt.Printf("closureEscape(7)()  = %d  (heap: captured by closure)\n", fn())

	// sum: slice header is on the stack; backing array is elsewhere but
	// sum itself introduces no new allocations.
	nums := []int{1, 2, 3, 4, 5}
	fmt.Printf("sum([1..5])         = %d  (no allocs in sum itself)\n", sum(nums))

	fmt.Println()
	fmt.Println("tip: run 'go build -gcflags=\"-m\" ./performance/' to see escape decisions")
}
