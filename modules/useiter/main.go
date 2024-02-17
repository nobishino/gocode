package main

func seq(n int) func(yield func(int) bool) {
	return func(yield func(int) bool) {
		for i := range n {
			if !yield(i) {
				break
			}
		}
	}
}

func main() {
	for i := range seq(10) {
		println(i)
	}
}
