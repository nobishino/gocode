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
		{"add\nsub\nneg", []string{"add", "sub", "neg"}},
		{"neg\neq\ngt\nlt\nand\nor\nnot\n", []string{"neg", "eq", "gt", "lt", "and", "or", "not"}},
	}
	for _, tc := range testcases {
		p := parser.New(strings.NewReader(tc.in))
		for idx := 0; idx < len(tc.arg); idx++ {
			if !p.HasMoreCommands() {
				t.Fatalf("there should be %d commands, but got only %d commands",
					len(tc.arg), idx)
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
		if p.HasMoreCommands() {
			t.Errorf("there should not be more than %d commands", len(tc.arg))
		}
	}
}
