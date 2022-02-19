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
	}
	for _, c := range testcases {
		var buf bytes.Buffer // 書き込み先(ファイルの代わりだけどテスト用にBufferを使う)
		writer := codewriter.New(&buf)

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
