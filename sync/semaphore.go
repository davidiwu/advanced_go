package main

import (
	"fmt"
	"sync"
	"time"
)

// --- Channel-as-Semaphore ---
//
// A buffered channel of size N acts as a counting semaphore that limits
// how many goroutines may do something at the same time.
//
// Pattern:
//   sem := make(chan struct{}, N)  // N = max concurrent workers
//   sem <- struct{}{}             // acquire: blocks when N slots are taken
//   defer func() { <-sem }()      // release: always free the slot on exit
//
// This is useful when you fan out to many goroutines but want to cap the
// number performing an expensive operation (e.g. outbound HTTP calls,
// file I/O) at any one moment — without a full worker-pool.

func DemoSemaphore() {
	fmt.Println("=== Channel-as-Semaphore ===")

	const total = 10
	const maxConcurrent = 3

	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			sem <- struct{}{} // acquire slot
			defer func() { <-sem }()

			// Simulate a slow operation — at most 3 goroutines are here at once.
			fmt.Printf("  worker %2d started  (slots in use: %d/%d)\n",
				id, len(sem), maxConcurrent)
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("  worker %2d finished\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("all workers done")
}
