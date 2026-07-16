# pointers

Safe pointer usage patterns in Go: what makes Go pointers memory-safe, and the
runtime gotchas that still bite.

## What this covers

| File | Pattern |
|---|---|
| `nilsafety.go` | Nil checks, nil receiver methods, nil as an optional sentinel |
| `receivers.go` | Value vs pointer receivers, mutation semantics, interface satisfaction |
| `aliasing.go` | Unintended slice/map aliasing, defensive copy pattern |
| `escape.go` | Escape analysis, stack vs heap, interface boxing cost |
| `loopcapture.go` | Loop variable capture with goroutines and the two fixes |

## Why Go pointers are safe (by design)

| Guarantee | How Go provides it |
|---|---|
| No dangling pointers | Garbage collector keeps target alive while any pointer to it exists |
| No use-after-free | No manual `free()` — memory is reclaimed only when unreachable |
| No pointer arithmetic | `ptr + 1` is a compile error; `unsafe.Pointer` is required and visible |
| No raw memory casting | Cannot cast `*int` to `*float64` outside `unsafe` |
| Nil is explicit | Zero value of a pointer is `nil`, not a random address; dereference panics loudly |
| Stack safety | Escape analysis automatically heap-allocates locals whose address outlives the frame |

## Key ideas

### Nil safety

The zero value of every pointer type is `nil`. Dereferencing `nil` panics at
runtime — the compiler does not catch it. Guard all pointers from outside your
package.

```go
func effectiveTimeout(c Config) int {
    if c.Timeout == nil {
        return 30 // default
    }
    return *c.Timeout
}
```

A method with a pointer receiver can be called on `nil` safely if it guards first.
This is idiomatic for optional/chainable APIs (e.g. a nil linked-list node):

```go
func (n *Node) Len() int {
    if n == nil {
        return 0
    }
    return 1 + n.Next.Len()
}
```

### Value vs pointer receivers

```go
func (c counter) incrementValue()   { c.count++ }  // mutates copy — no effect
func (c *counter) incrementPointer(){ c.count++ }  // mutates original
```

Interface satisfaction follows the method set:

```go
// String has a pointer receiver → only *point satisfies Stringer, not point
func (p *point) String() string { ... }

printStringer(&point{3, 4})  // ok
// printStringer(point{3, 4}) // compile error
```

Go auto-dereferences for method calls on addressable values:
`c.incrementPointer()` is sugar for `(&c).incrementPointer()`.

### Aliasing

Slices share a backing array. Storing a caller's slice directly means the caller
can mutate your internal state:

```go
// Bad: caller's src and s.items share the same backing array
func newStoreBad(items []string) *store {
    return &store{items: items}
}

// Good: copy on the way in
func newStoreGood(items []string) *store {
    cp := make([]string, len(items))
    copy(cp, items)
    return &store{items: cp}
}
```

Maps are always reference types — assigning a map variable copies only the pointer:

```go
alias := original  // NOT a copy; both variables see the same map
```

### Escape analysis

The compiler decides stack vs heap automatically. Stack allocation is essentially
free; heap allocation has GC cost.

```go
// Stack: value returned by copy — no allocation
func newPairByValue(x, y int) smallPair { return smallPair{x, y} }

// Heap: address returned — pair must outlive the stack frame
func newPairByPointer(x, y int) *smallPair {
    p := smallPair{x, y}  // escapes to heap
    return &p
}
```

Interface boxing always escapes — storing a concrete value in `interface{}` causes
a heap allocation. Inspect escape decisions with:

```bash
go build -gcflags="-m" ./pointers/
```

Rule of thumb: prefer returning small structs by value; pass large structs (> a
few words) by pointer to avoid copy cost.

### Loop variable capture

Before Go 1.22, all iterations of a `for` range shared one variable. Goroutines
capturing it by reference would all see the final value:

```go
// Bug (pre-1.22): every goroutine prints the last url
for _, url := range urls {
    go func() { fetch(url) }()
}
```

Two fixes that work in all Go versions:

```go
// Fix 1 — shadow the variable
for _, url := range urls {
    url := url  // new variable per iteration
    go func() { fetch(url) }()
}

// Fix 2 — pass as argument
for _, url := range urls {
    go func(u string) { fetch(u) }(url)
}
```

Go 1.22+ gives each iteration its own variable automatically — no extra code needed.
