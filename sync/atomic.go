package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// --- sync/atomic ---
//
// The atomic package provides lock-free operations on integers and pointers.
// Operations like Add, Load, Store, and CompareAndSwap are guaranteed to
// complete without being interrupted by another goroutine — no mutex needed.
//
// Use atomic when:
//   - You only need to increment/decrement a single integer counter
//   - You need the absolute lowest overhead (atomic ops are ~10x faster than Mutex)
//
// For anything more complex (protecting a struct, multiple fields together),
// stick with a Mutex — atomic is easy to misuse.

func DemoAtomic() {
	fmt.Println("=== sync/atomic ===")

	var hits atomic.Int64 // zero value is ready to use (Go 1.19+)
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hits.Add(1) // atomic — no mutex, no data race
		}()
	}
	wg.Wait()
	fmt.Printf("hits (want 1000): %d\n", hits.Load())

	// CompareAndSwap: only store new value if current value matches expected.
	// Classic use-case: implement a "set once" flag.
	var flag atomic.Int32
	swapped := flag.CompareAndSwap(0, 1)
	fmt.Printf("first CAS swapped: %v, flag: %d\n", swapped, flag.Load())
	swapped = flag.CompareAndSwap(0, 1) // fails — value is now 1, not 0
	fmt.Printf("second CAS swapped: %v, flag: %d\n", swapped, flag.Load())
}
