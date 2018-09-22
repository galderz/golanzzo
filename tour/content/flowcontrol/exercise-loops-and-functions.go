// +build no-build OMIT

package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := float64(1)

	for i := 0; i < 20; i++ {
		prev := z
		z -= (z*z - x) / (2*z)
		diff := math.Abs(prev - z)
		//fmt.Println("Prev: ", prev)
		//fmt.Println("Try: ", z)
		//fmt.Println("diff: ", diff)
		//fmt.Println()

		//if diff <= 0.000000000000001 {
		if diff <= 1e-15 {
			return z
		}
	}

	return z
}

func main() {
	fmt.Println(math.Sqrt(2))
	fmt.Println(Sqrt(2))
}
