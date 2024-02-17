// GOEXPERIMENT=rangefunc
// (*)実行するとgo vetがエラーを出しますが、実行はできています
package main

func main() {
	for i := range seq() {
		println(i)
	}
}

func seq() func(yield func(int) bool) {
	return func(yield func(int) bool) {
		for i := range 10 {
			if !yield(i) {
				break
			}
		}
	}
}
