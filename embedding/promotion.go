package main

import "fmt"

// ========================================================================
// PATTERN: Struct embedding for method promotion
//
// When you embed a type T inside a struct, all of T's exported methods and
// fields are "promoted" to the outer struct. The outer struct satisfies any
// interface T satisfies — for free, with no forwarding boilerplate.
//
// Use this to share behaviour across types without inheritance.
// ========================================================================

type Animal struct {
	Name string
}

func (a Animal) Speak() string {
	return fmt.Sprintf("%s makes a sound", a.Name)
}

func (a Animal) Breathe() string {
	return fmt.Sprintf("%s breathes air", a.Name)
}

// Dog embeds Animal. It automatically promotes Speak and Breathe.
// Dog satisfies any interface that Animal satisfies.
type Dog struct {
	Animal
	Breed string
}

// Cat also embeds Animal, gaining the same promoted methods.
type Cat struct {
	Animal
	Indoor bool
}

func DemoPromotion() {
	fmt.Println("=== Struct Embedding: Method Promotion ===")

	dog := Dog{Animal: Animal{Name: "Rex"}, Breed: "Labrador"}
	cat := Cat{Animal: Animal{Name: "Whiskers"}, Indoor: true}

	// Promoted methods called directly on the outer type.
	fmt.Println(dog.Speak())   // promoted from Animal
	fmt.Println(dog.Breathe()) // promoted from Animal
	fmt.Printf("Breed: %s\n", dog.Breed)

	fmt.Println(cat.Speak())
	fmt.Printf("Indoor: %v\n", cat.Indoor)

	// The embedded field is also accessible by name when needed.
	fmt.Printf("Embedded field: %+v\n", dog.Animal)
}
