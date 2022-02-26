package parser_test

import (
	"strings"
	"testing"

	"github.com/nobishino/gocode/hackvm/parser"
)

func TestParserArithmetic(t *testing.T) {
	const wantType = "C_ARITHMETIC"
	testcases := []struct {
		in  string
		arg []string
	}{
		{"add", []string{"add"}},
		{"add\n", []string{"add"}},
		{"sub", []string{"sub"}},
	}
	for _, tc := range testcases {
		p := parser.New(strings.NewReader(tc.in))
		var idx int
		if !p.HasMoreCommands() {
			if idx < len(tc.arg) {
				t.Fatalf("there should be %d commands, but got only %d commands",
					len(tc.arg), idx)
			}
			return // has read all commands
		}
		p.Advance()
		gotType := p.CommandType()
		if gotType != wantType {
			t.Errorf("want %q, but got %q", wantType, gotType)
		}
		gotArg1 := p.Arg1()
		if gotArg1 != tc.arg[idx] {
			t.Errorf("want %q, but got %q", tc.arg, gotArg1)
		}
	}
}
