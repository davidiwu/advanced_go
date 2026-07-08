package main

import "fmt"

func SumInts(m map[string]int64) int64 {
	var s int64
	for _, v := range m {
		s += v
	}
	return s
}

func SumFloats(f map[string]float64) float64 {
	var s float64

	for _, v := range f {
		s += v
	}

	return s
}

func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V

	for _, v := range m {
		s += v
	}
	return s
}

func main() {
	// original examples
	ints := map[string]int64{
		"first":  44,
		"second": 78,
	}
	floats := map[string]float64{
		"first":  44.75,
		"second": 78.23,
	}
	fmt.Println("=== Sum Functions ===")
	fmt.Println(SumInts(ints), SumFloats(floats))
	fmt.Println(SumIntsOrFloats(ints), SumIntsOrFloats(floats))
	fmt.Println()

	DemoStack()
	fmt.Println()
	DemoSliceUtils()
	fmt.Println()
	DemoOrdered()
	fmt.Println()
	DemoOption()
	fmt.Println()
	DemoCache()
}
