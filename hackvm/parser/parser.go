package parser

import (
	"io"
	"strconv"
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

const (
	add           = "add"
	sub           = "sub"
	neg           = "neg"
	eq            = "eq"
	gt            = "gt"
	lt            = "lt"
	and           = "and"
	or            = "or"
	not           = "not"
	push          = "push"
	pop           = "pop"
	cmdPush       = "C_PUSH"
	cmdArithmetic = "C_ARITHMETIC"
)

// 現コマンドの種類を返す。算術コマンドはすべてC_ARITHMETICが返される
func (p *Parser) CommandType() string {
	switch p.cmds[p.current][0] {
	case push:
		return cmdPush
	case add, sub, neg, eq, gt, lt, and, or, not:
		return cmdArithmetic
	}
	panic("undefined command type") // TODO: validate when advance is invoked
}

// 現コマンドの最初の引数を返す。　C_ARITHMETICの場合、コマンド自体(add,subなど)が返される。
// 現コマンドがC_RETURNの場合、本メソッドは呼ばないようにする
func (p *Parser) Arg1() string {
	if p.CommandType() == "C_ARITHMETIC" {
		return p.cmds[p.current][0] // if arithmetic
	}
	return p.cmds[p.current][1]
}

// 現コマンドの2番目の引数を返す。現コマンドがC_PUSH,C_POP,C_FUNCTION,C_CALLの場合のみ
// 本メソッドを呼ぶようにする。
func (p *Parser) Arg2() int {
	v, err := strconv.Atoi(p.cmds[p.current][2])
	if err != nil {
		panic(err) // FIXME validate when advance
	}
	return v
}
