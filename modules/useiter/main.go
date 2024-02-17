// GOEXPERIMENT=rangefunc
// (*)実行するとgo vetがエラーを出しますが、実行はできています
package main

import "github.com/nobishino/gocoro/iter"

func main() {
	for i := range seq() {
		println(i)
	}
}

func seq() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range 10 {
			if !yield(i) {
				break
			}
		}
	}
}
