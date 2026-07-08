package main

import "fmt"

// Filter, Map, and Reduce are classic functional utilities. Without generics you
// would need to write a separate version for every element type (FilterInts,
// FilterStrings, ...) or accept []interface{} and cast on every access.
// [T any] lets a single function work on any slice type while remaining type-safe.

func Filter[T any](s []T, fn func(T) bool) []T {
	var result []T
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map introduces a second type parameter U, because the input and output element
// types can differ (e.g. []string in, []int out). Both T and U are declared on
// this function — neither is shared with another type, so both have function scope.
func Map[T, U any](s []T, fn func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// Reduce also uses two type params: T for the slice element type and U for the
// accumulator type. This lets you reduce a []string down to an int, for example.
func Reduce[T, U any](s []T, init U, fn func(U, T) U) U {
	acc := init
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

func DemoSliceUtils() {
	fmt.Println("=== Slice Utilities (Filter / Map / Reduce) ===")

	nums := []int{1, 2, 3, 4, 5, 6}

	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println("evens:", evens)

	doubled := Map(nums, func(n int) int { return n * 2 })
	fmt.Println("doubled:", doubled)

	words := []string{"go", "is", "fun"}
	lengths := Map(words, func(s string) int { return len(s) })
	fmt.Println("word lengths:", lengths)

	sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
	fmt.Println("sum:", sum)
}
