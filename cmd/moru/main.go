package main

import "C"

//export sum
func sum(a, b C.int) C.int {
	return a + b
}

func main() {}
