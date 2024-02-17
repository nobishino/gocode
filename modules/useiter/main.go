package main

func seq() func(yield func(int) bool) {
	return func(yield func(int) bool) {
		for i := range 10 {
			if !yield(i) {
				break
			}
		}
	}
}

func main() {
	for i := range seq() {
		println(i)
	}
}
