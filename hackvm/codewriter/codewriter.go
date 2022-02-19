package codewriter

import "io"

type CodeWriter struct{}

//出力ファイル/ストリームを開き書き込む準備を行う
func New(w io.Writer) *CodeWriter {
	return new(CodeWriter)
}

// CodeWriterに、あたらしいVMファイルの変換が開始したことを伝える
func (c *CodeWriter) SetFileName(n string) {

}

//与えられた算術コマンドをアッセンブリーコードに変換し、それを書き込む
func (c *CodeWriter) WriteArithmetic(command string) error {
	return nil
}

//C_PUSHまたはC_POPコマンドをアッセンブリーに変換しそれを書き込む
func (c *CodeWriter) WritePushPop(command string, segment string, index int) error {
	return nil
}

// 出力ファイルを閉じる
func (c *CodeWriter) Close() error { return nil }
