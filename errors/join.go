package main

import (
	"errors"
	"fmt"
)

// --- Multi-Error Aggregation Pattern ---
//
// Often you want to report all failures rather than stopping at the first —
// batch validation, parallel operations, etc.
//
// errors.Join (Go 1.20) wraps multiple errors into one. It returns nil when
// all inputs are nil. The result implements Unwrap() []error so errors.Is/As
// traverse every member automatically.

var (
	ErrEmptyName = errors.New("name is empty")
	ErrBadEmail  = errors.New("email is invalid")
	ErrShortPass = errors.New("password too short")
)

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
	return errors.Join(errs...)
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
	return errors.Join(errs...)
}

func DemoJoin() {
	fmt.Println("=== Multi-Error Aggregation ===")

	// Case 1: all fields valid — errors.Join returns nil.
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
}
