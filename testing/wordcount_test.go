package testing_patterns

import "testing"

// --- Table-Driven Tests ---
//
// A table-driven test collects all cases in a slice of structs, then
// iterates over them with a single loop. This keeps test logic in one
// place and makes adding new cases trivial — just add a row to the table.
//
// t.Run creates a named subtest for each case. This gives:
//   - per-case PASS/FAIL lines in -v output
//   - the ability to run a single case with -run TestWordCount/empty_string

func TestWordCount(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want int
	}{
		{"empty string", "", 0},
		{"whitespace only", "   ", 0},
		{"single word", "hello", 1},
		{"multiple words", "the quick brown fox", 4},
		{"extra spaces", "  hello   world  ", 2},
		{"tabs and newlines", "a\tb\nc", 3},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel lets subtests run concurrently with each other.
			// Use it when subtests are independent and have no shared mutable state.
			t.Parallel()

			got := WordCount(tc.in)
			if got != tc.want {
				t.Errorf("WordCount(%q) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"single char", "a", "a"},
		{"ascii", "hello", "olleh"},
		{"unicode", "日本語", "語本日"},
		{"palindrome", "racecar", "racecar"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := Reverse(tc.in)
			if got != tc.want {
				t.Errorf("Reverse(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// --- testing.TB helper ---
//
// testing.TB is the interface satisfied by both *testing.T and *testing.B.
// Writing helpers that accept testing.TB means the same helper works in
// both unit tests and benchmarks — no duplication.

func requireEqual(tb testing.TB, got, want int) {
	tb.Helper() // marks this as a helper so failure lines point to the caller
	if got != want {
		tb.Fatalf("got %d, want %d", got, want)
	}
}

func TestWordCountWithHelper(t *testing.T) {
	requireEqual(t, WordCount("one two three"), 3)
	requireEqual(t, WordCount(""), 0)
}
