package main

import (
	"fmt"
	"sync"
)

// --- sync.Pool Pattern ---
//
// sync.Pool holds a set of temporary objects that can be reused across
// goroutines to reduce allocations and GC pressure. The pool is safe
// for concurrent use. Objects may be evicted by the GC at any time,
// so never store state you can't recreate.

// byteSlice is the object we pool. In practice this would be a
// buffer, parser context, or any object expensive to allocate.
type byteSlice []byte

var bufPool = sync.Pool{
	// New is called when the pool is empty.
	New: func() any {
		b := make(byteSlice, 0, 64)
		return &b
	},
}

// processWithPool borrows a buffer from the pool, uses it, then
// returns it so the next caller can reuse it without a fresh allocation.
func processWithPool(data string) string {
	bp := bufPool.Get().(*byteSlice)
	buf := (*bp)[:0] // reset length, keep capacity

	buf = append(buf, "processed: "...)
	buf = append(buf, data...)

	result := string(buf)

	*bp = buf
	bufPool.Put(bp) // return to pool for reuse
	return result
}

// processWithoutPool allocates a fresh slice every call — baseline for comparison.
func processWithoutPool(data string) string {
	buf := make([]byte, 0, 64)
	buf = append(buf, "processed: "...)
	buf = append(buf, data...)
	return string(buf)
}

func DemoPool() {
	fmt.Println("=== sync.Pool ===")

	inputs := []string{"alpha", "beta", "gamma", "delta"}

	fmt.Println("with pool (reuses buffer across calls):")
	for _, s := range inputs {
		fmt.Println(" ", processWithPool(s))
	}

	fmt.Println("without pool (fresh allocation each call):")
	for _, s := range inputs {
		fmt.Println(" ", processWithoutPool(s))
	}

	// Concurrency-safe: multiple goroutines can Get/Put simultaneously.
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			result := processWithPool(fmt.Sprintf("goroutine-%d", n))
			_ = result
		}(i)
	}
	wg.Wait()
	fmt.Println("concurrent pool usage: no data races")
}
