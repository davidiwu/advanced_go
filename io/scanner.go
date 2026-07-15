package main

// --- bufio.Scanner ---
//
// bufio.Scanner wraps any io.Reader and tokenizes it according to a split
// function. It handles buffering internally, so callers iterate one token
// at a time without loading the entire stream into memory.
//
// Pattern:
//
//	scanner := bufio.NewScanner(r)
//	for scanner.Scan() {
//	    process(scanner.Text()) // or scanner.Bytes() for a []byte view
//	}
//	if err := scanner.Err(); err != nil { ... }
//
// Built-in split functions: ScanLines (default), ScanWords, ScanBytes, ScanRunes.
//
// Custom split functions follow the signature:
//
//	func(data []byte, atEOF bool) (advance int, token []byte, err error)
//
//   - advance: how many bytes to consume from the front of data.
//   - token:   the next token to return (nil means "no token yet").
//   - err:     non-nil stops scanning with that error.
//   - return (0, nil, nil) to request more data from the reader.
//
// When to use: line-by-line log tailing, CSV/TSV parsing, or any stream where
// you want to iterate over tokens without reading everything into memory first.
//
// Key point: always check scanner.Err() after the loop. Scan() returns false
// for both EOF (no error) and a read error — they look identical without the check.
// For tokens larger than 64 KB, call scanner.Buffer(buf, maxSize) before the first Scan.

import (
	"bufio"
	"fmt"
	"strings"
	"unicode"
)

func DemoScanner() {
	fmt.Println("=== bufio.Scanner ===")

	// Default split: ScanLines — most common use case.
	input := "first line\nsecond line\nthird line"
	scanner := bufio.NewScanner(strings.NewReader(input))
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		fmt.Printf("line %d: %q\n", lineNum, scanner.Text())
	}

	// ScanWords: split on whitespace boundaries.
	fmt.Println()
	words := bufio.NewScanner(strings.NewReader("  the quick  brown fox  "))
	words.Split(bufio.ScanWords)
	var collected []string
	for words.Scan() {
		collected = append(collected, words.Text())
	}
	fmt.Printf("words: %v\n", collected)

	// Custom split function: split on comma, trimming surrounding whitespace.
	// A SplitFunc signature: func(data []byte, atEOF bool) (advance int, token []byte, err error)
	csvSplit := func(data []byte, atEOF bool) (int, []byte, error) {
		// If we're at EOF with no data, signal done.
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		for i, b := range data {
			if b == ',' {
				token := []byte(strings.TrimFunc(string(data[:i]), unicode.IsSpace))
				return i + 1, token, nil
			}
		}
		// No comma found; if at EOF return the remainder as the last token.
		if atEOF {
			token := []byte(strings.TrimFunc(string(data), unicode.IsSpace))
			return len(data), token, nil
		}
		// Request more data.
		return 0, nil, nil
	}

	fmt.Println()
	csv := bufio.NewScanner(strings.NewReader("alpha, beta , gamma,delta"))
	csv.Split(csvSplit)
	for csv.Scan() {
		fmt.Printf("field: %q\n", csv.Text())
	}

	// MaxScanTokenSize: the default buffer is 64 KB. For larger tokens,
	// call scanner.Buffer(buf, maxSize) before the first Scan call.
	fmt.Println()
	const bigLine = 70000
	big := strings.Repeat("x", bigLine) + "\nnormal line"
	s := bufio.NewScanner(strings.NewReader(big))
	s.Buffer(make([]byte, bigLine+1), bigLine+1)
	for s.Scan() {
		fmt.Printf("big scanner: %d bytes\n", len(s.Bytes()))
	}
	if err := s.Err(); err != nil {
		fmt.Printf("scanner error: %v\n", err)
	}
}
