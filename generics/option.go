package main

import "fmt"

// Option wraps a value that may or may not be present, eliminating the need to
// return (T, error) or use nil pointers for "not found" cases. Without generics
// this would require a separate OptionInt, OptionString, etc., or an interface{}
// wrapper that loses type information at the call site.
type Option[T any] struct {
	value T
	ok    bool
}

func Some[T any](v T) Option[T] { return Option[T]{value: v, ok: true} }
func None[T any]() Option[T]    { return Option[T]{} }

func (o Option[T]) Unwrap() (T, bool) { return o.value, o.ok }

func (o Option[T]) UnwrapOr(fallback T) T {
	if o.ok {
		return o.value
	}
	return fallback
}

// findFirst is a generic function that returns an Option rather than (T, bool),
// making it clear at the call site that absence is a normal, expected outcome —
// not an error. T is declared here on the function (function scope), not borrowed
// from any type, because findFirst is a standalone function.
func findFirst[T any](s []T, fn func(T) bool) Option[T] {
	for _, v := range s {
		if fn(v) {
			return Some(v)
		}
	}
	return None[T]()
}

func DemoOption() {
	fmt.Println("=== Option Type ===")

	nums := []int{1, 3, 5, 8, 9}

	found := findFirst(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println("first even:", found.UnwrapOr(-1))

	notFound := findFirst(nums, func(n int) bool { return n > 100 })
	fmt.Println("first > 100:", notFound.UnwrapOr(-1))

	if v, ok := found.Unwrap(); ok {
		fmt.Println("unwrapped:", v)
	}
}
