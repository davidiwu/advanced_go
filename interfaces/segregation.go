package main

import (
	"fmt"
	"strings"
)

// --- Interface Segregation ---
//
// Prefer many small, focused interfaces over one large "god interface."
// Callers declare only the methods they actually use, so:
//   - Implementations stay easy to swap (a mock only needs the small surface)
//   - Code is easier to test (no unused methods to stub out)
//   - Adding methods to the concrete type doesn't break existing callers
//
// This is the Go analogue of the Interface Segregation Principle (ISP).

// Three narrow interfaces instead of one big DataStore interface.
type Reader interface {
	Read(key string) (string, bool)
}

type Writer interface {
	Write(key, value string)
}

type Deleter interface {
	Delete(key string)
}

// ReadWriter is composed only where both behaviors are needed.
type ReadWriter interface {
	Reader
	Writer
}

// MemStore is the one concrete implementation. It satisfies all three
// individual interfaces and the composed ReadWriter.
type MemStore struct {
	data map[string]string
}

func NewMemStore() *MemStore {
	return &MemStore{data: make(map[string]string)}
}

func (m *MemStore) Read(key string) (string, bool) {
	v, ok := m.data[key]
	return v, ok
}
func (m *MemStore) Write(key, value string) { m.data[key] = value }
func (m *MemStore) Delete(key string)        { delete(m.data, key) }

// cacheLayer only needs to read — accepting Reader keeps it decoupled
// from write and delete operations it has no business calling.
func cacheLayer(r Reader, keys []string) {
	fmt.Println("  reading from cache layer:")
	for _, k := range keys {
		if v, ok := r.Read(k); ok {
			fmt.Printf("    %s = %q\n", k, v)
		} else {
			fmt.Printf("    %s = <not found>\n", k)
		}
	}
}

// populate only needs to write — a read-only store, or a write-only
// sink, both satisfy this parameter type.
func populate(w Writer, entries map[string]string) {
	for k, v := range entries {
		w.Write(k, v)
	}
}

// audit needs both read and delete — it accepts the composed interface.
func audit(rw ReadWriter, d Deleter, flagged []string) {
	fmt.Println("  auditing and removing flagged keys:")
	for _, k := range flagged {
		if v, ok := rw.Read(k); ok {
			fmt.Printf("    removing %s=%q\n", k, v)
			d.Delete(k)
		}
	}
}

// DemoSegregation shows that the single MemStore can be passed as any
// of the narrow interfaces, and each function only sees what it needs.
func DemoSegregation() {
	fmt.Println("=== Interface Segregation ===")

	store := NewMemStore()

	populate(store, map[string]string{
		"user:1": "alice",
		"user:2": "bob",
		"spam:1": strings.Repeat("x", 10),
	})

	cacheLayer(store, []string{"user:1", "user:2", "user:3"})

	audit(store, store, []string{"spam:1"})

	fmt.Println("  after audit:")
	cacheLayer(store, []string{"user:1", "spam:1"})
}
