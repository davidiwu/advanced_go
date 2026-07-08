package main

import "fmt"

// --- Pipeline Pattern ---
//
// A pipeline is a series of stages connected by channels.
// Each stage runs in its own goroutine, receives values from
// the previous stage, transforms them, and forwards results
// to the next stage.
//
// Data flow:  generate → square → filter → print

// generate sends a fixed set of integers onto a channel,
// then closes it to signal downstream stages there is no more data.
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out) // closing signals the next stage to stop ranging
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// square reads integers from 'in', squares each one,
// and sends the result to a new channel.
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in { // range on a channel stops when the channel is closed
			out <- n * n
		}
	}()
	return out
}

// filterEven reads integers from 'in' and forwards only even values.
// This illustrates that a stage can drop items, not just transform them.
func filterEven(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 == 0 {
				out <- n
			}
		}
	}()
	return out
}

// DemoPipeline wires the stages together and prints the final output.
// The stages run concurrently: while 'square' processes item N,
// 'generate' is already producing item N+1.
func DemoPipeline() {
	fmt.Println("=== Pipeline ===")

	// Stage 1: produce integers 1..8
	nums := generate(1, 2, 3, 4, 5, 6, 7, 8)

	// Stage 2: square each integer
	squared := square(nums)

	// Stage 3: keep only even squares (4, 16, 36, 64)
	evens := filterEven(squared)

	// Sink: consume the final channel — ranging stops when evens is closed
	for v := range evens {
		fmt.Println(v)
	}
}
