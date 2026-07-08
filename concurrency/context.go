package main

import (
	"context"
	"fmt"
	"time"
)

// --- Context Propagation Pattern ---
//
// context.Context carries deadlines, cancellation signals, and
// request-scoped values across goroutines and API boundaries.
// It is the standard mechanism for controlling goroutine lifetime.
//
// Rule of thumb: every function that does I/O or spawns goroutines
// should accept ctx as its first parameter and check ctx.Done().

// taskWithContext simulates a unit of work that obeys cancellation.
// It uses a select loop: whichever case is ready first wins.
func taskWithContext(ctx context.Context, name string, duration time.Duration) error {
	select {
	case <-time.After(duration):
		// Happy path: the work finished before the context was cancelled.
		fmt.Printf("[%s] completed after %v\n", name, duration)
		return nil
	case <-ctx.Done():
		// The context was cancelled (timeout, deadline, or explicit cancel).
		// ctx.Err() returns context.Canceled or context.DeadlineExceeded.
		return fmt.Errorf("[%s] cancelled: %w", name, ctx.Err())
	}
}

// DemoContextTimeout shows context.WithTimeout.
// The context is cancelled automatically after the deadline,
// even if we never call cancel() ourselves — but we still must
// defer cancel() to release internal timer resources promptly.
func DemoContextTimeout() {
	fmt.Println("=== Context: Timeout ===")

	// Give the whole operation 100 ms.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel() // always defer cancel to free resources, even on success

	// Task A finishes in 50 ms — within the deadline.
	if err := taskWithContext(ctx, "task-A", 50*time.Millisecond); err != nil {
		fmt.Println(err)
	}

	// Task B needs 200 ms — it will be cancelled by the 100 ms deadline.
	if err := taskWithContext(ctx, "task-B", 200*time.Millisecond); err != nil {
		fmt.Println(err)
	}
}

// DemoContextCancel shows context.WithCancel.
// We launch a long-running goroutine and cancel it manually
// when an external event (the done channel here) fires.
func DemoContextCancel() {
	fmt.Println("=== Context: Manual Cancel ===")

	ctx, cancel := context.WithCancel(context.Background())

	// Simulate an external trigger that fires after 80 ms.
	done := make(chan struct{})
	go func() {
		time.Sleep(80 * time.Millisecond)
		close(done)
	}()

	// Monitor the external trigger; call cancel when it fires.
	// This is the typical pattern for wiring OS signals or
	// other shutdown mechanisms to a context tree.
	go func() {
		<-done
		fmt.Println("external trigger fired — cancelling context")
		cancel()
	}()

	// This task would run for 500 ms, but the cancel fires at ~80 ms.
	if err := taskWithContext(ctx, "long-task", 500*time.Millisecond); err != nil {
		fmt.Println(err)
	}
}

// DemoContextValues shows context.WithValue.
// Values should carry request-scoped metadata (trace IDs, auth tokens)
// not function parameters — avoid using WithValue as a side channel
// for optional function arguments.
func DemoContextValues() {
	fmt.Println("=== Context: Values ===")

	// Use an unexported struct key to avoid collisions with
	// other packages that might use the same context.
	type requestIDKey struct{}

	ctx := context.WithValue(context.Background(), requestIDKey{}, "req-abc-123")

	// Simulate middleware reading a trace/request ID out of the context.
	handleRequest := func(ctx context.Context) {
		id, ok := ctx.Value(requestIDKey{}).(string)
		if !ok {
			fmt.Println("no request ID in context")
			return
		}
		fmt.Printf("handling request with ID: %s\n", id)
	}

	handleRequest(ctx)
}
