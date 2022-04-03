package codewriter

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

type CodeWriter struct {
	out             io.Writer
	comparisonIndex int
	fileName        string
	currentFuncName string
	callCount       int // how many times "call f n" command are written
}

//出力ファイル/ストリームを開き書き込む準備を行う
func New(w io.Writer) *CodeWriter {
	return &CodeWriter{
		// WriteLabel, WriteGoto, WriteIfの単体テストで用いるための仮の関数名
		// 実際のVMコードではなんらかのfunctionの内部にしかbranching命令が書かれない
		currentFuncName: "global",
		out:             w,
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

func (c *CodeWriter) Init() error {
	initCode := `// init
@256
D=A
@SP
M=D
`
	if _, err := fmt.Fprint(c.out, initCode); err != nil {
		return errors.WithStack(err)
	}
	c.SetFileName("Sys.vm")
	if err := c.WriteCall("Sys.init", 0); err != nil {
		return err
	}
	terminateCode := `// termination
(TERMINATE)
@TERMINATE
0;JMP
`
	if _, err := fmt.Fprint(c.out, terminateCode); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
