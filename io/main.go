package main

import "fmt"

func main() {
	// Custom io.Reader: read from a non-standard source
	DemoCustomReader()
	fmt.Println()

	// Custom io.Writer: capture or transform written bytes
	DemoCustomWriter()
	fmt.Println()

	// Decorator pattern: wrap a Reader or Writer to add behavior
	DemoDecorator()
	fmt.Println()

	// io.Pipe: connect a writer to a reader with no intermediate buffer
	DemoPipe()
	fmt.Println()

	// bufio.Scanner: line-by-line reading with custom split functions
	DemoScanner()
}
