package main

import "fmt"

func main() {
	// Struct embedding: promoted methods and interface satisfaction
	DemoPromotion()
	fmt.Println()

	// Overriding a promoted method (no virtual dispatch)
	DemoOverride()
	fmt.Println()

	// Embedding an interface for test spies and partial decorators
	DemoInterfaceEmbed()
	fmt.Println()

	// Embedding vs. explicit composition: when to use each
	DemoComposition()
}
