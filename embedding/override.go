package main

import "fmt"

// ========================================================================
// PATTERN: Overriding promoted methods
//
// If the outer struct defines a method with the same name as a promoted one,
// the outer method wins in a direct call. The promoted method is still
// reachable through the embedded field name.
//
// This is Go's analogue of method overriding, with an important difference:
// the embedded type's method does NOT automatically call the outer one (no
// dynamic dispatch). It is plain field access, not virtual dispatch.
// ========================================================================

type Logger struct{}

func (l Logger) Log(msg string) {
	fmt.Printf("[base] %s\n", msg)
}

func (l Logger) Warn(msg string) {
	l.Log("WARN: " + msg)
}

// PrefixLogger embeds Logger and overrides Log to prepend a prefix.
// Warn is still promoted unchanged — but it calls Logger.Log, not the
// overriding PrefixLogger.Log. That is the "no virtual dispatch" trade-off.
type PrefixLogger struct {
	Logger
	Prefix string
}

func (p PrefixLogger) Log(msg string) {
	fmt.Printf("[%s] %s\n", p.Prefix, msg)
}

func DemoOverride() {
	fmt.Println("=== Overriding Promoted Methods ===")

	pl := PrefixLogger{Prefix: "app"}

	// Calls PrefixLogger.Log — outer method wins.
	pl.Log("server started")

	// Calls Logger.Warn, which internally calls Logger.Log (not PrefixLogger.Log).
	// There is no dynamic dispatch: embedding is field access, not inheritance.
	pl.Warn("disk space low")

	// Reach the shadowed method explicitly via the embedded field name.
	pl.Logger.Log("direct call to embedded Log")
}
