package main

import "fmt"

func main() {
	// Mutex: exclusive lock for shared state
	DemoMutex()
	fmt.Println()

	// RWMutex: many readers, one writer
	DemoRWMutex()
	fmt.Println()

	// Once: lazy initialization, guaranteed to run exactly once
	DemoOnce()
	fmt.Println()

	// atomic: lock-free integer operations
	DemoAtomic()
	fmt.Println()

	// Channel-as-semaphore: bounded concurrency without a full worker pool
	DemoSemaphore()
}
