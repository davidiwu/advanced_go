# embedding

Go's primary composition mechanism: embedding types and interfaces inside structs.

## What this covers

| File | Pattern |
|---|---|
| `promotion.go` | Struct embedding and method promotion |
| `override.go` | Overriding a promoted method; no virtual dispatch |
| `interface_embed.go` | Embedding an interface for test spies and partial decorators |
| `composition.go` | Embedding vs. explicit composition — when to use each |

## Key ideas

### Method promotion

Embedding a type `T` in a struct promotes all of `T`'s exported methods to the outer struct. The outer struct automatically satisfies any interface `T` satisfies.

```go
type Animal struct{ Name string }
func (a Animal) Speak() string { return a.Name + " speaks" }

type Dog struct {
    Animal        // embedded — Speak is promoted
    Breed string
}

d := Dog{Animal: Animal{Name: "Rex"}, Breed: "Lab"}
d.Speak()        // calls Animal.Speak directly
```

### No virtual dispatch

When the outer struct overrides a promoted method, the promoted method does **not** automatically call the outer version. Embedding is field access, not inheritance.

```go
type PrefixLogger struct {
    Logger
    Prefix string
}
func (p PrefixLogger) Log(msg string) { /* own impl */ }

pl.Log("x")        // → PrefixLogger.Log (outer wins)
pl.Warn("y")       // → Logger.Warn, which calls Logger.Log (not PrefixLogger.Log)
pl.Logger.Log("z") // reach the embedded method explicitly
```

### Interface embedding in a struct

A struct can embed an interface to satisfy it without implementing every method. Unimplemented methods panic at runtime — which is intentional in test doubles.

```go
// Test spy: only override the method under test.
type spyNotifier struct {
    Notifier       // satisfies the interface; unimplemented methods panic
    sent []string
}
func (s *spyNotifier) Send(msg string) error {
    s.sent = append(s.sent, msg)
    return nil
}
```

### Embedding vs. explicit composition

| | Embedding | Explicit field |
|---|---|---|
| Method promotion | Yes | No |
| Interface satisfaction | Automatic | Manual (forward methods) |
| Relationship | IS-A | HAS-A |
| Multiple same-type fields | Only one | Unlimited |
| Call site | `s.Method()` | `s.field.Method()` |
