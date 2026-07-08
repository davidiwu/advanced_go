package main

import "fmt"

// Ordered is a constraint interface that restricts T to types that support the
// < and > operators. "any" would be too broad here — you can't compare arbitrary
// types with < , so the compiler would reject Min/Max if T were unconstrained.
// The ~ prefix means "any type whose underlying type is X", which allows custom
// named types like Celsius (defined below) to satisfy the constraint.
type Ordered interface {
	~int | ~float64 | ~string
}

// Without the Ordered constraint, the expression "a < b" would not compile because
// Go cannot guarantee that an arbitrary T supports comparison operators.
func Min[T Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Clamp reuses Min and Max — both accept any T Ordered, so they compose naturally.
func Clamp[T Ordered](v, lo, hi T) T {
	return Max(lo, Min(v, hi))
}

// Celsius has float64 as its underlying type, so ~float64 in Ordered covers it.
// Without the ~ prefix, Celsius would NOT satisfy the constraint even though it
// behaves exactly like float64 numerically.
type Celsius float64

func DemoOrdered() {
	fmt.Println("=== Ordered Constraint ===")

	fmt.Println("min(3, 7):", Min(3, 7))
	fmt.Println("max(3.14, 2.72):", Max(3.14, 2.72))
	fmt.Println(`min("apple","banana"):`, Min("apple", "banana"))

	var temp Celsius = 120.0
	clamped := Clamp(temp, Celsius(0), Celsius(100))
	fmt.Println("clamped temp:", clamped)
}
