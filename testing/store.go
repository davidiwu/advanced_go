package testing_patterns

import "fmt"

// Storer is the interface our business logic depends on.
// Accepting an interface (not a concrete type) lets tests
// substitute a lightweight fake without a real database.
type Storer interface {
	Save(id string, val int) error
	Load(id string) (int, error)
}

// Service contains business logic that uses a Storer.
type Service struct {
	store Storer
}

func NewService(s Storer) *Service { return &Service{store: s} }

// Increment loads the current value for id, adds delta, and saves it back.
func (s *Service) Increment(id string, delta int) (int, error) {
	cur, err := s.store.Load(id)
	if err != nil {
		return 0, fmt.Errorf("load %q: %w", id, err)
	}
	next := cur + delta
	if err := s.store.Save(id, next); err != nil {
		return 0, fmt.Errorf("save %q: %w", id, err)
	}
	return next, nil
}
