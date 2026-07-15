# io — Custom Readers, Writers, and Stream Pipelines

Go's `io` package is built around two interfaces — `io.Reader` and `io.Writer` — that nearly every I/O function in the stdlib accepts. This module shows how to implement them on custom types, how to stack decorators, how to connect a writer to a reader with `io.Pipe`, and how to tokenize streams with `bufio.Scanner`.

## Patterns

### Custom `io.Reader` (`reader.go`)

Implement `Read(p []byte) (n int, err error)` to make any type a source of bytes.

```go
type RepeatReader struct {
    data   string
    repeat int
    pos    int
}

func (r *RepeatReader) Read(p []byte) (int, error) {
    total := len(r.data) * r.repeat
    if r.pos >= total {
        return 0, io.EOF
    }
    n := 0
    for n < len(p) && r.pos < total {
        p[n] = r.data[r.pos%len(r.data)]
        r.pos++
        n++
    }
    if r.pos >= total {
        return n, io.EOF
    }
    return n, nil
}
```

**When to use:** When your data source doesn't come from a file, network, or string — e.g. generated data, in-memory structures, protocol encoding.  
**Key point:** Return `io.EOF` alongside the last bytes in a single call — callers (including `io.ReadAll`) handle that correctly. Never return `(0, nil)` unless you have no data yet and expect more later.

---

### Custom `io.Writer` (`writer.go`)

Implement `Write(p []byte) (n int, err error)` to capture, transform, or route written bytes.

```go
type CountingWriter struct {
    w     io.Writer
    total int64
}

func (c *CountingWriter) Write(p []byte) (int, error) {
    n, err := c.w.Write(p)
    c.total += int64(n)
    return n, err
}
```

`io.MultiWriter` fans a single write out to multiple destinations simultaneously:

```go
mw := io.MultiWriter(fileWriter, hashWriter, logWriter)
fmt.Fprintf(mw, "written once, received by all three")
```

**When to use:** Byte counting, quota enforcement, test capture, or fan-out to multiple sinks.  
**Key point:** Always return the exact number of bytes written. Returning fewer than `len(p)` without an error violates the `io.Writer` contract and will cause `io.Copy` to return `io.ErrShortWrite`.

---

### Decorator Pattern (`decorator.go`)

A decorator wraps an existing `Reader` or `Writer` and adds behaviour without changing the interface.

```go
// UpperReader transforms bytes as they pass through.
type UpperReader struct{ r io.Reader }

func (u *UpperReader) Read(p []byte) (int, error) {
    n, err := u.r.Read(p)
    for i := range n {
        if p[i] >= 'a' && p[i] <= 'z' {
            p[i] -= 32
        }
    }
    return n, err
}
```

The stdlib ships several ready-made decorators:

| Function | What it does |
|---|---|
| `io.LimitReader(r, n)` | truncates `r` to at most `n` bytes |
| `io.TeeReader(r, w)` | reads from `r` while simultaneously writing to `w` |
| `io.MultiReader(rs...)` | concatenates readers end-to-end |
| `io.MultiWriter(ws...)` | fans writes out to all writers |

**When to use:** Any time you need to add logging, compression, encryption, or transformation to an existing stream without changing the producer or consumer.  
**Key point:** Decorators compose. `NewUpperReader(io.LimitReader(r, 100))` chains two transformations over the same underlying stream.

---

### `io.Pipe` (`pipe.go`)

`io.Pipe` creates a synchronous, in-memory connection between a writer and a reader with no buffering. The writer blocks until the reader consumes.

```go
pr, pw := io.Pipe()

go func() {
    defer pw.Close()
    fmt.Fprint(pw, "hello through the pipe")
}()

data, _ := io.ReadAll(pr) // blocks until the goroutine writes and closes
```

Errors propagate in both directions:

```go
pw.CloseWithError(fmt.Errorf("upstream failure"))
// The reader receives that error from its next Read call.
```

**When to use:** Connecting a function that writes (e.g. `json.Encoder`, `gzip.Writer`) to a function that reads (e.g. `http.Request.Body`) without allocating an intermediate buffer.  
**Key point:** Because there is no buffer, writer and reader must run in separate goroutines or the write will deadlock waiting for a read that never comes.

---

### `bufio.Scanner` (`scanner.go`)

`Scanner` wraps any `io.Reader` and tokenizes it according to a split function.

```go
scanner := bufio.NewScanner(r)
for scanner.Scan() {
    fmt.Println(scanner.Text()) // one token per iteration
}
if err := scanner.Err(); err != nil { ... }
```

Built-in split functions: `bufio.ScanLines` (default), `bufio.ScanWords`, `bufio.ScanBytes`, `bufio.ScanRunes`.

Custom split functions follow the signature `func(data []byte, atEOF bool) (advance int, token []byte, err error)`:

```go
// Split on commas, trimming surrounding whitespace.
csv.Split(func(data []byte, atEOF bool) (int, []byte, error) {
    for i, b := range data {
        if b == ',' {
            return i + 1, bytes.TrimSpace(data[:i]), nil
        }
    }
    if atEOF && len(data) > 0 {
        return len(data), bytes.TrimSpace(data), nil
    }
    return 0, nil, nil // request more data
})
```

For tokens larger than 64 KB, call `scanner.Buffer` before the first `Scan`:

```go
scanner.Buffer(make([]byte, maxSize), maxSize)
```

**When to use:** Line-by-line log processing, CSV parsing, or any stream where you want to iterate over tokens without loading the entire input into memory.  
**Key point:** Always check `scanner.Err()` after the loop. A `false` return from `Scan` means either EOF (no error) or an error — they look the same without the explicit check.

---

## Running

```bash
cd io && go run .
```

## Key Takeaways

| Concept | Rule |
|---|---|
| `io.Reader` | Return `io.EOF` with the last bytes; never return `(0, nil)` at the end |
| `io.Writer` | Return exactly `len(p)` bytes written or an error |
| Decorators | Compose by wrapping — same interface in, same interface out |
| `io.Pipe` | Writer and reader must be in separate goroutines; use `CloseWithError` to propagate failures |
| `bufio.Scanner` | Always check `Err()` after the loop; call `Buffer()` for tokens > 64 KB |
