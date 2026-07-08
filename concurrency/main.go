package main

import "fmt"

func main() {
	// Pipeline: chain goroutine stages via channels
	DemoPipeline()
	fmt.Println()

	// Fan-out / Fan-in: spread work across N workers, merge results
	DemoFanOut()
	fmt.Println()

	// Worker Pool: bounded concurrency via a fixed goroutine pool
	DemoWorkerPool()
	fmt.Println()

	// errgroup: coordinate goroutines with automatic error propagation
	DemoErrGroupNoError()
	fmt.Println()
	DemoErrGroup()
	fmt.Println()

	// Context: timeouts, cancellation, and request-scoped values
	DemoContextTimeout()
	fmt.Println()
	DemoContextCancel()
	fmt.Println()
	DemoContextValues()
}
