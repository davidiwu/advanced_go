package main

import "fmt"

func main() {
	// Functional Options: configure structs with composable option functions
	DemoFunctionalOptions()
	fmt.Println()

	// Option Type: represent absent values explicitly without pointer semantics
	DemoOptionType()
	fmt.Println()

	// Table-Driven Tests: define test cases as data, run with a single loop
	DemoTableDriven()
	fmt.Println()

	// Middleware Chaining: layer cross-cutting concerns around a core handler
	DemoMiddleware()
}
