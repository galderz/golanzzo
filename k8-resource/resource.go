package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/resource"
)

func main() {
	cores := resource.NewMilliQuantity(500, resource.DecimalSI)
	fmt.Printf("cores = %v\n", cores)

	cores = resource.NewMilliQuantity(500 / 2, resource.DecimalSI)
	fmt.Printf("cores = %v\n", cores)

	var cpuLimit = 700

	cores = resource.NewMilliQuantity(int64(cpuLimit), resource.DecimalSI)
	fmt.Printf("cores = %v\n", cores)

	cores = resource.NewMilliQuantity(int64(cpuLimit / 2), resource.DecimalSI)
	fmt.Printf("cores = %v\n", cores)

	parsedCores := resource.MustParse("250m")
	halfParsedCores := resource.NewMilliQuantity(parsedCores.MilliValue() / 2, resource.DecimalSI)
	fmt.Printf("cores = %v\n", parsedCores)
	fmt.Printf("half cores = %v\n", halfParsedCores)

	parsedCores = resource.MustParse("0.25")
	halfParsedCores = resource.NewMilliQuantity(parsedCores.MilliValue() / 2, resource.DecimalSI)
	fmt.Printf("cores = %v\n", parsedCores)
	fmt.Printf("half cores = %v\n", halfParsedCores)
}
