package main

import "fmt"

func main() {
	var m map[string]int
	m = make(map[string]int)
	m["route"] = 66
	i, ok := m["route"]
	fmt.Printf("Contains? %t\n", ok)
	fmt.Printf("Element: %s\n", i)
}
