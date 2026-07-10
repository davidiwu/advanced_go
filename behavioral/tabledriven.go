package main

import "fmt"

// --- Table-Driven Tests Pattern ---
//
// Test cases are defined as a slice of structs (the "table").
// A single loop runs each case, calling t.Run to create a named subtest.
// This is the dominant Go testing idiom because:
//   - Adding a new case = adding one row, not a new function
//   - Each case is isolated: a failure in one doesn't stop the others
//   - Failures are self-documenting: the test name includes the case name
//
// This file demonstrates the pattern with runnable (non-test) code so you
// can see the mechanics without needing `go test`. The real pattern lives
// in *_test.go files — see the bottom of this file for the test version.

// --- Functions under test ---

// Add adds two integers. Trivial, but useful for a clean demo.
func Add(a, b int) int { return a + b }

// Clamp returns v clamped to [lo, hi].
// Demonstrates a function with multiple interesting boundary cases.
func Clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// --- Simulated table-driven runner ---
// In real code this would be inside TestAdd(t *testing.T) using t.Run.
// Here we replicate the same structure so the pattern is visible without
// needing the testing package.

type testCase[In, Out any] struct {
	name     string
	input    In
	expected Out
	// 'got' is filled in by the runner; we compare it to expected.
}

// runAddTable demonstrates the table-driven approach for Add.
func runAddTable() {
	fmt.Println("--- Add ---")

	cases := []struct {
		name     string
		a, b     int
		expected int
	}{
		// Each row is one scenario. The name appears in test output on failure.
		{"both positive", 1, 2, 3},
		{"both zero", 0, 0, 0},
		{"both negative", -1, -2, -3},
		{"mixed signs", -5, 10, 5},
		{"large numbers", 1_000_000, 2_000_000, 3_000_000},
	}

	for _, tc := range cases {
		// t.Run(tc.name, func(t *testing.T) { ... }) in real test code.
		got := Add(tc.a, tc.b)
		status := "PASS"
		if got != tc.expected {
			status = fmt.Sprintf("FAIL: got %d, want %d", got, tc.expected)
		}
		fmt.Printf("  %-20s %s\n", tc.name, status)
	}
}

// runClampTable demonstrates table-driven tests for a function with boundaries.
func runClampTable() {
	fmt.Println("--- Clamp ---")

	cases := []struct {
		name     string
		v        int
		lo, hi   int
		expected int
	}{
		// Boundary cases are where bugs hide — the table makes them explicit.
		{"below lo", -5, 0, 10, 0},
		{"at lo", 0, 0, 10, 0},
		{"within range", 5, 0, 10, 5},
		{"at hi", 10, 0, 10, 10},
		{"above hi", 15, 0, 10, 10},
		{"lo equals hi", 7, 5, 5, 5},
	}

	for _, tc := range cases {
		got := Clamp(tc.v, tc.lo, tc.hi)
		status := "PASS"
		if got != tc.expected {
			status = fmt.Sprintf("FAIL: got %d, want %d", got, tc.expected)
		}
		fmt.Printf("  %-20s %s\n", tc.name, status)
	}
}

// DemoTableDriven runs both tables and prints pass/fail for each case.
func DemoTableDriven() {
	fmt.Println("=== Table-Driven Tests ===")
	runAddTable()
	runClampTable()
}

/*
Real test file equivalent (would live in tabledriven_test.go):

func TestAdd(t *testing.T) {
	cases := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"both positive", 1, 2, 3},
		{"both zero", 0, 0, 0},
		{"both negative", -1, -2, -3},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel() here if cases are independent
			got := Add(tc.a, tc.b)
			if got != tc.expected {
				t.Errorf("got %d, want %d", got, tc.expected)
			}
		})
	}
}
*/
