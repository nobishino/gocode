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
(function_Test_f)
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
