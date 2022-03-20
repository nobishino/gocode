package codewriter

import (
	"fmt"
	"io"
	"strings"
)

type CodeWriter struct {
	out             io.Writer
	comparisonIndex int
	fileName        string
}

//出力ファイル/ストリームを開き書き込む準備を行う
func New(w io.Writer) *CodeWriter {
	return &CodeWriter{
		out: w,
	}
}

// CodeWriterに、あたらしいVMファイルの変換が開始したことを伝える
// .vm拡張子以外を渡すとpanicする
func (c *CodeWriter) SetFileName(name string) {
	if !strings.HasSuffix(name, ".vm") {
		panic(fmt.Sprintf("file name must have extension '.vm'. got: %q", name))
	}
	c.fileName = name[:len(name)-3]
}

// 出力ファイルを閉じる
func (c *CodeWriter) Close() error { return nil }
