package codewriter_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func TestWriteArithmetic(t *testing.T) {
	const wantErr = true
	testcases := []struct {
		command string
		want    string
		err     bool
	}{
		{
			command: "neg",
			want: `// neg
@SP
M=M-1
A=M
D=M
M=-D
@SP
M=M+1
`,
		},
		{
			command: "not",
			want: `// not
@SP
M=M-1
A=M
D=M
M=!D
@SP
M=M+1
`,
		},
		{
			command: "add",
			want: `// add
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=D+M
@SP
M=M+1
`,
		},
		{
			command: "sub",
			want: `// sub
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=M-D
@SP
M=M+1
`,
		},
		{
			command: "and",
			want: `// and
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=D&M
@SP
M=M+1
`,
		},
		{
			command: "or",
			want: `// or
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=D|M
@SP
M=M+1
`,
		},
		{
			command: "eq",
			want: `// eq
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=M-D
@eq_0_true
D;JEQ
(eq_0_false)
@SP
A=M
M=0
@SP
M=M+1
@eq_0_end
0;JMP
(eq_0_true)
@SP
A=M
M=-1
@SP
M=M+1
(eq_0_end)
`,
		},
		{
			command: "gt",
			want: `// gt
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=M-D
@gt_0_true
D;JGT
(gt_0_false)
@SP
A=M
M=0
@SP
M=M+1
@gt_0_end
0;JMP
(gt_0_true)
@SP
A=M
M=-1
@SP
M=M+1
(gt_0_end)
`,
		},
		{
			command: "lt",
			want: `// lt
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=M-D
@lt_0_true
D;JLT
(lt_0_false)
@SP
A=M
M=0
@SP
M=M+1
@lt_0_end
0;JMP
(lt_0_true)
@SP
A=M
M=-1
@SP
M=M+1
(lt_0_end)
`,
		},
		{command: "ADD", err: wantErr},
		{command: "", err: wantErr},
	}
	for _, c := range testcases {
		t.Run(c.command, func(t *testing.T) {

			var buf bytes.Buffer // 書き込み先(ファイルの代わりだけどテスト用にBufferを使う)
			writer := codewriter.New(&buf)

			err := writer.WriteArithmetic(c.command)
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
