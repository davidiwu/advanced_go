package main

import "fmt"

// Stack is a generic data structure. Without generics, you would need a separate
// IntStack, StringStack, etc., or use interface{} and lose type safety (requiring
// type assertions on every Pop call). With [T any], one Stack works for all types
// and the compiler enforces that you don't mix types within a single instance.
type Stack[T any] struct {
	items []T
}

// T is not re-declared here — it is already in scope from the type definition above.
// The receiver *Stack[T] borrows T; no constraint is needed because it was set to
// "any" when the type was declared.
func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		// "var zero T" is the idiomatic way to get the zero value of a type parameter,
		// since we can't write nil or 0 — we don't know T's concrete type at this point.
		var zero T
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}

func DemoStack() {
	fmt.Println("=== Generic Stack ===")

	var intStack Stack[int]
	intStack.Push(1)
	intStack.Push(2)
	intStack.Push(3)
	for intStack.Size() > 0 {
		v, _ := intStack.Pop()
		fmt.Println("popped:", v)
	}

	var strStack Stack[string]
	strStack.Push("hello")
	strStack.Push("world")
	for strStack.Size() > 0 {
		v, _ := strStack.Pop()
		fmt.Println("popped:", v)
	}
}
