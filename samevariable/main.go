package main

import (
	"log"
	"sync"
)

func main() {
	exec()
}

var x int

func exec() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		x = 1
	}()
	go func() {
		defer wg.Done()
		x1 := x
		x2 := x
		log.Println(x1, x2)
		if x1 == 1 && x2 == 0 {
			panic("HELLO!")
		}
	}()
	wg.Wait()
}
