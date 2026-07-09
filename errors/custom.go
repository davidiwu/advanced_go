package main

import (
	"errors"
	"fmt"
)

// --- Custom Error Types with Is / As / Unwrap Pattern ---
//
// A struct error type carries structured context (status code, field name,
// request ID, etc.) that callers can extract programmatically.
//
// Three optional methods extend the default behaviour:
//   Unwrap() error        — makes the type transparent to errors.Is/As chains
//   Is(target error) bool — custom equality: match by code, not by identity
//   As(target any) bool   — custom extraction (rare; errors.As handles most cases)

// APIError is a structured error carrying an HTTP-style status code and message.
type APIError struct {
	Code    int
	Message string
	Cause   error // wrapped inner error
}

func (e *APIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("API %d %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("API %d %s", e.Code, e.Message)
}

// Unwrap lets errors.Is/As traverse into Cause.
func (e *APIError) Unwrap() error { return e.Cause }

// Is enables matching by code: errors.Is(err, &APIError{Code: 404}) returns
// true for any *APIError with Code==404, regardless of message or cause.
func (e *APIError) Is(target error) bool {
	var t *APIError
	if errors.As(target, &t) {
		return e.Code == t.Code
	}
	return false
}

// Sentinel targets used with errors.Is code-based matching.
var (
	ErrAPI404 = &APIError{Code: 404}
	ErrAPI500 = &APIError{Code: 500}
)

func fetchResource(id int) error {
	if id == 0 {
		return &APIError{Code: 404, Message: "resource not found"}
	}
	if id < 0 {
		root := fmt.Errorf("connection refused")
		return &APIError{Code: 500, Message: "upstream failure", Cause: root}
	}
	return nil
}

func DemoCustom() {
	fmt.Println("=== Custom Error Types (Is / As / Unwrap) ===")

	cases := []int{1, 0, -1}
	for _, id := range cases {
		err := fetchResource(id)
		if err == nil {
			fmt.Printf("  id=%-3d → ok\n", id)
			continue
		}

		// errors.As: extract the *APIError to read its fields.
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("  id=%-3d → APIError code=%d msg=%q\n", id, apiErr.Code, apiErr.Message)
		}

		// errors.Is with code-based matching via the custom Is() method.
		switch {
		case errors.Is(err, ErrAPI404):
			fmt.Printf("          matched ErrAPI404 (code-based Is)\n")
		case errors.Is(err, ErrAPI500):
			fmt.Printf("          matched ErrAPI500 (code-based Is)\n")
			// Unwrap chain: errors.Is keeps going into Cause.
			rootCause := errors.Unwrap(err)
			if rootCause != nil {
				fmt.Printf("          root cause: %v\n", rootCause)
			}
		}
	}
}
