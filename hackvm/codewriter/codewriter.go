package codewriter

import (
	"fmt"
	"io"
)

type CodeWriter struct {
	out io.Writer
}

//出力ファイル/ストリームを開き書き込む準備を行う
func New(w io.Writer) *CodeWriter {
	return &CodeWriter{
		out: w,
	}
}

// CodeWriterに、あたらしいVMファイルの変換が開始したことを伝える
func (c *CodeWriter) SetFileName(n string) {

}

//与えられた算術コマンドをアッセンブリーコードに変換し、それを書き込む
func (c *CodeWriter) WriteArithmetic(command string) error {
	code := `// add
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
`
	_, err := io.WriteString(c.out, code)
	if err != nil {
		return err
	}
	return nil
}

//C_PUSHまたはC_POPコマンドをアッセンブリーに変換しそれを書き込む
func (c *CodeWriter) WritePushPop(command string, segment string, index int) error {
	format := `// push constant %[1]d
@%[1]d
D=A
@SP
A=M
M=D
@SP
M=M+1
`
	code := fmt.Sprintf(format, index)
	_, err := io.WriteString(c.out, code)
	if err != nil {
		return err
	}

	return nil
}

// 出力ファイルを閉じる
func (c *CodeWriter) Close() error { return nil }
