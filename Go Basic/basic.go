package main

import (
	"fmt"
)

// TODO: use of range in go
func main() {
	var num = []int{2, 3, 4}
	for i, n := range num {
		fmt.Printf("index: %v, value: %v\n", i, n)
	}
}

//GET request  :-here data is directly visible in thier URL
//POST request :-here request data is hidden.
