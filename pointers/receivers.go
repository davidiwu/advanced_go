package main

import "fmt"

// --- Value vs Pointer Receivers ---
//
// A method can be declared with either a value receiver (func (t T) M()) or a
// pointer receiver (func (t *T) M()). The choice has two consequences:
//
//   1. Mutation: only a pointer receiver can modify the original value.
//      A value receiver operates on a copy — changes are invisible to the caller.
//
//   2. Interface satisfaction: if any method of T uses a pointer receiver,
//      the method set of *T satisfies the interface but T does not.
//      Storing a plain T value in the interface variable will not compile.
//
// Rules of thumb:
//   - Use a pointer receiver when the method mutates state, or when T is large
//     (avoid copying on every call).
//   - Use a value receiver when T is a small, immutable value type (coordinates,
//     colors, options) and the method only reads fields.
//   - Be consistent: mixing both on the same type compiles but confuses readers
//     about which methods are "safe to copy."

type counter struct {
	count int
}

// incrementValue takes a copy — the original counter is unchanged.
func (c counter) incrementValue() {
	c.count++ // mutates the local copy only
}

// incrementPointer mutates the original through the pointer.
func (c *counter) incrementPointer() {
	c.count++
}

// Stringer demonstrates the interface satisfaction rule.
type Stringer interface {
	String() string
}

type point struct {
	x, y int
}

// String uses a pointer receiver — only *point satisfies Stringer, not point.
func (p *point) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
}

func printStringer(s Stringer) {
	fmt.Println(s.String())
}

func DemoReceivers() {
	fmt.Println("=== Value vs Pointer Receivers ===")

	c := counter{count: 0}
	c.incrementValue()   // no effect on c
	fmt.Printf("after value receiver:   %d (want 0)\n", c.count)
	c.incrementPointer() // mutates c
	fmt.Printf("after pointer receiver: %d (want 1)\n", c.count)

	// Interface satisfaction
	p := &point{3, 4}
	printStringer(p) // *point satisfies Stringer

	// The compiler rejects this (uncomment to see the error):
	//   printStringer(point{3, 4})
	// → cannot use point{...} (type point) as type Stringer:
	//      point does not implement Stringer (String method has pointer receiver)

	// Go auto-dereferences for method calls on addressable values,
	// so this works even though c is not a pointer:
	c2 := counter{}
	c2.incrementPointer() // sugar for (&c2).incrementPointer()
	fmt.Printf("auto-deref on addressable value: %d (want 1)\n", c2.count)
}
