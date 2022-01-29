package hackasm_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nobishino/gocode/hackasm"
)

func TestAssemble(t *testing.T) {
	testcase := []struct {
		name string
	}{
		{"PongL"},
	}
	for _, tt := range testcase {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			src := open(t, filepath.Join("testdata", tt.name+".asm"))
			dest := new(bytes.Buffer)

			if err := hackasm.Assemble(src, dest); err != nil {
				t.Fatal(err)
			}

			expect := open(t, filepath.Join("testdata", tt.name+".hack"))

			equal(t, expect, dest)
		})
	}
}

func open(t *testing.T, path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})
	return f
}

func equal(t *testing.T, want, got io.Reader) {
	wantB := new(strings.Builder)
	gotB := new(strings.Builder)
	if _, err := io.Copy(wantB, want); err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(gotB, got); err != nil {
		t.Fatal(err)
	}
	if wantB.String() != gotB.String() {
		t.Errorf("does not match. want:%s, got:%s", wantB.String(), gotB.String())
	}
}
