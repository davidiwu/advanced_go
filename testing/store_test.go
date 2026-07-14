package testing_patterns

import (
	"errors"
	"testing"
)

// --- Interface-Based Test Doubles ---
//
// A test double is a lightweight stand-in for a real dependency.
// In Go the idiomatic approach is: define the dependency as an interface,
// then implement a minimal fake that satisfies it in the test file.
//
// No third-party mocking framework needed — the interface IS the contract,
// and a struct with the right methods IS the mock.

// fakeStore is a test double for Storer backed by an in-memory map.
// It also records which IDs were saved so tests can assert on side-effects.
type fakeStore struct {
	data    map[string]int
	savelog []string // tracks every Save call in order
}

func newFakeStore() *fakeStore {
	return &fakeStore{data: make(map[string]int)}
}

func (f *fakeStore) Save(id string, val int) error {
	f.data[id] = val
	f.savelog = append(f.savelog, id)
	return nil
}

func (f *fakeStore) Load(id string) (int, error) {
	v, ok := f.data[id]
	if !ok {
		return 0, errors.New("not found")
	}
	return v, nil
}

func TestService_Increment(t *testing.T) {
	t.Run("increments existing value", func(t *testing.T) {
		store := newFakeStore()
		store.data["x"] = 10

		svc := NewService(store)
		got, err := svc.Increment("x", 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 15 {
			t.Errorf("got %d, want 15", got)
		}
		// Assert the side-effect: Save was called with the right ID.
		if len(store.savelog) != 1 || store.savelog[0] != "x" {
			t.Errorf("savelog = %v, want [x]", store.savelog)
		}
	})

	t.Run("returns error when key missing", func(t *testing.T) {
		store := newFakeStore() // empty — Load will return "not found"
		svc := NewService(store)
		_, err := svc.Increment("missing", 1)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// errorStore is a test double that always returns an error from Save.
// Having separate, purpose-built fakes keeps each scenario crystal clear.
type errorStore struct{ fakeStore }

func (e *errorStore) Save(_ string, _ int) error {
	return errors.New("disk full")
}

func TestService_Increment_SaveError(t *testing.T) {
	store := &errorStore{}
	store.data = map[string]int{"y": 1}

	svc := NewService(store)
	_, err := svc.Increment("y", 1)
	if err == nil {
		t.Fatal("expected error from Save, got nil")
	}
}
