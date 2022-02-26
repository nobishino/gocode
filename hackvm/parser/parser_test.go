package parser_test

import (
	"strings"
	"testing"

	"github.com/nobishino/gocode/hackvm/parser"
)

func TestParser_Arithmetic(t *testing.T) {
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
		{"// begin code\nadd // comment", []string{"add"}},
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

func TestParser_MemoryAccess(t *testing.T) {
	type args struct {
		arg1 string
		arg2 int
	}
	testcases := []struct {
		in      string
		cmdType string
		arg     []args
	}{
		{"push constant 7", "C_PUSH", []args{{"constant", 7}}},
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
			if gotType != tc.cmdType {
				t.Errorf("want %q, but got %q", tc.cmdType, gotType)
			}
			gotArg1 := p.Arg1()
			wantArg1 := tc.arg[idx].arg1
			if gotArg1 != wantArg1 {
				t.Errorf("want %q, but got %q", wantArg1, gotArg1)
			}
			gotArg2 := p.Arg2()
			wantArg2 := tc.arg[idx].arg2
			if gotArg2 != wantArg2 {
				t.Errorf("want %q, but got %q", wantArg2, gotArg2)
			}
		}
		if p.HasMoreCommands() {
			t.Errorf("there should not be more than %d commands", len(tc.arg))
		}
	}
}
