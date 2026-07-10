package main

import "fmt"

// --- Slice / Map Preallocation Pattern ---
//
// When you know (or can estimate) the final size of a slice or map upfront,
// preallocating with make avoids the repeated doubling reallocations that
// happen when appending to a nil or zero-capacity collection.
//
// Cost model:
//   growing nil slice to N elements  → O(log N) reallocations, O(N) total copies
//   preallocated slice of cap N       → 0 reallocations, 0 copies

// doubleWithGrowth builds a slice by appending to a nil slice.
// Each time the backing array is full, Go doubles its capacity and copies.
func doubleWithGrowth(src []int) []int {
	var result []int // nil — zero capacity
	for _, v := range src {
		result = append(result, v*2)
	}
	return result
}

// doublePreallocated builds the same slice but preallocates the exact capacity.
// No reallocation ever occurs.
func doublePreallocated(src []int) []int {
	result := make([]int, 0, len(src)) // cap = len(src) → no growth needed
	for _, v := range src {
		result = append(result, v*2)
	}
	return result
}

// countWords counts word frequency. Growing map vs preallocated map.
func countWordsGrowth(words []string) map[string]int {
	counts := map[string]int{} // starts at default small capacity
	for _, w := range words {
		counts[w]++
	}
	return counts
}

func countWordsPreallocated(words []string) map[string]int {
	counts := make(map[string]int, len(words)) // hint avoids early rehashing
	for _, w := range words {
		counts[w]++
	}
	return counts
}

func DemoPrealloc() {
	fmt.Println("=== Slice / Map Preallocation ===")

	nums := make([]int, 100)
	for i := range nums {
		nums[i] = i + 1
	}

	grown := doubleWithGrowth(nums)
	prealloc := doublePreallocated(nums)
	fmt.Printf("growth result len=%d cap=%d\n", len(grown), cap(grown))
	fmt.Printf("prealloc result len=%d cap=%d  (cap == len: no wasted space)\n",
		len(prealloc), cap(prealloc))

	// Demonstrate capacity growth steps for a nil slice.
	var s []int
	fmt.Println("\ncapacity growth steps (nil start, appending 1..16):")
	prevCap := 0
	for i := range 16 {
		s = append(s, i)
		if cap(s) != prevCap {
			fmt.Printf("  len=%-3d cap=%d  ← reallocation\n", len(s), cap(s))
			prevCap = cap(s)
		}
	}

	words := []string{"go", "is", "fast", "go", "is", "simple", "go"}
	gc := countWordsGrowth(words)
	pc := countWordsPreallocated(words)
	fmt.Printf("\nword counts match: %v  (growth=%v, prealloc=%v)\n",
		fmt.Sprint(gc) == fmt.Sprint(pc), gc, pc)
}
