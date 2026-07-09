package main

import "fmt"

func main() {
	// Implicit Satisfaction: types satisfy interfaces without declaring it
	DemoImplicit()
	fmt.Println()

	// Interface Composition: build larger contracts from smaller ones
	DemoComposition()
	fmt.Println()

	// Type Assertions and Switches: inspect and extract concrete types at runtime
	DemoTypeAssert()
	fmt.Println()

	// Interface Segregation: prefer many small interfaces over one large one
	DemoSegregation()
	fmt.Println()

	// Stringer and Error: built-in interfaces for formatting and error handling
	DemoStringer()
	fmt.Println()

	// Accept Interfaces, Return Structs: the canonical Go API design principle
	DemoAcceptReturn()
}
