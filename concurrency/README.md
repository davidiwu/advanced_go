# Concurrency Patterns in Go

Go's concurrency model is built on goroutines and channels (CSP — Communicating Sequential Processes). The patterns below are the idiomatic building blocks for real-world concurrent programs.

---

## Table of Contents

1. [Pipeline](#1-pipeline)
2. [Fan-out / Fan-in](#2-fan-out--fan-in)
3. [Worker Pool](#3-worker-pool)
4. [errgroup](#4-errgroup)
5. [Context Propagation](#5-context-propagation)

---

## 1. Pipeline

### What it is
A pipeline is a series of stages connected by channels. Each stage receives values from upstream, transforms them, and sends results downstream. Each stage runs in its own goroutine.

### When to use
- ETL-style data processing (read → transform → write)
- Streaming workloads where you want to process items as they arrive, not all at once
- When you want separation of concerns — each stage has one job

### How it works

```
source → stage1 → stage2 → stage3 → sink
```

Each stage looks like:

```go
func stage(in <-chan T) <-chan U {
    out := make(chan U)
    go func() {
        defer close(out)
        for v := range in {
            out <- transform(v)
        }
    }()
    return out
}
```

The caller wires stages together:

```go
nums   := generate(2, 3, 4, 5)     // source
sq     := square(nums)              // stage 1
result := filter(sq)               // stage 2
for v := range result {
    fmt.Println(v)
}
```

### Key points
- Close the output channel when done (`defer close(out)`) — this signals the next stage to stop ranging.
- Pass a `context.Context` to each stage for cancellation (see [Context Propagation](#5-context-propagation)).
- Unbuffered channels give you backpressure: a slow downstream stage naturally slows the upstream ones.

---

## 2. Fan-out / Fan-in

### What it is
- **Fan-out**: spread work from one channel across multiple goroutines working in parallel.
- **Fan-in**: merge multiple output channels into a single channel for the consumer.

### When to use
- CPU-bound or I/O-bound work that can be parallelized (e.g., making N HTTP calls, processing N files)
- When a single pipeline stage is a bottleneck — fan out that stage across `runtime.NumCPU()` workers
- Aggregating results from independent goroutines

### How it works

```
         ┌─ worker1 ─┐
source ──┤─ worker2 ─├──► merged output
         └─ worker3 ─┘
```

**Fan-out** — launch N goroutines reading from the same input channel:

```go
func fanOut(in <-chan Job, n int) []<-chan Result {
    channels := make([]<-chan Result, n)
    for i := 0; i < n; i++ {
        channels[i] = process(in) // each goroutine reads from the shared 'in'
    }
    return channels
}
```

**Fan-in** — merge N channels into one using a goroutine per channel + a WaitGroup:

```go
func fanIn(channels ...<-chan Result) <-chan Result {
    out := make(chan Result)
    var wg sync.WaitGroup
    wg.Add(len(channels))
    for _, ch := range channels {
        go func(c <-chan Result) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(ch)
    }
    go func() { wg.Wait(); close(out) }()
    return out
}
```

### Key points
- Output order is non-deterministic — if order matters, fan out but collect results into a slice indexed by position.
- Always close the merged channel after all sources are drained (`wg.Wait()` then `close`).
- Fan-out degree should typically be `runtime.NumCPU()` for CPU-bound work, or higher for I/O-bound work.

---

## 3. Worker Pool

### What it is
A fixed pool of N goroutines that all read from a shared job channel. The pool size caps concurrency, preventing resource exhaustion.

### When to use
- You have a large or unbounded number of tasks but want to limit how many run concurrently
- Controlling goroutine count when dealing with limited resources (DB connections, file handles, rate-limited APIs)
- Long-running services that process a stream of incoming work

### How it works

```
jobs channel ──► [worker 1]
              ──► [worker 2]  ──► results channel
              ──► [worker 3]
```

```go
func startPool(numWorkers int, jobs <-chan Job) <-chan Result {
    results := make(chan Result)
    var wg sync.WaitGroup
    wg.Add(numWorkers)
    for i := 0; i < numWorkers; i++ {
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }
    go func() { wg.Wait(); close(results) }()
    return results
}

// Caller:
jobs := make(chan Job, 100)
results := startPool(runtime.NumCPU(), jobs)

for _, j := range allJobs {
    jobs <- j
}
close(jobs) // signal workers that no more jobs are coming

for r := range results {
    fmt.Println(r)
}
```

### Key points
- The caller closes `jobs` to signal workers to stop — workers see `range jobs` end naturally.
- Buffer the jobs channel if producers can burst ahead of workers.
- Combine with `context.Context` so the pool can be cancelled mid-flight.
- Difference from fan-out: fan-out wires goroutines to individual channels; a worker pool uses a single shared channel.

---

## 4. errgroup

### What it is
`golang.org/x/sync/errgroup` coordinates a group of goroutines and collects the first non-nil error. It's a higher-level replacement for a manual `sync.WaitGroup` + error channel when you need error propagation.

### When to use
- Launching N independent goroutines and failing fast on the first error
- Replacing boilerplate `WaitGroup` + `errCh` patterns
- When you want automatic context cancellation on first error (`errgroup.WithContext`)

### How it works

```go
import "golang.org/x/sync/errgroup"

g, ctx := errgroup.WithContext(context.Background())

for _, url := range urls {
    url := url // capture loop variable
    g.Go(func() error {
        return fetch(ctx, url) // ctx is cancelled if any goroutine returns an error
    })
}

if err := g.Wait(); err != nil {
    log.Fatal(err) // first non-nil error from any goroutine
}
```

### Key points
- `g.Wait()` blocks until all goroutines finish and returns the **first** error (not all errors).
- `errgroup.WithContext` cancels the derived `ctx` when the first error occurs — pass it into your goroutines so they can exit early.
- Use `SetLimit(n)` (Go 1.20+) to cap how many goroutines run at once.
- For collecting *all* errors (not just the first), you need a different approach — e.g., a mutex-guarded slice.

---

## 5. Context Propagation

### What it is
`context.Context` carries deadlines, cancellation signals, and request-scoped values across goroutines and API boundaries. It is the standard mechanism for controlling goroutine lifetime in Go.

### When to use
- Any concurrent operation that should be cancellable (user request, timeout, shutdown signal)
- Propagating deadlines from an HTTP handler into downstream goroutines, DB calls, or RPCs
- Coordinating a tree of goroutines: cancelling the root cancels all children

### How it works

**Timeout:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel() // always call cancel to release resources

result, err := doWork(ctx)
```

**Manual cancellation:**
```go
ctx, cancel := context.WithCancel(context.Background())
go func() {
    <-shutdownSignal
    cancel() // triggers cancellation for all goroutines using this ctx
}()

doWork(ctx)
```

**Inside a goroutine — checking for cancellation:**
```go
func doWork(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err() // context.Canceled or context.DeadlineExceeded
        case job := <-jobs:
            process(job)
        }
    }
}
```

**Values (use sparingly):**
```go
type key struct{}
ctx = context.WithValue(ctx, key{}, "request-id-123")
val := ctx.Value(key{}).(string)
```

### Key points
- Always pass `ctx` as the **first parameter** of functions that do I/O or spawn goroutines — this is a strong Go convention.
- Always `defer cancel()` immediately after creating a cancelable context — even if you expect it to time out, this ensures resources are freed.
- Check `ctx.Done()` in long-running loops via a `select` statement so goroutines respond promptly to cancellation.
- Use `context.WithValue` only for request-scoped metadata (trace IDs, auth tokens) — never for optional function parameters.
- The zero value (`context.Background()`) is the root; `context.TODO()` signals "I haven't wired this up yet."

---

## Combining Patterns

These patterns compose naturally:

| Combination | Use case |
|---|---|
| Pipeline + Context | Cancellable data processing pipeline |
| Fan-out + errgroup | Parallel I/O (HTTP, disk) with fail-fast error handling |
| Worker Pool + Context | Bounded concurrency with graceful shutdown |
| Fan-out + Fan-in + Context | Parallel workers, merged output, cancellable |

A typical production pattern: a worker pool reading from a job channel, each worker propagating a shared context, with an errgroup coordinating errors and cancellation.

---

## Further Reading

- [Go Blog: Pipelines and Cancellation](https://go.dev/blog/pipelines)
- [Go Blog: Context](https://go.dev/blog/context)
- [`golang.org/x/sync/errgroup`](https://pkg.go.dev/golang.org/x/sync/errgroup)
- [Go Concurrency Patterns (Rob Pike, Google I/O 2012)](https://talks.golang.org/2012/concurrency.slide)
- [Advanced Go Concurrency Patterns (Sameer Ajmani, Google I/O 2013)](https://talks.golang.org/2013/advconc.slide)
