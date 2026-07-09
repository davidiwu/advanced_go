package main

import (
	"errors"
	"fmt"
)

// --- Sentinel Errors Pattern ---
//
// A sentinel error is a package-level variable declared with errors.New.
// Callers compare against it using errors.Is, which traverses the full
// wrapping chain — so a sentinel stays detectable even after being wrapped
// with fmt.Errorf("%w", ...).
//
// Rule: declare sentinels at the package level; never compare with ==
// on error values that may have been wrapped.

var (
	ErrNotFound   = errors.New("not found")
	ErrPermission = errors.New("permission denied")
	ErrTimeout    = errors.New("operation timed out")
)

// findUser simulates a lookup that may fail with a known sentinel.
func findUser(id int) error {
	if id <= 0 {
		return ErrNotFound
	}
	if id == 99 {
		return ErrPermission
	}
	return nil
}

// loadProfile wraps the sentinel with context. errors.Is still matches
// the underlying ErrNotFound even through the wrapper.
func loadProfile(id int) error {
	if err := findUser(id); err != nil {
		return fmt.Errorf("loadProfile(%d): %w", id, err)
	}
	return nil
}

func DemoSentinel() {
	fmt.Println("=== Sentinel Errors ===")

	cases := []struct {
		id   int
		desc string
	}{
		{1, "valid user"},
		{0, "missing user"},
		{99, "restricted user"},
	}

	for _, c := range cases {
		err := loadProfile(c.id)
		switch {
		case err == nil:
			fmt.Printf("  %-16s → ok\n", c.desc)
		case errors.Is(err, ErrNotFound):
			// errors.Is traverses the chain: loadProfile wraps findUser's ErrNotFound
			fmt.Printf("  %-16s → ErrNotFound   (%v)\n", c.desc, err)
		case errors.Is(err, ErrPermission):
			fmt.Printf("  %-16s → ErrPermission (%v)\n", c.desc, err)
		default:
			fmt.Printf("  %-16s → unexpected: %v\n", c.desc, err)
		}
	}

	// Demonstrate that == fails on a wrapped sentinel but errors.Is succeeds.
	wrapped := fmt.Errorf("outer: %w", ErrNotFound)
	fmt.Printf("\nwrapped == ErrNotFound:        %v  (direct == fails on wrapped errors)\n",
		wrapped == ErrNotFound)
	fmt.Printf("errors.Is(wrapped, ErrNotFound): %v  (Is traverses the chain)\n",
		errors.Is(wrapped, ErrNotFound))
}
