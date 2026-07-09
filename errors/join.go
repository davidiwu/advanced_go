package main

import (
	"errors"
	"fmt"
	"strings"
)

// --- Multi-Error Aggregation Pattern ---
//
// Often you want to report all failures rather than stopping at the first —
// batch validation, parallel operations, etc.
//
// Go 1.20 added errors.Join for this. On earlier toolchains the same
// behaviour is achieved with a small multiError type that implements
// Unwrap() []error, which errors.Is / errors.As also traverse.
//
// This file shows the manual pattern (compatible with Go 1.19) and
// notes where errors.Join would replace joinErrs.

var (
	ErrEmptyName = errors.New("name is empty")
	ErrBadEmail  = errors.New("email is invalid")
	ErrShortPass = errors.New("password too short")
)

// multiError holds several errors and presents them as one.
// Unwrap() []error is the Go 1.20 multi-unwrap protocol; errors.Is/As
// understand it even in the stdlib shipped with 1.20+.
type multiError struct {
	errs []error
}

func (m *multiError) Error() string {
	msgs := make([]string, len(m.errs))
	for i, e := range m.errs {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "; ")
}

// Unwrap returns the slice — the Go 1.20 multi-unwrap protocol.
// errors.Is/As in Go 1.20+ traverse it automatically.
func (m *multiError) Unwrap() []error { return m.errs }

// Is provides Go 1.19-compatible traversal: errors.Is calls this when it
// finds the method, allowing it to search all members.
func (m *multiError) Is(target error) bool {
	for _, e := range m.errs {
		if errors.Is(e, target) {
			return true
		}
	}
	return false
}

// joinErrs returns nil if all errors are nil, otherwise a *multiError.
// In Go 1.20+ you would write: return errors.Join(errs...)
func joinErrs(errs ...error) error {
	var nonNil []error
	for _, e := range errs {
		if e != nil {
			nonNil = append(nonNil, e)
		}
	}
	if len(nonNil) == 0 {
		return nil
	}
	return &multiError{errs: nonNil}
}

// validateUser collects all validation failures instead of stopping at the first.
func validateUser(name, email, password string) error {
	var errs []error
	if name == "" {
		errs = append(errs, ErrEmptyName)
	}
	if len(email) < 5 || email[len(email)-1] == '@' {
		errs = append(errs, ErrBadEmail)
	}
	if len(password) < 8 {
		errs = append(errs, ErrShortPass)
	}
	return joinErrs(errs...)
}

// processItems runs an operation on each item and collects all failures.
func processItems(items []string) error {
	var errs []error
	for _, item := range items {
		if item == "" {
			errs = append(errs, fmt.Errorf("empty item"))
		} else if item == "bad" {
			errs = append(errs, fmt.Errorf("invalid item %q", item))
		}
	}
	return joinErrs(errs...)
}

func DemoJoin() {
	fmt.Println("=== Multi-Error Aggregation ===")

	// Case 1: all fields valid — joinErrs returns nil.
	err := validateUser("alice", "alice@example.com", "s3cr3t!!")
	fmt.Printf("valid user:   err=%v\n\n", err)

	// Case 2: multiple failures — all are reported in one error.
	err = validateUser("", "bad@", "short")
	fmt.Printf("invalid user: %v\n", err)

	// errors.Is searches all members via Unwrap() []error.
	fmt.Printf("  Is(ErrEmptyName):  %v\n", errors.Is(err, ErrEmptyName))
	fmt.Printf("  Is(ErrBadEmail):   %v\n", errors.Is(err, ErrBadEmail))
	fmt.Printf("  Is(ErrShortPass):  %v\n\n", errors.Is(err, ErrShortPass))

	// Case 3: batch processing collects all item errors.
	err = processItems([]string{"ok", "bad", "", "also-ok"})
	fmt.Printf("batch errors: %v\n", err)
	fmt.Printf("  non-nil: %v\n\n", err != nil)

	fmt.Println("note: Go 1.20 adds errors.Join which implements this same")
	fmt.Println("      Unwrap()[]error protocol without a custom type.")
}
