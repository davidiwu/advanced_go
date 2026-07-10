package main

import (
	"cmp"
	"fmt"
)

// cmp.Ordered (Go 1.21) is the stdlib constraint for types that support <, >.
// It covers all integer types, float types, and string — with ~ so named types
// like Celsius below satisfy it automatically.
// Before 1.21 you had to declare this constraint yourself.

// Clamp uses the built-in min/max (Go 1.21) directly.
// The separate Min[T]/Max[T] wrapper functions this file used to define are
// now redundant — min and max are generic built-ins that work on any cmp.Ordered type.
func Clamp[T cmp.Ordered](v, lo, hi T) T {
	return max(lo, min(v, hi))
}

// Celsius has float64 as its underlying type, so cmp.Ordered covers it via ~float64.
type Celsius float64

func DemoOrdered() {
	fmt.Println("=== Ordered Constraint ===")

	// Built-in min/max work on any ordered type — no wrapper needed.
	fmt.Println("min(3, 7):", min(3, 7))
	fmt.Println("max(3.14, 2.72):", max(3.14, 2.72))
	fmt.Println(`min("apple","banana"):`, min("apple", "banana"))

	var temp Celsius = 120.0
	clamped := Clamp(temp, Celsius(0), Celsius(100))
	fmt.Println("clamped temp:", clamped)
}
