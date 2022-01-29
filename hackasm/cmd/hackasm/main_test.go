package main

import (
	"path/filepath"
	"testing"
)

func TestExec(t *testing.T) {
	srcPath := "./testdata/PongL.asm"
	tmp := t.TempDir()
	outPath := filepath.Join(tmp, "out.hack")
	exec(srcPath, outPath)
}
