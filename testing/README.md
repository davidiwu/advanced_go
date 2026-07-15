# testing — Benchmarks and Test Infrastructure

This module demonstrates idiomatic Go test patterns: table-driven tests,
parallel subtests, shared helpers via `testing.TB`, benchmarks with allocation
reporting, and interface-based test doubles — all without third-party frameworks.

> **Note:** Unlike other modules in this repo, there is no `main.go`.
> The demo *is* `go test`. See the commands below.

## Key Types in the `testing` Package

Go's `testing` package uses short, single-letter type names. Here's what each one does:

| Type | Full meaning | Passed to | Purpose |
|------|-------------|-----------|---------|
| `*testing.T` | **T**est | `func TestXxx(t *testing.T)` | Unit and integration tests — reports failures, runs subtests, controls parallelism |
| `*testing.B` | **B**enchmark | `func BenchmarkXxx(b *testing.B)` | Performance benchmarks — exposes `b.N` (iteration count), timer controls, and allocation reporting |
| `testing.TB` | **T**est or **B**enchmark | Helper parameters | Interface satisfied by both `*T` and `*B`; write shared helpers against `TB` so they work in both contexts |
| `*testing.F` | **F**uzz | `func FuzzXxx(f *testing.F)` | Fuzz tests — manages the seed corpus (`f.Add`) and delegates to a fuzz target via `f.Fuzz` |
| `*testing.M` | **M**ain | `func TestMain(m *testing.M)` | Test binary lifecycle — runs global setup/teardown around the entire test suite; call `os.Exit(m.Run())` |

### When each type is the receiver

```go
// T — every ordinary test function
func TestAdd(t *testing.T) { ... }

// B — every benchmark function
func BenchmarkAdd(b *testing.B) {
    for range b.N { Add(1, 2) }
}

// TB — shared helpers (accepts either T or B)
func setup(tb testing.TB) {
    tb.Helper()
    ...
}

// F — fuzz targets (go test -fuzz=FuzzAdd)
func FuzzAdd(f *testing.F) {
    f.Add(1, 2)                          // seed corpus
    f.Fuzz(func(t *testing.T, a, b int) { // T inside the fuzz target
        Add(a, b)
    })
}

// M — test binary entry point (optional, one per package)
func TestMain(m *testing.M) {
    setup()
    os.Exit(m.Run()) // runs all Test*, Benchmark*, and Fuzz* in the package
}
```

**Key point:** `T` and `B` share a common ancestor (`TB`) because assertion helpers should not care whether they're called from a test or a benchmark. Always accept `testing.TB` in helper functions — never `*testing.T` alone.

---

## Patterns

### Table-Driven Tests (`wordcount_test.go`)

Collect all cases in a slice of structs, iterate with a single loop, and use
`t.Run` for named subtests.

```go
cases := []struct {
    name string
    in   string
    want int
}{
    {"empty string", "", 0},
    {"single word", "hello", 1},
    {"multiple words", "the quick brown fox", 4},
}

for _, tc := range cases {
    t.Run(tc.name, func(t *testing.T) {
        got := WordCount(tc.in)
        if got != tc.want {
            t.Errorf("WordCount(%q) = %d, want %d", tc.in, got, tc.want)
        }
    })
}
```

**When to use:** Any function with more than two interesting inputs.  
**Key point:** Adding a new case = adding one row. No new test function needed.

---

### `t.Parallel()` in Subtests (`wordcount_test.go`)

Call `t.Parallel()` inside `t.Run` to let independent subtests run concurrently,
reducing total test time for slow operations.

```go
t.Run(tc.name, func(t *testing.T) {
    t.Parallel() // this subtest runs concurrently with its siblings
    got := WordCount(tc.in)
    ...
})
```

**When to use:** Subtests with no shared mutable state and no I/O ordering requirements.  
**Key point:** The parent `TestXxx` function waits for all parallel subtests to finish before returning.

---

### `testing.TB` Shared Helpers (`wordcount_test.go`)

`testing.TB` is the interface satisfied by both `*testing.T` and `*testing.B`.
Write helpers against `TB` so they work in both unit tests and benchmarks.

```go
func requireEqual(tb testing.TB, got, want int) {
    tb.Helper() // failure lines point to the caller, not this function
    if got != want {
        tb.Fatalf("got %d, want %d", got, want)
    }
}
```

**When to use:** Any assertion or setup logic shared between tests and benchmarks.  
**Key point:** `tb.Helper()` is essential — without it, failure output points into the helper, not the test that called it.

---

### Interface-Based Test Doubles (`store_test.go`)

Define the dependency as an interface. Implement a minimal fake in the test file
that satisfies the interface. No mocking framework needed.

```go
// Production interface
type Storer interface {
    Save(id string, val int) error
    Load(id string) (int, error)
}

// In-test fake — lives in store_test.go
type fakeStore struct{ data map[string]int }

func (f *fakeStore) Save(id string, val int) error { f.data[id] = val; return nil }
func (f *fakeStore) Load(id string) (int, error)   { ... }
```

**When to use:** Any time you want to test business logic without a real database, HTTP server, or external service.  
**Key point:** The interface IS the contract. A struct with the right methods IS the mock. Go's implicit satisfaction means no registration or annotation is needed.

---

### Benchmarks & `b.ReportAllocs()` (`bench_test.go`)

Benchmark functions start with `Benchmark` and accept `*testing.B`. The
framework adjusts `b.N` automatically — you never set it yourself.

```go
func BenchmarkWordCount(b *testing.B) {
    b.ReportAllocs() // print allocs/op and B/op without -benchmem flag
    input := "the quick brown fox jumps over the lazy dog"
    for range b.N {
        WordCount(input)
    }
}
```

Sub-benchmarks via `b.Run` compare the same function across different inputs:

```go
func BenchmarkWordCount_Sizes(b *testing.B) {
    for _, tc := range cases {
        b.Run(tc.name, func(b *testing.B) {
            for range b.N { WordCount(tc.input) }
        })
    }
}
```

**When to use:** Any time you need to measure performance or track allocations.  
**Key point:** Zero allocs/op means no heap pressure — use `b.ReportAllocs()` to confirm after optimising a hot path.

---

## Running

```bash
# Run all tests with verbose per-case output
go test -v ./testing/

# Run a specific test
go test -v -run TestWordCount ./testing/

# Run a specific subtest
go test -v -run TestWordCount/empty_string ./testing/

# Run all benchmarks (include allocation stats)
go test -bench=. -benchmem ./testing/

# Run a specific benchmark, 5 times
go test -bench=BenchmarkReverse -benchmem -count=5 ./testing/
```
