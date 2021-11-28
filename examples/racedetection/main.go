package main

import "fmt"

func main() {
	var i int
	go func() {
		i++
	}()
	i++
	fmt.Println(i)
}

// from https://go.dev/ref/mem
// modified a little bit
// printing 2, 0 is possible
func incorrectSync() {
	var a, b int

	f := func() {
		a = 1
		b = 2
	}

	g := func() {
		print(b)
		print(a)
	}

	go f()
	g()
}
