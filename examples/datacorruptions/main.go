package main

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

func main() {
	// fmt.Println(structCorrption())
	// fmt.Println(interfaceCorruption())
	// Corruption2()
	// stringCorruption()
	for {
		storeBuffer()
		// messagePassing()
		// IRIW()
		// loadBuffer()
		// coherence()
	}
}

// Goの仕様(メモリーモデル)によるとpanicする可能性があり実際ときどきpanicする
func storeBuffer() {
	var eg errgroup.Group
	// 共有変数
	x, y := 0, 0
	r1WasZero, r2WasZero := false, false
	eg.Go(func() error {
		x = 1
		r1 := y
		r1WasZero = r1 == 0
		return nil
	})
	eg.Go(func() error {
		y = 1
		r2 := x
		r2WasZero = r2 == 0
		return nil
	})
	eg.Wait() // エラー処理略
	if r1WasZero && r2WasZero {
		panic("Store Buffer Test Failed")
	}
}

// メモリーモデル上はpanicする可能性があるがpanicするのを確認できない
func messagePassing() {
	var eg errgroup.Group
	// 共有変数
	x, y := 0, 0
	eg.Go(func() error {
		if y == 1 && x == 0 {
			panic("Message Passing Test Fails")
		}
		return nil
	})
	eg.Go(func() error {
		x = 1
		y = 1
		return nil
	})
	eg.Wait() // エラー処理略
}

// メモリーモデル上はpanicする可能性があるがpanicするのを確認できない
func IRIW() {
	var eg errgroup.Group
	// 共有変数
	x, y := 0, 0
	// 読み取り結果の記録用
	var r1r2WasOneZero, r3r4WasOneZero bool
	eg.Go(func() error {
		x = 1
		return nil
	})
	eg.Go(func() error {
		y = 1
		return nil
	})
	eg.Go(func() error {
		r1 := x
		r2 := y
		r1r2WasOneZero = r1 == 1 && r2 == 0
		return nil
	})
	eg.Go(func() error {
		r3 := y
		r4 := x
		r3r4WasOneZero = r3 == 1 && r4 == 0
		return nil
	})
	eg.Wait() // エラー処理略
	if r1r2WasOneZero && r3r4WasOneZero {
		panic("IRIW Test Failed")
	}
}

// Goの仕様(メモリーモデル)によるとpanicする可能性があるがpanic確認できてない
func loadBuffer() {
	var eg errgroup.Group
	// 共有変数
	x, y := 0, 0
	var r1WasOne, r2WasOne bool
	eg.Go(func() error {
		r1 := x
		y = 1
		r1WasOne = r1 == 1
		return nil
	})
	eg.Go(func() error {
		r2 := y
		x = 1
		r2WasOne = r2 == 1
		return nil
	})
	eg.Wait() // エラー処理略
	if r1WasOne && r2WasOne {
		panic("Load Buffer Test Failed")
	}
}

// Goの仕様(メモリーモデル)によるとpanicする可能性があるがpanic確認できてない
func coherence() {
	var eg errgroup.Group
	// 共有変数
	x := 0
	// 読み取り結果の記録用
	var r1r2WasOneTwo, r3r4WasTwoOne bool
	eg.Go(func() error {
		x = 1
		return nil
	})
	eg.Go(func() error {
		x = 2
		return nil
	})
	eg.Go(func() error {
		r1 := x
		r2 := x
		r1r2WasOneTwo = r1 == 1 && r2 == 2
		return nil
	})
	eg.Go(func() error {
		r3 := x
		r4 := x
		r3r4WasTwoOne = r3 == 2 && r4 == 1
		return nil
	})
	eg.Wait() // エラー処理略
	if r1r2WasOneTwo && r3r4WasTwoOne {
		panic("Coherence Test Failed")
	}
}

func stringCorruption() string {
	var s string
	// writer
	go func() {
		arr := [2]string{"", "hello"}
		for i := 0; ; i++ {
			s = arr[i%2]
		}
	}()
	// reader
	for {
		fmt.Println(s)
	}
}

func structCorrption() string {
	type Pair struct {
		X int
		Y int
	}
	arr := []Pair{{X: 0, Y: 0}, {X: 1, Y: 1}}
	var p Pair
	go func() {
		for i := 0; ; i++ {
			p = arr[i%2]
		}
	}()
	for {
		read := p
		switch read.X + read.Y {
		case 0, 2:
		default:
			return fmt.Sprintf("struct corruption detected: %+v", read)
		}
	}
}

func interfaceCorruption() string {
	var x any

	go func() { // writer
		arr := []any{1, "hello"}
		for i := 0; ; i++ {
			x = arr[i%2]
		}
	}()
	// reader
	for {
		read := x
		switch r := read.(type) {
		case int:
			if r != 1 {
				return fmt.Sprintf("unexpected int value: %d", r)
			}
		case string:
			if len(r) != 5 {
				return fmt.Sprintf("unexpected string length :%d", len(r))
			}
		case nil:
		default:
			return fmt.Sprintf("strange type detected: %+v", read)
		}
	}
}

func sliceCorruption() {
	underlying := [5]int{1, 2, 3, 4, 5}
	var s []int

	go func() { // writer
		for i := 0; ; i++ {
			// rは1から5までの整数
			r := i%5 + 1
			// len == capであるようなスライスを新たに作り、
			// sに代入する
			s = underlying[:r:r]
		}
	}()
	// reader
	for {
		// len(s) == cap(s)は常に成り立つと期待する？
		if len(s) != cap(s) {
			panic(fmt.Sprintf("len(s) == %d and cap(s) == %d", len(s), cap(s)))
		}
	}
}

func mapCorruption() {
	// 共有変数
	m := map[int]int{}

	// writer
	go func() {
		for i := 0; ; i++ {
			m[i] = i
		}
	}()
	// reader
	for {
		if m[len(m)] > 10000 {
			break
		}
	}
}
func mapCorruption2() {
	// 共有変数
	m := map[int]int{}

	// writer
	go func() {
		for i := 0; ; i++ {
			m[i] = i
		}
	}()
	// reader
	for {
		// len(m)にだけアクセスする
		// 要素にはアクセスしない
		if len(m) > 10000 {
			break
		}
	}
}
