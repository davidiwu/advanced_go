package testing_patterns

import "testing"

// --- Benchmarks ---
//
// Benchmark functions start with Benchmark and accept *testing.B.
// The testing framework runs the body b.N times, automatically adjusting N
// until the total run time is stable — you never set N yourself.
//
// Run benchmarks with:
//   go test -bench=. ./testing/
//   go test -bench=. -benchmem ./testing/      # include allocation stats
//   go test -bench=BenchmarkReverse -count=5   # run a specific bench 5 times
//
// b.ReportAllocs() is equivalent to passing -benchmem for this benchmark only —
// it prints allocs/op and B/op without affecting other benchmarks.

func BenchmarkWordCount(b *testing.B) {
	b.ReportAllocs()
	input := "the quick brown fox jumps over the lazy dog"
	for range b.N {
		WordCount(input)
	}
}

func BenchmarkReverse(b *testing.B) {
	b.ReportAllocs()
	input := "the quick brown fox jumps over the lazy dog"
	for range b.N {
		Reverse(input)
	}
}

// BenchmarkReverse_Unicode demonstrates that Reverse handles multi-byte runes.
// Comparing it against the ASCII benchmark shows the cost of rune conversion.
func BenchmarkReverse_Unicode(b *testing.B) {
	b.ReportAllocs()
	input := "日本語のテスト文字列"
	for range b.N {
		Reverse(input)
	}
}

// --- Sub-benchmarks ---
//
// Just like t.Run, b.Run creates named sub-benchmarks.
// Useful for comparing the same function across different input sizes.

func BenchmarkWordCount_Sizes(b *testing.B) {
	cases := []struct {
		name  string
		input string
	}{
		{"small", "one two three"},
		{"medium", "the quick brown fox jumps over the lazy dog"},
		{"large", "word one two three four five six seven eight nine ten eleven twelve thirteen fourteen fifteen sixteen seventeen eighteen nineteen twenty"},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				WordCount(tc.input)
			}
		})
	}
}
