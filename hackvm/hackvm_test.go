package hackvm_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nobishino/gocode/hackvm"
)

func TestTranslate(t *testing.T) {
	testcases := []struct {
		in   string
		want string
	}{
		{in: "SimpleAdd.vm", want: "SimpleAdd.asm"},
		{in: "StackTest.vm", want: "StackTest.asm"},
	}
	for _, tc := range testcases {
		r := openFile(t, tc.in)
		out := new(strings.Builder)
		hackvm.Translate(out, r, tc.in)

		want := new(strings.Builder)
		if _, err := io.Copy(want, openFile(t, tc.want)); err != nil {
			t.Fatal(err)
		}
		if out.String() != want.String() {
			t.Errorf("want:\n%s\nbut got:\n%s\n", want.String(), out.String())
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
