# Interface Patterns

Go interfaces are satisfied implicitly, composed freely, and intentionally kept small. This module covers the six core patterns that appear most often in real Go code.

## Patterns

### 1. Implicit Interface Satisfaction (`implicit.go`)

A type satisfies an interface the moment it has all the required methods — no `implements` keyword needed. This decouples producers from consumers and lets you retrofit an existing type onto a new interface without touching its source.

```go
type Animal interface {
    Sound() string
    Name() string
}

// Dog satisfies Animal with no declaration.
type Dog struct{ name string }
func (d Dog) Sound() string { return "Woof" }
func (d Dog) Name() string  { return d.name }
```

**When to use:** Always — it's how all Go interfaces work. The key insight is that the interface and the type can live in completely different packages.

### 2. Interface Composition (`composition.go`)

Embed smaller interfaces into larger ones to build composable contracts. Functions declare exactly the surface they need — a function that only reads should accept `io.Reader`, not `io.ReadWriter`.

```go
type BufferedLogger interface {
    Logger   // Log(msg string)
    Flusher  // Flush() error
}
```

**When to use:** When you have two or more behaviors that are useful independently but sometimes needed together. The standard library's `io.ReadWriter`, `io.ReadWriteCloser`, etc. are the canonical example.

### 3. Type Assertions and Type Switches (`typeassert.go`)

Extract the concrete value from an interface variable. Use the two-value form to avoid panics, and a type switch when branching on multiple possible types.

```go
// Two-value assertion — safe, never panics.
if c, ok := s.(Circle); ok {
    fmt.Println("radius:", c.Radius)
}

// Type switch — cleaner than a chain of assertions.
switch v := s.(type) {
case Circle:    ...
case Rectangle: ...
default:        ...
}
```

**When to use:** When you genuinely need type-specific behavior not expressed by the interface. Use sparingly — heavy use of type switches is a signal the interface may be too narrow.

### 4. Interface Segregation (`segregation.go`)

Prefer many small, focused interfaces over one large "god interface." Each caller declares only the methods it uses, making implementations easier to swap and test.

```go
type Reader  interface { Read(key string) (string, bool) }
type Writer  interface { Write(key, value string) }
type Deleter interface { Delete(key string) }

// cacheLayer only needs to read:
func cacheLayer(r Reader, keys []string) { ... }

// populate only needs to write:
func populate(w Writer, entries map[string]string) { ... }
```

**When to use:** Whenever you're designing a function's parameter types. Ask: "which methods does this function actually call?" Accept only those.

### 5. Stringer and Error Interfaces (`stringer.go`)

`fmt.Stringer` (`String() string`) and `error` (`Error() string`) are Go's two most important single-method interfaces. Implementing them lets your types integrate naturally with `fmt`, `log`, `errors.As`/`Is`, and the rest of the standard library.

```go
// fmt.Stringer: controls how the type prints.
func (d Direction) String() string { ... }

// Custom error type: carries structured context.
type ValidationError struct{ Field, Message string }
func (e *ValidationError) Error() string { ... }

// Caller extracts structured fields without parsing strings.
var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Println(ve.Field, ve.Message)
}
```

**When to use:** Always implement `Stringer` on types that have a meaningful text form. Use custom error types when callers need to inspect the error's content programmatically.

### 6. Accept Interfaces, Return Structs (`acceptreturn.go`)

The canonical Go API design rule: **function parameters should be interface types; return types should be concrete structs.**

- Accepting an interface makes the function easy to test (inject a fake) and easy to extend (swap implementations).
- Returning a concrete struct gives callers access to all fields and future methods without type assertions.

```go
// Parameter: interface — any Sorter implementation works.
// Return: concrete struct — caller accesses all fields directly.
func ProcessData(s Sorter, data []int) Result { ... }

r := ProcessData(StdSorter{}, data)
fmt.Println(r.Max) // no type assertion needed
```

**When to use:** As a default when designing any function or method. Exceptions exist (factory functions sometimes return an interface to hide implementation details), but this rule covers the vast majority of cases.

## Running the examples

```bash
cd interfaces && go run .
```

## Key Takeaways

| Principle | Rule |
|---|---|
| Implicit satisfaction | No `implements` declaration — types satisfy interfaces by having the methods |
| Composition | Embed small interfaces into larger ones; keep each one focused |
| Type assertions | Use the `v, ok` form; use type switches over chains of assertions |
| Segregation | Accept the narrowest interface that covers what you need |
| Stringer / error | Implement these to integrate with the standard library |
| Accept/return | Accept interfaces, return structs |
