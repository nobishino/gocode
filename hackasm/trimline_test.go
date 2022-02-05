package hackasm_test

import (
	"testing"

	"github.com/nobishino/gocode/hackasm"
)

func TestExportTrimLine(t *testing.T) {
	testcases := []struct {
		in   string
		want string
	}{
		{" @R0 ", "@R0"},
		{"@R0  ", "@R0"},
		{"   @R0             \r", "@R0"},
	}
	for _, tt := range testcases {
		tt := tt
		t.Run(tt.in+"->"+tt.want, func(t *testing.T) {
			got := hackasm.ExportTrimLine(tt.in)
			if got != tt.want {
				t.Errorf("want %s, got %s", tt.want, tt.in)
				for i, r := range []rune(tt.in) {
					t.Logf("%dth rune = %v", i, r)
				}
			}
		})
	}

}
