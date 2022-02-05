package main

import "fmt"

func main() {
	lines := []int{1, 2, 3}
	for i, v := range lines {
		fmt.Println(i, v)
	}
}
