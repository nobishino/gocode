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
			label:    "x_.:123",
			fileName: "Test.vm",
			want: `// label x_.:123
(label_Test_x_.:123)
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
