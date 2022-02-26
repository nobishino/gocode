package parser

import (
	"io"
	"strings"
)

type Parser struct {
	cmds    [][]string
	current int
}

func New(r io.Reader) *Parser {
	sb := new(strings.Builder)
	io.Copy(sb, r) // TODO handle error
	lines := strings.Split(sb.String(), "\n")
	cmds := make([][]string, 0, len(lines))
	for _, line := range lines {
		if idx := strings.Index(line, "//"); idx != -1 { // コメントを無視する
			line = line[:idx]
		}
		line := strings.Trim(line, "\r\n\t ")
		if line == "" {
			continue // 空行は無視する
		}
		words := strings.Split(line, " ")
		cmds = append(cmds, words)
	}
	return &Parser{
		cmds:    cmds,
		current: -1,
	}
}

// HasMoreCommandsは、入力において更にコマンドが存在するかを返す
func (p *Parser) HasMoreCommands() bool {
	return p.current+1 < len(p.cmds)
}

// 入力から次のコマンドを読み、それを現コマンドとする。
// HasMoreCommandsがtrueを返したときだけAdvanceを呼ぶべき。
// 最初は現コマンドは空になる
func (p *Parser) Advance() {
	p.current++
}

// 現コマンドの種類を返す。算術コマンドはすべてC_ARITHMETICが返される
func (p *Parser) CommandType() string {
	return "C_ARITHMETIC"
}

// 現コマンドの最初の引数を返す。　C_ARITHMETICの場合、コマンド自体(add,subなど)が返される。
// 現コマンドがC_RETURNの場合、本メソッドは呼ばないようにする
func (p *Parser) Arg1() string {
	return p.cmds[p.current][0] // if arithmetic
}

// 現コマンドの2番目の引数を返す。現コマンドがC_PUSH,C_POP,C_FUNCTION,C_CALLの場合のみ
// 本メソッドを呼ぶようにする。
func (p *Parser) Arg2() int {
	return 0
}
