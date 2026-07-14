package main

import (
	"fmt"
	"sync"
)

// --- sync.RWMutex ---
//
// An RWMutex is a reader/writer mutual exclusion lock. It allows many
// goroutines to hold a read lock simultaneously, but only one goroutine
// can hold the write lock — and only when no readers hold the lock.
//
// Use RWMutex instead of Mutex when:
//   - Reads significantly outnumber writes
//   - Read operations are slow enough that the overhead of serializing them matters
//
// If writes are frequent, a plain Mutex is usually simpler and just as fast.

type cache struct {
	mu    sync.RWMutex
	items map[string]string
}

func newCache() *cache {
	return &cache{items: make(map[string]string)}
}

// Set acquires the exclusive write lock to modify the map.
func (c *cache) Set(key, val string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = val
}

// Get acquires only a read lock — multiple goroutines can call Get concurrently.
func (c *cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.items[key]
	return v, ok
}

func DemoRWMutex() {
	fmt.Println("=== sync.RWMutex ===")

	c := newCache()
	var wg sync.WaitGroup

	// One writer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.Set("lang", "Go")
		c.Set("version", "1.22")
	}()
	wg.Wait()

	// Many concurrent readers — all hold RLock at the same time
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if v, ok := c.Get("lang"); ok {
				fmt.Printf("reader %d: lang=%s\n", id, v)
			}
		}(i)
	}
	wg.Wait()
}
