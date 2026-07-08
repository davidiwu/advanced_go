package main

import "fmt"

// --- Option Type Pattern (Maybe / Optional) ---
//
// Optional[T] explicitly represents a value that may or may not be present.
// It is a cleaner alternative to:
//   - pointer semantics (*T) where nil means "absent" — but pointers also
//     mean "heap-allocated" which muddies the intent
//   - (T, bool) return pairs — correct but verbose at every call site
//   - sentinel zero values (e.g. -1 for "no age") — fragile and implicit
//
// The zero value of Optional[T] is naturally None — safe without init.

// Optional wraps a value of any type with an explicit presence flag.
type Optional[T any] struct {
	value T    // the wrapped value; only meaningful when present is true
	present bool
}

// Some wraps a value that is present.
func Some[T any](v T) Optional[T] {
	return Optional[T]{value: v, present: true}
}

// None returns an absent optional of type T.
func None[T any]() Optional[T] {
	return Optional[T]{} // zero value: present == false
}

// Unwrap returns the value and whether it is present.
// Callers are forced to handle both cases, unlike a plain pointer dereference.
func (o Optional[T]) Unwrap() (T, bool) {
	return o.value, o.present
}

// UnwrapOr returns the value if present, or the provided fallback otherwise.
// Useful when a default makes sense and the absence case needs no special logic.
func (o Optional[T]) UnwrapOr(fallback T) T {
	if o.present {
		return o.value
	}
	return fallback
}

// String makes Optional print readably in fmt calls.
func (o Optional[T]) String() string {
	if o.present {
		return fmt.Sprintf("Some(%v)", o.value)
	}
	return "None"
}

// --- Example usage ---

// User represents a record that may have an optional nickname and age.
type User struct {
	Name     string
	Nickname Optional[string]
	Age      Optional[int] // 0 is a valid age, so zero value is ambiguous without Optional
}

// findUser simulates a lookup that may return nothing.
func findUser(name string) Optional[User] {
	db := map[string]User{
		"alice": {
			Name:     "alice",
			Nickname: Some("ali"),
			Age:      Some(30),
		},
		"bob": {
			Name:     "bob",
			Nickname: None[string](), // bob has no nickname
			Age:      Some(0),        // age 0 is present and valid — not confused with absence
		},
	}
	if u, ok := db[name]; ok {
		return Some(u)
	}
	return None[User]()
}

// printUser demonstrates both Unwrap and UnwrapOr at the call site.
func printUser(name string) {
	result := findUser(name)
	user, ok := result.Unwrap()
	if !ok {
		fmt.Printf("%s: not found\n", name)
		return
	}
	// UnwrapOr gives a clean default without an explicit if block.
	nick := user.Nickname.UnwrapOr("(none)")
	age, _ := user.Age.Unwrap()
	fmt.Printf("name=%-6s  nickname=%-8s  age=%d\n", user.Name, nick, age)
}

// DemoOptionType exercises the Optional type across found, not-found,
// and zero-value-present cases.
func DemoOptionType() {
	fmt.Println("=== Option Type ===")
	printUser("alice")
	printUser("bob")
	printUser("carol") // not in the db
}
