package parser

import "io"

type Parser struct {
}

func New(r io.Reader) *Parser {
	return new(Parser)
}

// HasMoreCommandsは、入力において更にコマンドが存在するかを返す
func (p *Parser) HasMoreCommands() bool {
	return false
}

// 入力から次のコマンドを読み、それを現コマンドとする。
// HasMoreCommandsがtrueを返したときだけAdvanceを呼ぶべき。
// 最初は現コマンドは空になる
func (p *Parser) Advance() {

}

// 現コマンドの種類を返す。算術コマンドはすべてC_ARITHMETICが返される
func (p *Parser) CommandType() string {
	return ""
}

// 現コマンドの最初の引数を返す。　C_ARITHMETICの場合、コマンド自体(add,subなど)が返される。
// 現コマンドがC_RETURNの場合、本メソッドは呼ばないようにする
func (p *Parser) Arg1() string {
	return ""
}

// 現コマンドの2番目の引数を返す。現コマンドがC_PUSH,C_POP,C_FUNCTION,C_CALLの場合のみ
// 本メソッドを呼ぶようにする。
func (p *Parser) Arg2() int {
	return 0
}
