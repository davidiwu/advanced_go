package main

import "fmt"

// ========================================================================
// PATTERN: Embedding vs. explicit composition
//
// Embedding promotes methods and makes the outer type satisfy interfaces
// automatically. Explicit composition (a named field) keeps the boundary
// visible and avoids method-set surprises.
//
// Choose embedding when:
//   - The outer type IS-A specialisation of the inner type (conceptually).
//   - You want the outer type to satisfy the inner type's interfaces.
//   - You want clean callers: s.Method() instead of s.inner.Method().
//
// Choose explicit composition when:
//   - The relationship is HAS-A, not IS-A.
//   - Accidental interface satisfaction would be confusing or dangerous.
//   - You need multiple fields of the same type (only one can be embedded).
//   - You want to make the delegation explicit in the call site.
// ========================================================================

// --- Shared behaviour type ---

type Auditor struct{}

func (a Auditor) RecordAction(action string) {
	fmt.Printf("  [audit] %s\n", action)
}

// --- Embedding approach (IS-A / interface promotion) ---

// AuditedStore embeds Auditor. Callers can call store.RecordAction directly,
// and AuditedStore satisfies any interface Auditor satisfies.
type AuditedStore struct {
	Auditor
	data map[string]string
}

func NewAuditedStore() *AuditedStore {
	return &AuditedStore{data: make(map[string]string)}
}

func (s *AuditedStore) Set(k, v string) {
	s.RecordAction(fmt.Sprintf("set %q = %q", k, v)) // promoted method
	s.data[k] = v
}

// --- Explicit composition approach (HAS-A) ---

// ReportingStore has an Auditor as a named field. The delegation is explicit
// and RecordAction is NOT promoted to the outer type's method set.
type ReportingStore struct {
	auditor Auditor // named field — no promotion
	data    map[string]string
}

func NewReportingStore() *ReportingStore {
	return &ReportingStore{data: make(map[string]string)}
}

func (s *ReportingStore) Set(k, v string) {
	s.auditor.RecordAction(fmt.Sprintf("set %q = %q", k, v)) // explicit delegation
	s.data[k] = v
}

func DemoComposition() {
	fmt.Println("=== Embedding vs. Explicit Composition ===")

	as := NewAuditedStore()
	as.Set("env", "production")
	// Promoted: callers can also call as.RecordAction() directly.
	as.RecordAction("manual audit entry")

	fmt.Println()

	rs := NewReportingStore()
	rs.Set("env", "staging")
	// Not promoted: rs.RecordAction() would not compile.
	// The auditor field is encapsulated; callers must go through Set.
	fmt.Println("  (RecordAction is not accessible on ReportingStore)")
}
