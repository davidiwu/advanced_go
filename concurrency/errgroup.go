package main

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

// --- errgroup Pattern ---
//
// errgroup.Group coordinates a set of goroutines and collects the
// first non-nil error. It is a higher-level alternative to a manual
// sync.WaitGroup + error channel.
//
// errgroup.WithContext returns a derived context that is cancelled
// automatically when the first goroutine returns an error, so all
// other goroutines can exit early.

// fetchURL simulates an HTTP fetch. If the URL contains "bad",
// it returns an error so we can observe error propagation.
func fetchURL(ctx context.Context, url string) error {
	// In real code: use http.NewRequestWithContext(ctx, ...)
	// and respect ctx.Done() in long operations.
	if errors.Is(ctx.Err(), context.Canceled) {
		// Another goroutine already failed; the context was cancelled.
		return fmt.Errorf("fetch %s: context cancelled before start", url)
	}
	if url == "https://bad.example.com" {
		return fmt.Errorf("fetch %s: connection refused", url)
	}
	fmt.Printf("fetched %s\n", url)
	return nil
}

// DemoErrGroup launches one goroutine per URL.
// errgroup.WithContext wires the group to a context so that when
// "bad.example.com" fails, the context is cancelled and all
// subsequent goroutines can detect it via ctx.Err().
func DemoErrGroup() {
	fmt.Println("=== errgroup ===")

	urls := []string{
		"https://example.com",
		"https://go.dev",
		"https://bad.example.com", // this one will fail
		"https://pkg.go.dev",
	}

	// WithContext: cancels ctx when the first goroutine returns a non-nil error.
	g, ctx := errgroup.WithContext(context.Background())

	for _, url := range urls {
		url := url // capture loop variable — required before Go 1.22
		g.Go(func() error {
			return fetchURL(ctx, url)
		})
	}

	// Wait blocks until all goroutines finish.
	// It returns the first non-nil error (not all errors).
	if err := g.Wait(); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Println("all fetches succeeded")
}

// DemoErrGroupNoError shows the success path: all goroutines succeed
// and g.Wait() returns nil.
func DemoErrGroupNoError() {
	fmt.Println("=== errgroup (no error) ===")

	urls := []string{
		"https://example.com",
		"https://go.dev",
		"https://pkg.go.dev",
	}

	g, ctx := errgroup.WithContext(context.Background())
	for _, url := range urls {
		url := url
		g.Go(func() error {
			return fetchURL(ctx, url)
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Println("all fetches succeeded")
}
