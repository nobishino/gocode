package codewriter_test

import (
	"bytes"
	"testing"

	"github.com/nobishino/gocode/hackvm/codewriter"
)

func TestWritePushPop(t *testing.T) {
	testcases := []struct {
		command string
		segment string
		index   int
		want    string
	}{
		{
			command: "C_PUSH", segment: "constant", index: 7,
			want: `// push constant 7
@7
D=A
@SP
A=M
M=D
@SP
M=M+1
`,
		},
		{
			command: "C_POP", segment: "static", index: 5,
			want: `// pop static 5
@SP
M=M-1
A=M
D=M
@filename.5
M=D
`,
		},
		{
			command: "C_PUSH", segment: "static", index: 7,
			want: `// push static 7
@filename.7
D=M
@SP
A=M
M=D
@SP
M=M+1
`,
		},
		{
			command: "C_POP", segment: "local", index: 1,
			want: `// pop local 1
@LCL
D=M
@1
D=D+A
@R13 // general register
M=D
@SP
M=M-1
A=M
D=M
@R13 // general register
A=M
M=D
`,
		},
		{
			command: "C_PUSH", segment: "local", index: 3,
			want: `// push local 3
@3
D=A
@LCL
D=D+M
A=D
D=M
@SP
A=M
M=D
@SP
M=M+1
`,
		},
		{
			command: "C_POP", segment: "argument", index: 4,
			want: `// pop argument 4
@ARG
D=M
@4
D=D+A
@R13 // general register
M=D
@SP
M=M-1
A=M
D=M
@R13 // general register
A=M
M=D
`,
		},
		{
			command: "C_PUSH", segment: "this", index: 9,
			want: `// push this 9
@9
D=A
@THIS
D=D+M
A=D
D=M
@SP
A=M
M=D
@SP
M=M+1
`,
		},
		{
			command: "C_PUSH", segment: "that", index: 7,
			want: `// push that 7
@7
D=A
@THAT
D=D+M
A=D
D=M
@SP
A=M
M=D
@SP
M=M+1
`,
		},
		{
			command: "C_POP", segment: "temp", index: 6,
			want: `// pop temp 6
@SP
M=M-1
A=M
D=M
@R11
M=D
`,
		},
		{
			command: "C_PUSH", segment: "temp", index: 7,
			want: `// push temp 7
@R12
D=M
@SP
A=M
M=D
@SP
M=M+1
`,
		},
		{
			command: "C_POP", segment: "pointer", index: 0,
			want: `// pop pointer 0
@SP
M=M-1
A=M
D=M
@R3
M=D
`,
		},
		{
			command: "C_POP", segment: "pointer", index: 1,
			want: `// pop pointer 1
@SP
M=M-1
A=M
D=M
@R4
M=D
`,
		},
		{
			command: "C_PUSH", segment: "pointer", index: 0,
			want: `// push pointer 0
@R3
D=M
@SP
A=M
M=D
@SP
M=M+1
`,
		},
		{
			command: "C_PUSH", segment: "pointer", index: 1,
			want: `// push pointer 1
@R4
D=M
@SP
A=M
M=D
@SP
M=M+1
`,
		},
	}
	for _, c := range testcases {
		var buf bytes.Buffer // 書き込み先(ファイルの代わりだけどテスト用にBufferを使う)
		writer := codewriter.New(&buf)
		writer.SetFileName("filename.vm")

		err := writer.WritePushPop(c.command, c.segment, c.index)
		if err != nil {
			t.Fatal(err)
		}

		got := buf.String() // 結果を取り出す
		if got != c.want {
			t.Errorf("want: \n%s\nbut got:\n%s\n", c.want, got)
		}
	}
}
