// +build OMIT

package main

import "fmt"

func main() {
	s := []int{2, 3, 5, 7, 11, 13}
	printSlice(s)

	// Slice the slice to give it zero length.
	s = s[:0]
	printSlice(s)

	// Extend its length.
	s = s[:4]
	printSlice(s)

	// Drop its first two values.
	s = s[2:] // lenght of slice is default value so upper would be 4 (instead of 6)...
	printSlice(s)

	// Get final 3 elements, having drop 2 and skip the first one
	s = s[1:4]
	printSlice(s)

	// Get 4 elements...
	//s = s[:4]
	//printSlice(s)
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
