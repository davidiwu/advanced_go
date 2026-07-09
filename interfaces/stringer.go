package main

import (
	"errors"
	"fmt"
)

// --- Stringer and Error Interfaces ---
//
// Two of Go's most used built-in interfaces:
//
//   fmt.Stringer  — any type with String() string controls its own fmt output.
//   error         — any type with Error() string is a Go error.
//
// Implementing these lets your types integrate naturally with fmt, log,
// errors.As/Is, and the rest of the standard library without any glue code.

// Direction is an enum-style type that prints as a readable string instead
// of a raw integer when used with fmt.Println, %v, or %s.
type Direction int

const (
	North Direction = iota
	East
	South
	West
)

// String satisfies fmt.Stringer. fmt will call this automatically for %v and %s.
func (d Direction) String() string {
	switch d {
	case North:
		return "North"
	case East:
		return "East"
	case South:
		return "South"
	case West:
		return "West"
	default:
		return fmt.Sprintf("Direction(%d)", int(d))
	}
}

// --- Custom error types ---

// ValidationError carries structured context about what failed.
// Callers can use errors.As to extract the details, rather than
// parsing an error string.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s — %s", e.Field, e.Message)
}

// NotFoundError is a sentinel-style error type.
// Wrapping it with fmt.Errorf("%w", ...) preserves the chain for errors.Is/As.
type NotFoundError struct {
	Resource string
	ID       int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id=%d not found", e.Resource, e.ID)
}

// lookup returns a wrapped NotFoundError so callers can use errors.As
// to get the full NotFoundError value, or errors.Is against a sentinel.
func lookup(id int) error {
	if id != 42 {
		return fmt.Errorf("lookup: %w", &NotFoundError{Resource: "user", ID: id})
	}
	return nil
}

// validate returns a *ValidationError when input is blank.
func validate(name string) error {
	if name == "" {
		return &ValidationError{Field: "name", Message: "must not be empty"}
	}
	return nil
}

// DemoStringer shows fmt.Stringer in action and demonstrates extracting
// concrete error types with errors.As.
func DemoStringer() {
	fmt.Println("=== Stringer and Error Interfaces ===")

	// fmt.Stringer: directions print as words, not numbers.
	dirs := []Direction{North, East, South, West}
	fmt.Println("  directions:", dirs)
	fmt.Printf("  heading: %v\n", North)

	// Custom errors: errors.As lets callers access structured fields
	// without parsing the error string.
	fmt.Println("  --- custom errors ---")

	if err := lookup(99); err != nil {
		var nfe *NotFoundError
		if errors.As(err, &nfe) {
			fmt.Printf("  not found: resource=%s id=%d\n", nfe.Resource, nfe.ID)
		}
		fmt.Println("  full error:", err)
	}

	if err := validate(""); err != nil {
		var ve *ValidationError
		if errors.As(err, &ve) {
			fmt.Printf("  validation: field=%s message=%s\n", ve.Field, ve.Message)
		}
	}

	// No error path
	if err := lookup(42); err == nil {
		fmt.Println("  lookup(42): ok")
	}
}
