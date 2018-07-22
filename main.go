package main

import (
	"bufio"
	"fmt"
	. "github.com/konart/gos/lib"

	"os"
)

func main() {
	worker := NewWorker(5)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		worker.Run(scanner.Text())
	}
	worker.Wait()

	fmt.Printf("Total: %d", worker.GetSum())
}