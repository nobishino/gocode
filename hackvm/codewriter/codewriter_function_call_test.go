package codewriter_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nobishino/gocode/hackvm/codewriter"
)

func TestWriteFunctionDeclaration(t *testing.T) {
	testcases := []struct {
		funcName  string
		numLocals int
		fileName  string
		want      string
		err       bool
	}{
		{
			funcName:  "f",
			numLocals: 3,
			fileName:  "Test.vm",
			want: `// function f 3
(f)
@0
D=A
@SP
A=M
M=D
@SP
M=M+1
@SP
A=M
M=D
@SP
M=M+1
@SP
A=M
M=D
@SP
M=M+1
`,
		},
	}
	for _, c := range testcases {
		t.Run(c.funcName, func(t *testing.T) {
			var buf bytes.Buffer
			writer := codewriter.New(&buf)

			writer.SetFileName(c.fileName)
			err := writer.WriteFunction(c.funcName, c.numLocals)
			if err != nil && !c.err {
				t.Fatal(err)
			}
			if c.err {
				if err == nil {
					t.Fatal("want non-nil error but got nil")
				}
				return
			}

			got := buf.String() // 結果を取り出す
			if diff := cmp.Diff(c.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestWriteFunctionReturn(t *testing.T) {
	testcases := []struct {
		fileName string
		want     string
	}{
		{
			fileName: "Test.vm",
			want: `// return
// FRAME = LCL
@LCL
D=M
@R13
M=D
// RET = *(FRAME-5)
@5
A=D-A
D=M
@R14
M=D
// *ARG = pop()
@SP
M=M-1
@SP
A=M
D=M
@ARG
A=M
M=D
// SP = ARG+1
@ARG
D=M+1
@SP
M=D
// THAT = *(FRAME-1)
@R13
D=M
@1
A=D-A
D=M
@THAT
M=D
// THIS = *(FRAME-2)
@R13
D=M
@2
A=D-A
D=M
@THIS
M=D
// ARG = *(FRAME-3)
@R13
D=M
@3
A=D-A
D=M
@ARG
M=D
// LCL = *(FRAME-4)
@R13
D=M
@4
A=D-A
D=M
@LCL
M=D
// goto RET
@R14
A=M
0;JMP
`,
		},
	}
	for _, c := range testcases {
		t.Run("return", func(t *testing.T) {
			var buf bytes.Buffer
			writer := codewriter.New(&buf)

			writer.SetFileName(c.fileName)
			err := writer.WriteReturn()
			if err != nil {
				t.Fatal(err)
			}

			got := buf.String() // 結果を取り出す
			if diff := cmp.Diff(c.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestWriteFunctionCall(t *testing.T) {
	testcases := []struct {
		funcName string
		numArgs  int
		fileName string
		want     string
		err      bool
	}{
		{
			funcName: "f",
			numArgs:  1,
			fileName: "Test.vm",
			want: `// call f 1
// push return-address 
@return_address_0
D=A
@SP
A=M
M=D
@SP
M=M+1
// push LCL
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1
// push ARG
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
// push THIS
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
// push THAT
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1
// ARG = SP-n-5
@6
D=A
@SP
D=M-D
@ARG
M=D
// LCL = SP
@SP
D=M
@LCL
M=D
// goto f
@f
0;JMP
(return_address_0)
`,
		},
	}
	for _, c := range testcases {
		t.Run(c.funcName, func(t *testing.T) {
			var buf bytes.Buffer
			writer := codewriter.New(&buf)

			writer.SetFileName(c.fileName)
			err := writer.WriteCall(c.funcName, c.numArgs)
			if err != nil && !c.err {
				t.Fatal(err)
			}
			if c.err {
				if err == nil {
					t.Fatal("want non-nil error but got nil")
				}
				return
			}

			got := buf.String() // 結果を取り出す
			if diff := cmp.Diff(c.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
