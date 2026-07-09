package main

import (
	"fmt"
	"strings"
)

// --- strings.Builder Pattern ---
//
// String concatenation with + creates a new allocation per operation
// because strings are immutable in Go. strings.Builder accumulates bytes
// in a mutable buffer and only allocates the final string once.
//
// Rule of thumb:
//   N = 2..3 pieces → + is fine
//   N > 3 pieces, or N is unknown → use strings.Builder

// joinWithPlus concatenates items using + inside a loop.
// Each iteration allocates a new string (O(N²) total byte copies).
func joinWithPlus(items []string, sep string) string {
	result := ""
	for i, item := range items {
		if i > 0 {
			result += sep
		}
		result += item
	}
	return result
}

// joinWithBuilder uses strings.Builder — one allocation for the final string.
func joinWithBuilder(items []string, sep string) string {
	var b strings.Builder
	// Preallocate to avoid internal buffer growth.
	total := 0
	for _, s := range items {
		total += len(s)
	}
	total += len(sep) * (len(items) - 1)
	b.Grow(total)

	for i, item := range items {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(item)
	}
	return b.String()
}

// buildTable demonstrates fmt.Fprintf into a Builder to avoid fmt.Sprintf
// allocations when composing a multi-line report.
func buildTable(rows [][]string) string {
	var b strings.Builder
	b.WriteString("name            value\n")
	b.WriteString("--------------- -----\n")
	for _, row := range rows {
		fmt.Fprintf(&b, "%-15s %s\n", row[0], row[1])
	}
	return b.String()
}

func DemoStringBuilder() {
	fmt.Println("=== strings.Builder ===")

	words := []string{"the", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"}

	plus := joinWithPlus(words, " ")
	builder := joinWithBuilder(words, " ")
	fmt.Printf("joinWithPlus:    %q\n", plus)
	fmt.Printf("joinWithBuilder: %q\n", builder)
	fmt.Printf("results match: %v\n", plus == builder)

	rows := [][]string{
		{"latency_p99", "12ms"},
		{"throughput", "50k/s"},
		{"error_rate", "0.01%"},
	}
	fmt.Println()
	fmt.Print(buildTable(rows))

	// Show that strings.Join is the stdlib convenience wrapper for the same pattern.
	stdlib := strings.Join(words, " ")
	fmt.Printf("\nstrings.Join matches: %v\n", stdlib == builder)
}
