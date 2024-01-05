package corruptslice_test

import (
	"testing"
)

// メモリーモデル関係のテスト
func TestAppendConsistency_ExpectFailOccationally(t *testing.T) {
	t.Skip()
	for i := 1; ; i++ {
		s := []int{1}
		go func() {
			s = append(s, 2)
		}()
		for len(s) < 2 {
		}
		a := (*[2]int)(s)

		if a[1] != 2 {
			t.Fatalf("expect a[1] == 2 but got %d, at %dth trial", a[1], i)
		}
	}
}
