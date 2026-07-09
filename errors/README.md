# Error Handling Patterns

Go treats errors as values. This module covers five patterns that appear in production Go code: signalling known conditions, adding context through wrapping, building inspectable error types, recovering from panics, and aggregating multiple failures.

## Patterns

### 1. Sentinel Errors (`sentinel.go`)

A sentinel is a package-level `var` declared with `errors.New`. Callers compare against it with `errors.Is`, which traverses wrapping chains — so a sentinel remains detectable even after being wrapped with `fmt.Errorf("%w", ...)`.

```go
var (
    ErrNotFound   = errors.New("not found")
    ErrPermission = errors.New("permission denied")
)

// Wrapping preserves the sentinel.
func loadProfile(id int) error {
    if err := findUser(id); err != nil {
        return fmt.Errorf("loadProfile(%d): %w", id, err)
    }
    return nil
}

// errors.Is walks the chain — == would miss wrapped sentinels.
if errors.Is(err, ErrNotFound) { ... }
```

**When to use:** When callers need to branch on a specific, named condition. Keep sentinel vars in the package that owns the concept; never compare raw error values with `==` when wrapping is possible.

### 2. Wrapping and Unwrapping Chains (`wrapping.go`)

`fmt.Errorf("context: %w", err)` attaches a message to an error while keeping the original reachable. Each layer of the call stack adds its own context; `errors.Unwrap` lets you step through the chain manually.

```go
func writeChunk(path string, size int) error {
    if size > 100 { return ErrDiskFull }
    return nil
}

func saveFile(path string, data []byte) error {
    if err := writeChunk(path, len(data)); err != nil {
        return fmt.Errorf("saveFile %q: %w", path, err)
    }
    return nil
}

// errors.Is finds ErrDiskFull anywhere in the chain.
errors.Is(err, ErrDiskFull) // true

// Walk the chain manually.
for e := err; e != nil; e = errors.Unwrap(e) {
    fmt.Println(e)
}
```

**When to use:** Always wrap when crossing a meaningful boundary (function, package, subsystem) where the additional context would help a reader of a log line understand what went wrong. Don't wrap if you have nothing useful to add — just return the error unchanged.

### 3. Custom Error Types with `Is` / `As` / `Unwrap` (`custom.go`)

Struct error types carry structured fields (status codes, request IDs, field names) that callers extract with `errors.As`. Three optional methods extend the default behaviour:

| Method | Purpose |
|---|---|
| `Unwrap() error` | makes the type transparent to `errors.Is`/`errors.As` chains |
| `Is(target error) bool` | custom equality — match by code or category, not by pointer identity |
| `As(target any) bool` | custom extraction (rarely needed; `errors.As` handles most cases) |

```go
type APIError struct {
    Code    int
    Message string
    Cause   error
}

func (e *APIError) Error() string  { return fmt.Sprintf("API %d %s", e.Code, e.Message) }
func (e *APIError) Unwrap() error  { return e.Cause }

// Custom Is: match any *APIError with the same code.
func (e *APIError) Is(target error) bool {
    var t *APIError
    return errors.As(target, &t) && e.Code == t.Code
}

var ErrAPI404 = &APIError{Code: 404}

// Caller checks by code, not by message.
if errors.Is(err, ErrAPI404) { ... }

// Or extract the full struct.
var apiErr *APIError
if errors.As(err, &apiErr) {
    log.Printf("code=%d msg=%s", apiErr.Code, apiErr.Message)
}
```

**When to use:** When callers need to inspect error fields programmatically, not just compare against a sentinel. The `Is` override is useful when equality should be semantic (same code) rather than identity (same pointer).

### 4. panic / recover (`panic_recover.go`)

`panic` is for unrecoverable states and programmer errors — violated invariants, nil dereferences, out-of-bounds access. It is **not** a substitute for error returns.

`recover()` stops a panic and returns its value, but only inside a deferred function in the same goroutine. The canonical library pattern wraps user-supplied callbacks to ensure panics never propagate to callers as crashes:

```go
func safeCall(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            switch v := r.(type) {
            case error:
                err = v
            default:
                err = fmt.Errorf("panic: %v", v)
            }
        }
    }()
    fn()
    return nil
}

err := safeCall(func() {
    _ = divide(10, 0) // panics
})
// err is now a non-nil error instead of a crash
```

**When to use:** Use `panic` for logic errors that indicate a bug (never for expected runtime conditions). Use `recover` at package boundaries where callers must not be exposed to internal panics — HTTP handler middleware and test frameworks are the most common examples.

### 5. Multi-Error Aggregation (`join.go`)

Aggregate multiple errors into one by implementing `Unwrap() []error` — the multi-unwrap protocol that `errors.Is` and `errors.As` traverse. Go 1.20 ships `errors.Join` built in; on earlier toolchains the same contract is a small custom type.

```go
type multiError struct{ errs []error }

func (m *multiError) Error() string   { /* join messages */ }
func (m *multiError) Unwrap() []error { return m.errs }

// Is provides Go 1.19-compatible traversal (1.20+ handles Unwrap[]error natively).
func (m *multiError) Is(target error) bool {
    for _, e := range m.errs {
        if errors.Is(e, target) { return true }
    }
    return false
}

func validateUser(name, email, password string) error {
    var errs []error
    if name == ""        { errs = append(errs, ErrEmptyName) }
    if !validEmail(email){ errs = append(errs, ErrBadEmail)  }
    if len(password) < 8 { errs = append(errs, ErrShortPass) }
    return joinErrs(errs...) // nil when errs is empty
}

err := validateUser("", "bad", "x")
errors.Is(err, ErrEmptyName)  // true — searches all members
errors.Is(err, ErrBadEmail)   // true
```

In Go 1.20+ `errors.Join(errs...)` replaces the custom `joinErrs` wrapper.

**When to use:** Batch validation, parallel operations, or any scenario where you want to report all failures rather than stopping at the first.

## Running the examples

```bash
cd errors && go run .
```

## Key Takeaways

| Pattern | Rule |
|---|---|
| Sentinel errors | Declare as `var Err... = errors.New(...)`, compare with `errors.Is` |
| Wrapping | Use `%w` to add context; return unwrapped if you have nothing to add |
| Custom types | Implement `Unwrap` for chain transparency; `Is` for semantic equality |
| panic / recover | `panic` for bugs only; `recover` at library boundaries |
| Multi-error | Implement `Unwrap() []error`; use `errors.Join` on Go 1.20+ |
