package main

import "fmt"

func main() {
	fmt.Println(structCorrption())
}

func structCorrption() string {
	type Pair struct {
		X int
		Y int
	}
	arr := []Pair{{X: 0, Y: 0}, {X: 1, Y: 1}}
	var p Pair
	done := make(chan struct{})
	go func() {
		for i := 0; ; i++ {
			select {
			case <-done:
				return
			default:
				p = arr[i%2]
			}
		}
	}()
	for {
		read := p
		switch read.X + read.Y {
		case 0, 2:
		default:
			close(done)
			return fmt.Sprintf("struct corruption detected: %+v", read)
		}
	}
}
