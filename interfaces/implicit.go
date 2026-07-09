package main

import "fmt"

// --- Implicit Interface Satisfaction ---
//
// Go interfaces are satisfied implicitly: no "implements" keyword needed.
// A type satisfies an interface the moment it has all the required methods.
//
// This decouples producers from consumers — the package that defines a type
// never needs to import or know about the interfaces it satisfies. You can
// retrofit an interface onto an existing type without touching its source.

// Animal is a simple interface. Any type with a Sound() method satisfies it,
// regardless of what package it lives in.
type Animal interface {
	Sound() string
	Name() string
}

type Dog struct{ name string }
type Cat struct{ name string }
type Parrot struct{ name string }

func (d Dog) Sound() string    { return "Woof" }
func (d Dog) Name() string     { return d.name }
func (c Cat) Sound() string    { return "Meow" }
func (c Cat) Name() string     { return c.name }
func (p Parrot) Sound() string { return "Squawk" }
func (p Parrot) Name() string  { return p.name }

// makeNoise works with any Animal — it has no knowledge of Dog, Cat, or Parrot.
// Adding a new animal type requires zero changes here.
func makeNoise(a Animal) {
	fmt.Printf("%s says: %s\n", a.Name(), a.Sound())
}

// DemoImplicit shows that all three types satisfy Animal without any
// explicit declaration, and that a single function works on all of them.
func DemoImplicit() {
	fmt.Println("=== Implicit Interface Satisfaction ===")

	animals := []Animal{
		Dog{name: "Rex"},
		Cat{name: "Whiskers"},
		Parrot{name: "Polly"},
	}
	for _, a := range animals {
		makeNoise(a)
	}
}
