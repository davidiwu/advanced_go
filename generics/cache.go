package main

import "fmt"

// Cache uses two type parameters: K for the key and V for the value.
// K is constrained to "comparable" (not "any") because Go map keys must be
// comparable — using "any" here would be a compile error since not all types
// support == and !=. V has no such restriction so it stays "any".
// Without generics this would require interface{} for both, losing type safety
// on every Get call and requiring type assertions by the caller.
type Cache[K comparable, V any] struct {
	m map[K]V
}

// NewCache uses a constructor function so the internal map is always initialized.
// T parameters on the function match the type parameters declared on the struct —
// the caller specifies the concrete types when calling NewCache[string, string]().
func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{m: make(map[K]V)}
}

// Methods below reference K and V from the type definition — they are in type scope.
// No constraints are repeated; the compiler already knows K is comparable and V is any.
func (c *Cache[K, V]) Set(k K, v V)      { c.m[k] = v }
func (c *Cache[K, V]) Get(k K) (V, bool) { v, ok := c.m[k]; return v, ok }
func (c *Cache[K, V]) Delete(k K)        { delete(c.m, k) }
func (c *Cache[K, V]) Len() int          { return len(c.m) }

func DemoCache() {
	fmt.Println("=== Generic Cache ===")

	strCache := NewCache[string, string]()
	strCache.Set("lang", "Go")
	strCache.Set("version", "1.19")

	if v, ok := strCache.Get("lang"); ok {
		fmt.Println("lang:", v)
	}

	intCache := NewCache[int, float64]()
	intCache.Set(1, 3.14)
	intCache.Set(2, 2.72)

	if v, ok := intCache.Get(1); ok {
		fmt.Println("key 1:", v)
	}
	if _, ok := intCache.Get(99); !ok {
		fmt.Println("key 99: not found")
	}

	fmt.Println("cache size:", intCache.Len())
}
