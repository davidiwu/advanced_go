package main

import (
	"fmt"
	"unsafe"
)

// --- Struct Field Alignment Pattern ---
//
// The CPU reads memory in aligned words (typically 8 bytes on 64-bit platforms).
// If a field is placed at an unaligned offset the compiler inserts invisible
// padding bytes before it. Reordering fields from largest to smallest
// eliminates padding and shrinks the struct's memory footprint.
//
// Rule of thumb: order fields largest-to-smallest alignment requirement,
// then group fields by access locality.

// Unaligned has poor field ordering — padding inserted by the compiler.
//
//	bool (1) + 7-byte pad + int64 (8) + int32 (4) + bool (1) + 3-byte pad = 24 bytes
type Unaligned struct {
	Flag1  bool
	Count  int64
	Score  int32
	Flag2  bool
}

// Aligned reorders the same fields to eliminate padding.
//
//	int64 (8) + int32 (4) + bool (1) + bool (1) + 2-byte pad = 16 bytes
type Aligned struct {
	Count  int64
	Score  int32
	Flag1  bool
	Flag2  bool
}

// Cache line packing: keep hot fields together so they fit in one 64-byte
// cache line rather than straddling two.
type HotCold struct {
	// Hot path: read on every request
	RequestCount int64
	BytesIn      int64
	BytesOut     int64
	// Cold path: read only during health checks
	StartTime    int64
	BuildVersion [32]byte
}

func DemoStructAlign() {
	fmt.Println("=== Struct Field Alignment ===")

	fmt.Printf("Unaligned size: %d bytes\n", unsafe.Sizeof(Unaligned{}))
	fmt.Printf("Aligned   size: %d bytes  (%.0f%% smaller)\n",
		unsafe.Sizeof(Aligned{}),
		float64(unsafe.Sizeof(Unaligned{})-unsafe.Sizeof(Aligned{}))/
			float64(unsafe.Sizeof(Unaligned{}))*100)

	fmt.Println()
	fmt.Println("field offsets in Unaligned:")
	u := Unaligned{}
	fmt.Printf("  Flag1  offset=%d\n", unsafe.Offsetof(u.Flag1))
	fmt.Printf("  Count  offset=%d  (7 bytes of padding before this)\n", unsafe.Offsetof(u.Count))
	fmt.Printf("  Score  offset=%d\n", unsafe.Offsetof(u.Score))
	fmt.Printf("  Flag2  offset=%d\n", unsafe.Offsetof(u.Flag2))

	fmt.Println()
	fmt.Println("field offsets in Aligned:")
	a := Aligned{}
	fmt.Printf("  Count  offset=%d\n", unsafe.Offsetof(a.Count))
	fmt.Printf("  Score  offset=%d\n", unsafe.Offsetof(a.Score))
	fmt.Printf("  Flag1  offset=%d\n", unsafe.Offsetof(a.Flag1))
	fmt.Printf("  Flag2  offset=%d  (no padding between booleans)\n", unsafe.Offsetof(a.Flag2))

	fmt.Println()
	fmt.Printf("HotCold size: %d bytes\n", unsafe.Sizeof(HotCold{}))
	hc := HotCold{}
	fmt.Printf("  hot fields span bytes 0..%d (fit in one 64-byte cache line: %v)\n",
		unsafe.Offsetof(hc.BytesOut)+unsafe.Sizeof(hc.BytesOut)-1,
		unsafe.Offsetof(hc.BytesOut)+unsafe.Sizeof(hc.BytesOut) <= 64)
}
