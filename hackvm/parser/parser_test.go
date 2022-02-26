package parser_test

import (
	"os"
	"path/filepath"
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
		cmdType []string
		arg     []args
	}{
		{"push constant 7", []string{"C_PUSH"}, []args{{"constant", 7}}},
		{"push constant 8", []string{"C_PUSH"}, []args{{"constant", 8}}},
		{"push constant 7\npush constant 8", []string{"C_PUSH", "C_PUSH"}, []args{{"constant", 7}, {"constant", 8}}},
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

func TestParser_Real(t *testing.T) {
	type args struct {
		arg1 string
		arg2 int
	}
	testcases := []struct {
		in      string
		cmdType []string
		arg     []args
	}{
		{
			"SimpleAdd.vm",
			[]string{"C_PUSH", "C_PUSH", "C_ARITHMETIC"},
			[]args{{arg1: "constant", arg2: 7}, {arg1: "constant", arg2: 8}, {arg1: "add"}},
		},
	}
	for _, tc := range testcases {
		p := parser.New(openFile(t, tc.in))
		for idx := 0; idx < len(tc.arg); idx++ {
			if !p.HasMoreCommands() {
				t.Fatalf("there should be %d commands, but got only %d commands",
					len(tc.arg), idx)
			}
			p.Advance()
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
			if wantCmdType == "C_ARITHMETIC" {
				continue
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

func openFile(t *testing.T, filePath string) *os.File {
	f, err := os.Open(filepath.Join("testdata", filePath))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})
	return f
}
