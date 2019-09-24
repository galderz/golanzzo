// +build OMIT

package main

import "fmt"

type Vertex struct {
	X int
	Y int
}

func main() {
	v := Vertex{1, 2}
	v.X = 4
	pX := &v.X
	pY := &v.Y
	fmt.Println(v.X)
	fmt.Println(pX)
	fmt.Println(pY)
}
