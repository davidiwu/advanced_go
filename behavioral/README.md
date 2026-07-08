# Behavioral Patterns in Go

Behavioral patterns define how types and functions are structured and composed. Unlike concurrency patterns (which are about coordination), these patterns shape the API and design of your Go code.

---

## Table of Contents

1. [Functional Options](#1-functional-options)
2. [Option Type](#2-option-type)
3. [Table-Driven Tests](#3-table-driven-tests)
4. [Middleware Chaining](#4-middleware-chaining)

---

## 1. Functional Options

### What it is
A pattern for configuring a struct using a slice of `func(*T)` arguments. Instead of a growing parameter list or a config struct, callers pass named option functions.

### When to use
- A constructor has more than 2-3 optional parameters
- You want to add new options in the future without breaking callers
- You want zero values to be meaningful defaults (callers only pass what they want to change)
- Public library APIs where backwards compatibility matters

### How it works

```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
}

type Option func(*Server)

func WithPort(port int) Option {
    return func(s *Server) { s.port = port }
}

func NewServer(opts ...Option) *Server {
    s := &Server{host: "localhost", port: 8080} // sensible defaults
    for _, opt := range opts {
        opt(s) // each option mutates the struct
    }
    return s
}

// Caller:
srv := NewServer(WithPort(9090), WithTimeout(30*time.Second))
```

### Key points
- Default values live in the constructor, not scattered across callers.
- Adding a new option is backwards-compatible — existing callers don't change.
- Options are composable: you can define preset bundles as `[]Option` slices.
- Validate after applying all options so you get one coherent error, not partial state.

---

## 2. Option Type (Maybe / Optional)

### What it is
A generic wrapper that explicitly represents a value that may or may not be present, as an alternative to pointer semantics (`*T`) or sentinel zero values.

### When to use
- A function return value is legitimately absent (not an error, just missing)
- A struct field is truly optional and zero value is a valid real value (e.g., `0` age is valid)
- You want to force callers to handle the "absent" case explicitly, rather than silently using a zero value
- As a cleaner alternative to `(T, bool)` return pairs

### How it works

```go
type Option[T any] struct {
    value T
    valid bool
}

func Some[T any](v T) Option[T] { return Option[T]{value: v, valid: true} }
func None[T any]() Option[T]    { return Option[T]{} }

func (o Option[T]) Unwrap() (T, bool) { return o.value, o.valid }

// Caller:
age := findAge("Alice")        // returns Option[int]
if v, ok := age.Unwrap(); ok {
    fmt.Println(v)
}
```

### Key points
- Clearer intent than `*T`: a pointer often means "heap-allocated" as much as "optional".
- Zero value of `Option[T]` is naturally `None` — safe to use without initialization.
- Not idiomatic in all Go codebases; use where the absence/presence distinction is the primary concern.

---

## 3. Table-Driven Tests

### What it is
A testing style where test cases are defined as a slice of structs, and a single loop runs each case as a subtest via `t.Run`. It is the dominant Go testing idiom.

### When to use
- Any function with multiple input/output combinations to verify
- When adding a new test case should require only adding a row to a table, not writing a new function
- When you want named, isolated subtests (failures are easy to pinpoint)

### How it works

```go
func TestAdd(t *testing.T) {
    cases := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 1, 2, 3},
        {"zero", 0, 0, 0},
        {"negative", -1, -2, -3},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got := Add(tc.a, tc.b)
            if got != tc.expected {
                t.Errorf("got %d, want %d", got, tc.expected)
            }
        })
    }
}
```

### Key points
- `t.Run` isolates each case: a failure in one doesn't stop others from running.
- Run a single case: `go test -run TestAdd/zero`.
- Use `t.Parallel()` inside `t.Run` for independent cases to speed up the suite.
- The `name` field makes failures self-documenting in test output.

---

## 4. Middleware Chaining

### What it is
A pattern where behavior is layered around a core handler using `func(Handler) Handler` wrappers. Each middleware wraps the next, forming a chain. Common in HTTP servers but applicable to any handler-style interface.

### When to use
- Adding cross-cutting concerns (logging, auth, rate limiting, tracing) without touching core logic
- Building composable, reusable processing layers
- HTTP servers, CLI command handlers, gRPC interceptors, event processors

### How it works

```go
type Handler func(req Request) Response

type Middleware func(Handler) Handler

// A middleware wraps a handler: it can run logic before, after, or around it.
func Logger(next Handler) Handler {
    return func(req Request) Response {
        fmt.Printf("→ %s\n", req.Path)
        resp := next(req)           // call the inner handler
        fmt.Printf("← %d\n", resp.Status)
        return resp
    }
}

// Chain applies middlewares right-to-left so the first in the list
// is the outermost (first to run on the way in).
func Chain(h Handler, middlewares ...Middleware) Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}

// Caller:
handler := Chain(coreHandler, Logger, Auth, RateLimiter)
response := handler(request)
```

### Key points
- Middlewares are applied in reverse order so the first listed runs first.
- Each middleware is independent and reusable — no coupling to other middlewares.
- The pattern works for any `func(In) Out` shape, not just `http.Handler`.
- For HTTP, `net/http` uses `http.Handler` interface; the same chaining principle applies.

---

## Combining Patterns

| Combination | Use case |
|---|---|
| Functional Options + Middleware | Configure which middlewares to apply at construction time |
| Option Type + Table-Driven Tests | Test functions that return optional values clearly |
| Middleware + Context | Pass request-scoped values (trace ID, user) through the chain |

---

## Further Reading

- [Go Blog: Functional Options (Dave Cheney)](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)
- [Go Blog: Table Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [justforfunc: Middleware in Go](https://www.youtube.com/watch?v=xyDkyFjzFVc)
