package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	host := os.Getenv("HOST")
	addresses, err := net.LookupHost(host)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Resolved addresses: %+v\n", addresses)
}
