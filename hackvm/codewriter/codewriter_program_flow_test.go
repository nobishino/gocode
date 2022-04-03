package codewriter_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nobishino/gocode/hackvm/codewriter"
)

func TestWriteLabel(t *testing.T) {
	const wantErr = true
	testcases := []struct {
		label    string
		fileName string
		want     string
		err      bool
	}{
		{
			label:    "xyz",
			fileName: "Test1.vm",
			want: `// label global$xyz
(global$xyz)
`,
		},
		{
			label:    "x_.:123",
			fileName: "Test.vm",
			want: `// label global$x_.:123
(global$x_.:123)
`,
		},
		{
			label:    "",
			fileName: "Test2.vm",
			err:      wantErr,
		},
		{
			label:    "3ABC",
			fileName: "Test2.vm",
			err:      wantErr,
		},
	}
	for _, c := range testcases {
		t.Run(c.label, func(t *testing.T) {
			var buf bytes.Buffer
			writer := codewriter.New(&buf)

			writer.SetFileName(c.fileName)
			err := writer.WriteLabel(c.label)
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
func TestWriteGoto(t *testing.T) {
	testcases := []struct {
		label    string
		fileName string
		want     string
		err      bool
	}{
		{
			label:    "xyz",
			fileName: "Test1.vm",
			want: `// goto global$xyz
@global$xyz
0;JMP
`,
		},
	}
	for _, c := range testcases {
		t.Run(c.label, func(t *testing.T) {
			var buf bytes.Buffer
			writer := codewriter.New(&buf)

			writer.SetFileName(c.fileName)
			err := writer.WriteGoto(c.label)
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

func TestWriteIf(t *testing.T) {
	testcases := []struct {
		label    string
		fileName string
		want     string
		err      bool
	}{
		{
			label:    "xyz",
			fileName: "Test1.vm",
			want: `// if-goto global$xyz
@SP
M=M-1
A=M
D=M
@global$xyz
D;JNE
`,
		},
	}
	for _, c := range testcases {
		t.Run(c.label, func(t *testing.T) {
			var buf bytes.Buffer
			writer := codewriter.New(&buf)

			writer.SetFileName(c.fileName)
			err := writer.WriteIf(c.label)
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
