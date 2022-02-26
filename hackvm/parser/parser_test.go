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
		arg string
	}{
		{in: "add\n", arg: "add"},
		{in: "sub\n", arg: "sub"},
	}
	for _, tc := range testcases {
		p := parser.New(strings.NewReader(tc.in))
		if !p.HasMoreCommands() {
			t.Fatal("parser should have at least 1 command")
		}
		p.Advance()
		gotType := p.CommandType()
		if gotType != wantType {
			t.Errorf("want %q, but got %q", wantType, gotType)
		}
		gotArg1 := p.Arg1()
		if gotArg1 != tc.arg {
			t.Errorf("want %q, but got %q", tc.arg, gotArg1)
		}
	}
}
