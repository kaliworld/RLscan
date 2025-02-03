package main

import "fmt"

func main() {
	bytes := []byte{'H', 'e', 'l', 'l', 'o'}
	for _, byteValue := range bytes {
		fmt.Printf("%d\n", byteValue)
	}
}
