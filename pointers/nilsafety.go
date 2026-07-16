package main

import "fmt"

// --- Nil Safety ---
//
// The zero value of any pointer type is nil. Dereferencing nil panics at runtime;
// the compiler does not catch it. Safe pointer code follows two rules:
//
//   1. Validate at boundaries — check for nil before dereferencing any pointer
//      that originates outside your package (function arguments, interface values,
//      decoded structs from JSON/gRPC, etc.).
//   2. Nil receivers are valid — a method with a pointer receiver can be called on
//      a nil pointer without panicking, as long as the method itself guards against it.
//      This is idiomatic for optional / chainable APIs.
//
// When to use nil as a sentinel:
//   - Return (*T, error): nil pointer signals "no value" alongside an explicit error.
//   - Optional config fields: a nil pointer means "not set"; a non-nil pointer means "set to this value".
//   - Linked lists, trees: nil marks the end of a chain.

// Config demonstrates using a nil pointer field as an optional setting.
type Config struct {
	Timeout *int // nil means "use default"; non-nil means "caller specified this"
}

func effectiveTimeout(c Config) int {
	if c.Timeout == nil {
		return 30 // default
	}
	return *c.Timeout
}

// Node is a simple linked list node where nil terminates the chain.
type Node struct {
	Value int
	Next  *Node
}

// String is safe to call on a nil *Node — nil represents the empty list.
func (n *Node) String() string {
	if n == nil {
		return "[]"
	}
	result := fmt.Sprintf("[%d", n.Value)
	cur := n.Next
	for cur != nil {
		result += fmt.Sprintf(" %d", cur.Value)
		cur = cur.Next
	}
	return result + "]"
}

// Len returns the length of the list. Safe to call on nil.
func (n *Node) Len() int {
	if n == nil {
		return 0
	}
	return 1 + n.Next.Len()
}

func DemoNilSafety() {
	fmt.Println("=== Nil Safety ===")

	// Optional field: nil means "not set"
	defaultCfg := Config{}
	customTimeout := 10
	customCfg := Config{Timeout: &customTimeout}
	fmt.Printf("default timeout: %d\n", effectiveTimeout(defaultCfg))
	fmt.Printf("custom timeout:  %d\n", effectiveTimeout(customCfg))

	// Nil receiver method: calling String/Len on nil is safe
	var empty *Node
	fmt.Printf("nil node string: %s, len: %d\n", empty.String(), empty.Len())

	list := &Node{1, &Node{2, &Node{3, nil}}}
	fmt.Printf("list string: %s, len: %d\n", list.String(), list.Len())
}
