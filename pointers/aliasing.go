package main

import "fmt"

// --- Aliasing and Unintended Mutation ---
//
// Two pointers alias when they point to the same memory. A write through either
// pointer changes what the other sees — which is surprising when the aliasing is
// accidental.
//
// The most common source of accidental aliasing in Go:
//
//   Slices — a slice header (pointer, length, capacity) is a value, but the
//   backing array is shared. Appending within capacity or writing by index
//   through a slice copy modifies the original's backing array.
//
//   Maps — always a reference type; copying a map variable copies the pointer,
//   not the data.
//
// Defensive copy pattern: when a function or struct should own its data
// exclusively, copy slices/maps on the way in (or out) rather than storing the
// caller's slice header directly.
//
// When aliasing IS intentional — pointer parameters, output parameters, shared
// buffers — document it clearly so callers know the ownership contract.

// store illustrates accidental slice aliasing via append-within-capacity.
type store struct {
	items []string
}

// addBad stores the slice header directly — the caller still shares the backing array.
func newStoreBad(items []string) *store {
	return &store{items: items}
}

// addGood copies the slice, so mutations to the original cannot reach s.items.
func newStoreGood(items []string) *store {
	cp := make([]string, len(items))
	copy(cp, items)
	return &store{items: cp}
}

func (s *store) set(i int, v string) { s.items[i] = v }
func (s *store) get(i int) string    { return s.items[i] }

// mapAliasing shows that assigning a map copies only the header pointer.
func mapAliasing() {
	original := map[string]int{"a": 1}
	alias := original // NOT a copy — both variables point to the same map
	alias["a"] = 99
	fmt.Printf("original[\"a\"] after alias write: %d (want 99, same map)\n", original["a"])

	// Defensive copy
	safeCopy := make(map[string]int, len(original))
	for k, v := range original {
		safeCopy[k] = v
	}
	safeCopy["a"] = 0
	fmt.Printf("original[\"a\"] after safeCopy write: %d (want 99, independent)\n", original["a"])
}

func DemoAliasing() {
	fmt.Println("=== Aliasing and Unintended Mutation ===")

	src := []string{"hello", "world"}

	bad := newStoreBad(src)
	src[0] = "mutated" // caller modifies src after handing it off
	fmt.Printf("bad store[0]: %q (want \"hello\", got caller's mutation)\n", bad.get(0))

	src2 := []string{"hello", "world"}
	good := newStoreGood(src2)
	src2[0] = "mutated"
	fmt.Printf("good store[0]: %q (want \"hello\", copy is independent)\n", good.get(0))

	mapAliasing()
}
