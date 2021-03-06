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
			if err := p.Advance(); err != nil {
				t.Fatal(err)
			}
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
		cmdType []string
		arg     []args
	}{
		{"push constant 7", []string{"C_PUSH"}, []args{{"constant", 7}}},
		{"push constant 8", []string{"C_PUSH"}, []args{{"constant", 8}}},
		{"push constant 7\npush constant 8", []string{"C_PUSH", "C_PUSH"}, []args{{"constant", 7}, {"constant", 8}}},
		{"pop static 5", []string{"C_POP"}, []args{{"static", 5}}},
		{"push local 7\npop argument 1", []string{"C_PUSH", "C_POP"}, []args{{"local", 7}, {"argument", 1}}},
	}
	for _, tc := range testcases {
		p := parser.New(strings.NewReader(tc.in))
		for idx := 0; idx < len(tc.arg); idx++ {
			if !p.HasMoreCommands() {
				t.Fatalf("there should be %d commands, but got only %d commands",
					len(tc.arg), idx)
			}
			if err := p.Advance(); err != nil {
				t.Fatal(err)
			}
			gotType := p.CommandType()
			wantCmdType := tc.cmdType[idx]
			if gotType != wantCmdType {
				t.Errorf("want %q, but got %q", wantCmdType, gotType)
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

func TestParser_ProgramFlow(t *testing.T) {
	type args struct {
		arg1 string
		arg2 int
	}
	testcases := []struct {
		in      string
		cmdType []string
		arg     []args
	}{
		{"label xxx", []string{"C_LABEL"}, []args{{arg1: "xxx"}}},
		{"goto yyy", []string{"C_GOTO"}, []args{{arg1: "yyy"}}},
		{"if-goto zzz", []string{"C_IF"}, []args{{arg1: "zzz"}}},
	}
	for _, tc := range testcases {
		p := parser.New(strings.NewReader(tc.in))
		for idx := 0; idx < len(tc.arg); idx++ {
			if !p.HasMoreCommands() {
				t.Fatalf("there should be %d commands, but got only %d commands",
					len(tc.arg), idx)
			}
			if err := p.Advance(); err != nil {
				t.Fatal(err)
			}
			gotType := p.CommandType()
			wantCmdType := tc.cmdType[idx]
			if gotType != wantCmdType {
				t.Errorf("want %q, but got %q", wantCmdType, gotType)
			}
			gotArg1 := p.Arg1()
			wantArg1 := tc.arg[idx].arg1
			if gotArg1 != wantArg1 {
				t.Errorf("want %q, but got %q", wantArg1, gotArg1)
			}
		}
		if p.HasMoreCommands() {
			t.Errorf("there should not be more than %d commands", len(tc.arg))
		}
	}
}

func TestParser_Function(t *testing.T) {
	type args struct {
		arg1 string
		arg2 int
	}
	testcases := []struct {
		in      string
		cmdType []string
		arg     []args
	}{
		{"function f 0", []string{"C_FUNCTION"}, []args{{arg1: "f", arg2: 0}}},
		{"function g 1", []string{"C_FUNCTION"}, []args{{arg1: "g", arg2: 1}}},
		{"function myFunc 2", []string{"C_FUNCTION"}, []args{{arg1: "myFunc", arg2: 2}}},
		{"call f 0", []string{"C_CALL"}, []args{{arg1: "f", arg2: 0}}},
		{"call g 1", []string{"C_CALL"}, []args{{arg1: "g", arg2: 1}}},
		{"call myFunc 2", []string{"C_CALL"}, []args{{arg1: "myFunc", arg2: 2}}},
		{"return", []string{"C_RETURN"}, []args{{}}},
	}
	for _, tc := range testcases {
		p := parser.New(strings.NewReader(tc.in))
		for idx := 0; idx < len(tc.arg); idx++ {
			if !p.HasMoreCommands() {
				t.Fatalf("there should be %d commands, but got only %d commands",
					len(tc.arg), idx)
			}
			if err := p.Advance(); err != nil {
				t.Fatal(err)
			}
			gotType := p.CommandType()
			wantCmdType := tc.cmdType[idx]
			if gotType != wantCmdType {
				t.Errorf("want %q, but got %q", wantCmdType, gotType)
			}

			if gotType == "C_RETURN" {
				continue
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
