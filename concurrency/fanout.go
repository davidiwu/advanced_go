package main

import (
	"fmt"
	"sync"
)

// --- Fan-out / Fan-in Pattern ---
//
// Fan-out: spread work from one input channel across N workers
//          running in parallel, each reading from the same channel.
// Fan-in:  merge N output channels into a single channel so the
//          consumer only has to range over one place.
//
//            ┌─ worker 0 ─┐
//  jobs ─────┤─ worker 1 ─├───► merged results
//            └─ worker 2 ─┘

// slowSquare simulates an expensive computation so that running
// multiple workers in parallel produces a visible speedup.
func slowSquare(n int) int {
	return n * n
}

// fanOutWorker reads jobs from the shared 'in' channel and sends
// results to its own private output channel.
// Multiple workers can safely read from the same channel because
// the Go channel implementation is concurrent-safe.
func fanOutWorker(id int, in <-chan int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for n := range in {
			result := slowSquare(n)
			out <- fmt.Sprintf("worker %d: %d² = %d", id, n, result)
		}
	}()
	return out
}

// fanIn merges an arbitrary number of channels into one.
// It launches a goroutine per source channel that forwards every
// value, and closes the merged channel once all sources are drained.
func fanIn(channels ...<-chan string) <-chan string {
	merged := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(channels))

	// One forwarder goroutine per source channel.
	for _, ch := range channels {
		go func(c <-chan string) {
			defer wg.Done()
			for v := range c {
				merged <- v
			}
		}(ch)
	}

	// Close the merged channel once all forwarders finish.
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

// DemoFanOut demonstrates fanning out 8 jobs across 3 workers,
// then fanning in their results into a single channel.
// Note: output order is non-deterministic because workers run in parallel.
func DemoFanOut() {
	fmt.Println("=== Fan-out / Fan-in ===")

	jobs := make(chan int, 8) // buffered so the producer never blocks
	for i := 1; i <= 8; i++ {
		jobs <- i
	}
	close(jobs) // no more jobs; workers will drain and exit

	// Fan-out: start 3 workers, all reading from the same jobs channel.
	numWorkers := 3
	workerChans := make([]<-chan string, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workerChans[i] = fanOutWorker(i, jobs)
	}

	// Fan-in: merge all worker output into one channel.
	for result := range fanIn(workerChans...) {
		fmt.Println(result)
	}
}
