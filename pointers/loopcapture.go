package main

import (
	"fmt"
	"sync"
)

// --- Loop Variable Capture ---
//
// Before Go 1.22, all iterations of a for-range loop shared a single loop
// variable. Any closure or goroutine that captured the variable by reference
// (i.e. held a pointer to it) would see the value at the time it ran — almost
// always the final value after the loop completed, not the value at the time the
// goroutine was launched.
//
// Go 1.22 fixed this: each iteration now has its own copy of the loop variable,
// so closures capture the per-iteration value automatically. Code targeting
// Go 1.22+ does not need the manual shadowing workaround.
//
// This file demonstrates both the old problem and the two fixes so the pattern
// is recognisable when reading pre-1.22 code.
//
// Fix 1 — explicit shadow (pre-1.22 workaround):
//   for _, v := range items {
//       v := v  // shadow: new variable per iteration, closes over this copy
//       go func() { use(v) }()
//   }
//
// Fix 2 — pass as argument (works in all Go versions):
//   for _, v := range items {
//       go func(v string) { use(v) }(v)
//   }
//
// Fix 3 — Go 1.22+: no extra code needed; each iteration already has its own v.

// captureByArgument launches goroutines that each receive their iteration value
// as a function argument — safe in all Go versions.
func captureByArgument(items []string) []string {
	var mu sync.Mutex
	var results []string
	var wg sync.WaitGroup

	for _, v := range items {
		wg.Add(1)
		go func(val string) { // val is a fresh copy per goroutine
			defer wg.Done()
			mu.Lock()
			results = append(results, val)
			mu.Unlock()
		}(v) // pass current v as argument
	}

	wg.Wait()
	return results
}

// captureGo122 relies on Go 1.22 per-iteration variables. No extra shadowing needed.
func captureGo122(items []string) []string {
	var mu sync.Mutex
	var results []string
	var wg sync.WaitGroup

	for _, v := range items {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			results = append(results, v) // safe: each iteration has its own v (Go 1.22+)
			mu.Unlock()
		}()
	}

	wg.Wait()
	return results
}

func DemoLoopCapture() {
	fmt.Println("=== Loop Variable Capture ===")

	items := []string{"a", "b", "c", "d"}

	r1 := captureByArgument(items)
	fmt.Printf("captureByArgument: got %d items (want 4), all unique: %v\n",
		len(r1), allUnique(r1))

	r2 := captureGo122(items)
	fmt.Printf("captureGo122:       got %d items (want 4), all unique: %v\n",
		len(r2), allUnique(r2))
}

func allUnique(ss []string) bool {
	seen := make(map[string]bool, len(ss))
	for _, s := range ss {
		if seen[s] {
			return false
		}
		seen[s] = true
	}
	return true
}
