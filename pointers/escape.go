package main

import "fmt"

// --- Escape Analysis ---
//
// The Go compiler decides at compile time whether a variable lives on the stack
// or the heap — a process called escape analysis. Stack allocation is essentially
// free (bump a pointer); heap allocation involves the garbage collector and is
// measurably slower on hot paths.
//
// A variable "escapes to the heap" when:
//   - Its address is returned from the function that created it.
//   - It is stored in a location that outlives the current stack frame
//     (e.g. a struct field, a slice element, an interface value, a channel).
//   - It is too large to fit on the stack (default limit: 64 KB).
//
// Returning a pointer to a local IS safe in Go (the compiler heap-allocates it
// automatically), but it costs an allocation. On a hot path that matters.
//
// To inspect escape decisions:
//   go build -gcflags="-m" ./pointers/
//
// Key points:
//   - Prefer returning values (not pointers) for small structs — copy is cheap
//     and keeps the value on the stack.
//   - Interface boxing always escapes: storing a concrete value in an interface{}
//     or any interface causes a heap allocation for the value.
//   - sync.Pool can amortise heap allocations for short-lived objects.

// smallPair is cheap to copy — returning by value keeps it on the stack.
type smallPair struct {
	x, y int
}

// byValue — compiler keeps the pair on the stack; no allocation.
func newPairByValue(x, y int) smallPair {
	return smallPair{x, y}
}

// byPointer — pair escapes to the heap because its address is returned.
func newPairByPointer(x, y int) *smallPair {
	p := smallPair{x, y} // escapes to heap
	return &p
}

// interfaceEscape shows that boxing a value into an interface causes it to escape.
func boxInInterface(v int) interface{} {
	return v // v escapes to heap here
}

// largeStruct illustrates that large values are worth passing by pointer to avoid
// copying on every call — but note that taking the address causes a heap allocation
// if the struct was created locally.
type largeStruct struct {
	data [1024]byte
}

// processLarge accepts a pointer to avoid copying 1 KB on every call.
func processLarge(s *largeStruct) int {
	sum := 0
	for _, b := range s.data {
		sum += int(b)
	}
	return sum
}

func DemoEscape() {
	fmt.Println("=== Escape Analysis ===")

	// Stack allocation — no GC pressure
	p1 := newPairByValue(1, 2)
	fmt.Printf("value pair: %+v (stack allocated)\n", p1)

	// Heap allocation — returned pointer forces escape
	p2 := newPairByPointer(3, 4)
	fmt.Printf("pointer pair: %+v (heap allocated)\n", *p2)

	// Interface boxing always allocates
	boxed := boxInInterface(42)
	fmt.Printf("boxed value: %v (heap allocated)\n", boxed)

	// Large struct: pass by pointer to avoid copy cost
	ls := largeStruct{}
	for i := range ls.data {
		ls.data[i] = byte(i % 256)
	}
	fmt.Printf("processLarge sum: %d\n", processLarge(&ls))

	fmt.Println("Run: go build -gcflags=\"-m\" ./pointers/ to see escape decisions")
}
