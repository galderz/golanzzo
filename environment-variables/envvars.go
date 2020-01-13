package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	os.Setenv("FOO", "1")
	fmt.Println("FOO:", os.Getenv("FOO"))
	fmt.Println("BAR:", os.Getenv("BAR"))

	fmt.Println()
	printEnvironmentVariables()
	fmt.Println()

	os.Setenv("FOO", "2")
	fmt.Println("FOO:", os.Getenv("FOO"))
	printEnvironmentVariables()
}

func printEnvironmentVariables() {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 1)
		fmt.Println(pair[0])
	}
}
