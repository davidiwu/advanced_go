package main

import (
	"fmt"
	"sync"
)

// --- sync.Mutex ---
//
// A Mutex (mutual exclusion lock) protects a shared resource so that
// only one goroutine can access it at a time. Use it whenever multiple
// goroutines read AND write the same variable — a data race otherwise.
//
// Rule of thumb: lock before touching shared state, unlock immediately after.
// defer mu.Unlock() is the idiomatic way to ensure the lock is always released,
// even if the function returns early or panics.

type safeCounter struct {
	mu    sync.Mutex
	value int
}

// Inc increments the counter safely from any goroutine.
func (c *safeCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// Value reads the counter safely.
func (c *safeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func DemoMutex() {
	fmt.Println("=== sync.Mutex ===")

	counter := &safeCounter{}
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Inc()
		}()
	}

	wg.Wait()
	// Without the mutex this would be a data race and the result unpredictable.
	fmt.Printf("final count (want 1000): %d\n", counter.Value())
}
