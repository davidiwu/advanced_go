# Performance Patterns

Go gives you low-level control over allocation, memory layout, and data structures. This module covers five patterns that appear repeatedly in performance-sensitive Go code.

## Patterns

### 1. sync.Pool (`pool.go`)

`sync.Pool` holds a set of reusable temporary objects. When a goroutine needs a short-lived object (buffer, parser context, scratch space), it borrows from the pool instead of allocating fresh. After use the object is returned, ready for the next caller.

```go
var bufPool = sync.Pool{
    New: func() any {
        b := make([]byte, 0, 64)
        return &b
    },
}

func process(data string) string {
    bp := bufPool.Get().(*[]byte)
    buf := (*bp)[:0]           // reset length, keep capacity

    buf = append(buf, data...)

    result := string(buf)
    *bp = buf
    bufPool.Put(bp)            // return for reuse
    return result
}
```

**When to use:** When the same type of object is allocated and discarded at high frequency — HTTP request buffers, JSON encoders, regexp match arrays. The pool is GC-safe: objects can be evicted at any time, so never store state you can't recreate via `New`.

### 2. Slice / Map Preallocation (`prealloc.go`)

Appending to a nil slice triggers repeated doublings of the backing array. If you know the final size, pass it as capacity to `make` and eliminate every reallocation.

```go
// Growing (multiple reallocations):
var result []int
for _, v := range src {
    result = append(result, v*2)
}

// Preallocated (zero reallocations):
result := make([]int, 0, len(src))
for _, v := range src {
    result = append(result, v*2)
}
```

The same principle applies to maps: `make(map[K]V, hint)` reduces the number of rehash operations during population.

**When to use:** Any time the final collection size is known before the loop. When it's an estimate, erring on the high side wastes a little memory but avoids all reallocations.

### 3. strings.Builder (`stringbuilder.go`)

`string` is immutable in Go. Concatenation with `+` inside a loop creates one new string per iteration — O(N²) total byte copies. `strings.Builder` accumulates bytes in a mutable buffer and allocates the final string exactly once.

```go
var b strings.Builder
b.Grow(estimatedLen)           // optional: preallocate the buffer

for i, item := range items {
    if i > 0 {
        b.WriteString(sep)
    }
    b.WriteString(item)
}
return b.String()              // one allocation
```

`fmt.Fprintf(&b, ...)` works directly on a Builder, so you can mix formatted output without extra allocations.

**When to use:** Whenever you're building a string from more than 3 pieces, or from a loop of unknown length. `strings.Join` is the stdlib shortcut for the common case of joining a slice with a separator.

### 4. Struct Field Alignment (`structalign.go`)

The compiler inserts invisible padding bytes before fields whose alignment requirement is larger than their current offset. Reordering fields from largest to smallest alignment eliminates padding and shrinks the struct.

```go
// Poor ordering — 32 bytes on 64-bit:
type Unaligned struct {
    Flag1  bool    // 1 byte + 7 bytes padding
    Count  int64   // 8 bytes
    Score  int32   // 4 bytes + 4 bytes padding
    Flag2  bool    // 1 byte + 7 bytes padding
}

// Better ordering — 16 bytes:
type Aligned struct {
    Count  int64   // 8 bytes
    Score  int32   // 4 bytes
    Flag1  bool    // 1 byte
    Flag2  bool    // 1 byte + 2 bytes padding
}
```

Use `unsafe.Sizeof` and `unsafe.Offsetof` to inspect sizes and offsets at runtime. Beyond raw size, group hot fields together so they fit in a single 64-byte cache line.

**When to use:** For structs that are allocated in large numbers (slices, maps, arena allocators) or whose hot fields are accessed on every request. Use `go vet -shadow` or `fieldalignment` (part of `golang.org/x/tools/go/analysis`) to catch suboptimal structs automatically.

### 5. Escape Analysis (`escape.go`)

The Go compiler determines at compile time whether each variable lives on the **stack** (cheap: no GC, allocated/freed with the function frame) or the **heap** (GC-managed, requires an allocation). A variable "escapes" to the heap when its lifetime may outlive the function call.

Common escape triggers:
| Trigger | Why it escapes |
|---|---|
| Returning a pointer `&x` | Address outlives the stack frame |
| Storing into an `interface{}` | Interface holds a pointer to a heap copy |
| Capturing in a returned closure | Closure lifetime is unknown at compile time |
| Passing to a function that stores the pointer | Compiler can't prove the pointer won't outlive the call |

```go
func stackAlloc(n int) int  { x := n * 2; return x  }  // stays on stack
func heapAlloc(n int) *int  { x := n * 2; return &x }  // escapes to heap
func ifaceEscape(n int) any { return n }                // int boxed onto heap
```

Observe escape decisions with:
```bash
go build -gcflags="-m" ./performance/
```

**When to use:** Escape analysis is a *diagnostic* — you read it to understand allocation behaviour, not as a design principle. Optimize only after profiling (`go test -bench . -memprofile mem.out`, then `go tool pprof`) reveals that heap allocations are the actual bottleneck.

## Running the examples

```bash
cd performance && go run .
```

## Key Takeaways

| Pattern | Rule |
|---|---|
| sync.Pool | Pool short-lived objects allocated at high frequency |
| Preallocation | Pass capacity to `make` when the final size is known |
| strings.Builder | Use instead of `+` for more than 3 concatenations or any loop |
| Struct alignment | Order fields largest-to-smallest; group hot fields together |
| Escape analysis | Read `-gcflags="-m"` output; profile before optimizing |
