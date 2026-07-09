package main

import "fmt"

// --- Type Assertions and Type Switches ---
//
// A type assertion extracts the concrete value from an interface variable.
// Use the two-value form (v, ok) to avoid a panic when the assertion fails.
//
// A type switch is a cleaner way to branch on the concrete type of an
// interface value. It's the idiomatic alternative to a chain of if/else
// type assertions.

// Shape is intentionally minimal so we can hold different concrete types
// in the same slice and inspect them at runtime.
type Shape interface {
	Area() float64
}

type Circle struct{ Radius float64 }
type Rectangle struct{ Width, Height float64 }
type Triangle struct{ Base, Height float64 }

func (c Circle) Area() float64    { return 3.14159 * c.Radius * c.Radius }
func (r Rectangle) Area() float64 { return r.Width * r.Height }
func (t Triangle) Area() float64  { return 0.5 * t.Base * t.Height }

// describe uses a type switch to print type-specific details that aren't
// part of the Shape interface. This is the right tool when you genuinely
// need to diverge behavior based on the concrete type — use it sparingly,
// since it creates a coupling to every listed type.
func describe(s Shape) string {
	switch v := s.(type) {
	case Circle:
		return fmt.Sprintf("circle with radius %.2f, area=%.2f", v.Radius, v.Area())
	case Rectangle:
		return fmt.Sprintf("rectangle %.2fx%.2f, area=%.2f", v.Width, v.Height, v.Area())
	case Triangle:
		return fmt.Sprintf("triangle base=%.2f height=%.2f, area=%.2f", v.Base, v.Height, v.Area())
	default:
		// default keeps the switch safe as new Shape types are added.
		return fmt.Sprintf("unknown shape, area=%.2f", v.Area())
	}
}

// extractCircle shows the two-value type assertion form.
// The single-value form (s.(Circle)) panics if s is not a Circle;
// always prefer the ok form unless you are certain of the type.
func extractCircle(s Shape) {
	if c, ok := s.(Circle); ok {
		fmt.Printf("  got a circle: radius=%.2f\n", c.Radius)
	} else {
		fmt.Println("  not a circle")
	}
}

// DemoTypeAssert shows both type assertions and type switches in action.
func DemoTypeAssert() {
	fmt.Println("=== Type Assertions and Type Switches ===")

	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 4, Height: 6},
		Triangle{Base: 3, Height: 8},
	}

	for _, s := range shapes {
		fmt.Println(" ", describe(s))
	}

	fmt.Println("  --- two-value type assertion ---")
	extractCircle(Circle{Radius: 2.5})
	extractCircle(Rectangle{Width: 1, Height: 1})
}
