package fermat

import (
	"fmt"
	"testing"
)

// x^n + y^n = z^n
// s.t x, y, z > 0, n >= 3

// corpus = (n, x, y, z)
func FuzzTestFermatFinalTheorem(f *testing.F) {
	f.Add(uint(3), uint(4), uint(5), uint(6))
	f.Fuzz(func(t *testing.T, n, x, y, z uint) {
		skipIfNotPositive(t, n, x, y, z)
		if n < 3 {
			t.Skip()
		}
		left := pow(x, n) + pow(y, n)
		right := pow(z, n)
		if left == right {
			t.Errorf("It's a counterexample for Fermat's final theorem!")
			// Of course it's not, because of overflow...
		}
	})
}

func skipIfNotPositive(t *testing.T, xs ...uint) {
	for _, x := range xs {
		if x == 0 {
			t.Skip()
		}
	}
}

func TestPow(t *testing.T) {
	t.Parallel()
	cases := []struct {
		x    uint
		y    uint
		want uint
	}{
		{1, 1, 1},
		{1, 2, 1},
		{2, 2, 4},
		{2, 3, 8},
		{5, 3, 125},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("%d^%d=%d", tc.x, tc.y, tc.want), func(t *testing.T) {
			t.Parallel()
			got := pow(tc.x, tc.y)
			if got != tc.want {
				t.Errorf("want %d, but got %d", tc.want, got)
			}
		})
	}
}
