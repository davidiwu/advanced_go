package main

import "fmt"

// ========================================================================
// PATTERN: Embedding an interface in a struct
//
// A struct can embed an interface type. This means the struct automatically
// satisfies that interface — but any unimplemented method panics at runtime
// if called. The value of this pattern is not in calling every method; it is
// in overriding only the methods you care about in tests or adapters.
//
// Two common uses:
//   1. Test spy: embed the interface, override only the methods under test,
//      let everything else panic (which reveals unexpected calls).
//   2. Partial decorator: wrap a real implementation, intercept a few methods,
//      delegate the rest through the embedded field.
// ========================================================================

// Notifier is the interface we want to partially implement.
type Notifier interface {
	Send(msg string) error
	Ping() bool
	Close() error
}

// realNotifier is the production implementation (simplified here).
type realNotifier struct{ addr string }

func (r *realNotifier) Send(msg string) error {
	fmt.Printf("  [real] sending %q to %s\n", msg, r.addr)
	return nil
}
func (r *realNotifier) Ping() bool  { return true }
func (r *realNotifier) Close() error { return nil }

// --- Test spy pattern ---

// spyNotifier embeds Notifier so it satisfies the interface with zero
// boilerplate. Only Send is overridden; Ping and Close would panic if called
// (which is intentional — it signals the test has unexpected side-effects).
type spyNotifier struct {
	Notifier // embedded interface — satisfies Notifier without full impl
	sent     []string
}

func (s *spyNotifier) Send(msg string) error {
	s.sent = append(s.sent, msg)
	return nil
}

// --- Partial decorator pattern ---

// loggingNotifier wraps a real Notifier and intercepts Send to add logging.
// Ping and Close are delegated through the embedded field automatically.
type loggingNotifier struct {
	Notifier // holds the real implementation
}

func (l *loggingNotifier) Send(msg string) error {
	fmt.Printf("  [log] about to send: %q\n", msg)
	err := l.Notifier.Send(msg) // delegate to wrapped impl
	fmt.Printf("  [log] send complete, err=%v\n", err)
	return err
}

func DemoInterfaceEmbed() {
	fmt.Println("=== Embedding an Interface in a Struct ===")

	// Test spy: capture calls, never touch a real network.
	spy := &spyNotifier{}
	var n Notifier = spy
	_ = n.Send("hello")
	_ = n.Send("world")
	fmt.Printf("spy captured %d messages: %v\n", len(spy.sent), spy.sent)

	// Partial decorator: log around a real notifier.
	real := &realNotifier{addr: "localhost:9000"}
	logged := &loggingNotifier{Notifier: real}
	_ = logged.Send("deploy complete")
	fmt.Printf("ping through decorator: %v\n", logged.Ping()) // delegated
}
