package main

import (
	"fmt"
	"strings"
)

type Mat [2][2]float64

func (m Mat) String() string {
	var s strings.Builder
	for i := 0; i < len(m); i++ {
		fmt.Fprint(&s, m[i])
		if i < len(m)-1 {
			fmt.Fprintln(&s)
		}
	}
	return s.String()
}

func Prod(x, y Mat) Mat {
	result := Mat{}
	result[0][0] = x[0][0]*y[0][0] + x[0][1]*y[1][0]
	result[0][1] = x[0][0]*y[0][1] + x[0][1]*y[1][1]
	result[1][0] = x[1][0]*y[0][0] + x[1][1]*y[1][0]
	result[1][1] = x[1][0]*y[0][1] + x[1][1]*y[1][1]
	return result
}

var E = Mat{
	{1, 0},
	{0, 1},
}

func Power(x Mat, n int) Mat {
	result := E
	for i := 0; i < n; i++ {
		result = Prod(result, x)
	}
	return result
}

func main() {
	A := Mat{
		{-2, 1},
		{5, 2},
	}
	fmt.Println("A^2 =\n", Power(A, 2))
	fmt.Println("A^3 =\n", Power(A, 3))
	fmt.Println("A^4 =\n", Power(A, 4))
}
