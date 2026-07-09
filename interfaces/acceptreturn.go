package main

import (
	"fmt"
	"sort"
)

// --- Accept Interfaces, Return Structs ---
//
// The canonical Go API design principle:
//   - Function PARAMETERS should be interface types — callers pass any
//     compatible implementation, including test doubles.
//   - Return types should be concrete structs — callers get the full type
//     with all methods, and the compiler catches misuse at the call site.
//
// Returning an interface hides the concrete type from the caller and forces
// them to type-assert when they need concrete behavior. It also makes it
// harder to add methods later without breaking the interface contract.

// Sorter is the interface this package depends on — just Sort.
// By depending on the interface, ProcessData works with any sort strategy
// including test fakes that don't actually sort (useful for determinism in tests).
type Sorter interface {
	Sort(data []int)
}

// Result is a concrete struct returned to callers.
// Returning a concrete type lets callers access all fields and future
// methods without needing a type assertion.
type Result struct {
	Sorted []int
	Count  int
	Min    int
	Max    int
}

func (r Result) String() string {
	return fmt.Sprintf("count=%d min=%d max=%d sorted=%v", r.Count, r.Min, r.Max, r.Sorted)
}

// StdSorter is the production implementation of Sorter.
type StdSorter struct{}

func (s StdSorter) Sort(data []int) { sort.Ints(data) }

// ReverseSorter is an alternative implementation — swappable without
// changing ProcessData because it satisfies the same interface.
type ReverseSorter struct{}

func (s ReverseSorter) Sort(data []int) {
	sort.Sort(sort.Reverse(sort.IntSlice(data)))
}

// ProcessData accepts an interface (flexible) and returns a struct (concrete).
// The caller decides which sort strategy to inject.
func ProcessData(s Sorter, data []int) Result {
	if len(data) == 0 {
		return Result{}
	}
	// Work on a copy so we don't mutate the caller's slice.
	cp := make([]int, len(data))
	copy(cp, data)
	s.Sort(cp)

	return Result{
		Sorted: cp,
		Count:  len(cp),
		Min:    cp[0],
		Max:    cp[len(cp)-1],
	}
}

// DemoAcceptReturn shows the same ProcessData function used with two
// different Sorter implementations, and highlights how returning a
// concrete Result lets callers access all fields directly.
func DemoAcceptReturn() {
	fmt.Println("=== Accept Interfaces, Return Structs ===")

	data := []int{5, 2, 8, 1, 9, 3}

	r1 := ProcessData(StdSorter{}, data)
	fmt.Println("  ascending: ", r1)

	r2 := ProcessData(ReverseSorter{}, data)
	fmt.Println("  descending:", r2)

	// Because Result is concrete, callers access fields without assertions.
	fmt.Printf("  max of ascending run: %d\n", r1.Max)
}
