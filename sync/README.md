# sync — Shared-Memory Concurrency Primitives

This module demonstrates the standard library tools for protecting shared state
in Go: mutexes, atomic operations, lazy initialization, and bounded concurrency.

## Patterns

### `sync.Mutex` (`mutex.go`)

Protects a shared variable so only one goroutine can access it at a time.

```go
type safeCounter struct {
    mu    sync.Mutex
    value int
}

func (c *safeCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}
```

**When to use:** Any time multiple goroutines read _and_ write the same variable.  
**Key point:** `defer mu.Unlock()` ensures the lock is always released, even on early return.

---

### `sync.RWMutex` (`rwmutex.go`)

Allows many concurrent readers but only one writer at a time.

```go
func (c *cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.items[key]
}

func (c *cache) Set(key, val string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = val
}
```

**When to use:** Read-heavy workloads (e.g. in-memory caches) where serializing all reads would be wasteful.  
**Key point:** If writes are frequent, a plain `Mutex` is usually simpler and equally fast.

---

### `sync.Once` (`once.go`)

Executes a function exactly once regardless of how many goroutines call `Do`.

```go
var (
    instance *dbConn
    once     sync.Once
)

func getDB() *dbConn {
    once.Do(func() {
        instance = &dbConn{dsn: "postgres://localhost/mydb"}
    })
    return instance
}
```

**When to use:** Lazy initialization of expensive resources (DB connections, config loading).  
**Key point:** The zero value of `sync.Once` is ready to use — no constructor needed.

---

### `sync/atomic` (`atomic.go`)

Lock-free integer operations backed by CPU-level instructions.

```go
var hits atomic.Int64

hits.Add(1)           // increment
hits.Load()           // read
hits.CompareAndSwap(0, 1) // only update if current value matches expected
```

**When to use:** A single shared integer counter; avoid for anything more complex.  
**Key point:** ~10× faster than a Mutex for simple counters, but easy to misuse. Prefer Mutex when protecting multiple fields or a struct.

---

### Channel-as-Semaphore (`semaphore.go`)

A buffered channel of capacity N limits how many goroutines run concurrently.

```go
sem := make(chan struct{}, 3) // allow at most 3 concurrent workers

go func() {
    sem <- struct{}{}         // acquire — blocks when full
    defer func() { <-sem }() // release
    doWork()
}()
```

**When to use:** Fan out to many goroutines but cap the number doing an expensive operation (outbound HTTP, file I/O) at any moment.  
**Key point:** Simpler than a full worker pool when you don't need persistent workers or a job queue.

---

## Running

```bash
go run ./sync/
```
