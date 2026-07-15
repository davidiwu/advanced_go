package main

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
