package main

import (
	"fmt"
	"sync"
)

// --- sync.Once ---
//
// sync.Once guarantees that a function is executed exactly once,
// no matter how many goroutines call Do concurrently.
// This is the idiomatic way to implement lazy initialization in Go —
// for example, opening a database connection or loading config only
// when it is first needed.
//
// The zero value of sync.Once is ready to use — no constructor needed.

type dbConn struct {
	dsn string
}

var (
	instance *dbConn
	once     sync.Once
)

// getDB returns the singleton connection, creating it on the first call.
// Safe to call from any number of goroutines simultaneously.
func getDB() *dbConn {
	once.Do(func() {
		// Expensive initialization runs exactly once.
		fmt.Println("  [once] opening database connection…")
		instance = &dbConn{dsn: "postgres://localhost/mydb"}
	})
	return instance
}

func DemoOnce() {
	fmt.Println("=== sync.Once ===")

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			db := getDB()
			fmt.Printf("  goroutine %d got db @ %p\n", id, db)
		}(i)
	}
	wg.Wait()
	// The "[once] opening…" line appears exactly once regardless of concurrency.
}
