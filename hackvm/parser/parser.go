package parser

import (
	"io"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

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
	label         = "label"
	goTo          = "goto"
	function      = "function"
	call          = "call"
	return_       = "return"
	ifGoTo        = "if-goto"
	cmdPush       = "C_PUSH"
	cmdPop        = "C_POP"
	cmdArithmetic = "C_ARITHMETIC"
	cmdLabel      = "C_LABEL"
	cmdGoto       = "C_GOTO"
	cmdIf         = "C_IF"
	cmdFunc       = "C_FUNCTION"
	cmdCall       = "C_CALL"
	cmdReturn     = "C_RETURN"
	invalidArg2   = -1
)

type Parser struct {
	cmds    [][]string
	current int
	command string
	arg1    string
	arg2    int
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
func (p *Parser) Advance() error {
	p.current++
	if err := p.validate(p.cmds[p.current]); err != nil {
		return errors.WithMessagef(err, "invalid command at %d th line:\n%s", p.current, p.cmds[p.current])
	}
	return nil
}

func (p *Parser) validate(cmd []string) error {
	if len(cmd) < 1 {
		return errors.Errorf("command should have at least 2 words, but got %d", len(cmd))
	}
	switch cmd[0] {
	case cmdPush, cmdPop:
		if len(cmd) != 3 {
			return errors.Errorf("push/pop command should have exactly 3 words, but got %d", len(cmd))
		}
		v, err := strconv.Atoi(cmd[2])
		if err != nil {
			return errors.Wrap(err, "failed to parse 2nd argument of push/pop command")
		}
		p.command = cmd[0]
		p.arg1 = cmd[1]
		p.arg2 = v
	case add, sub, neg, eq, gt, lt, and, or, not:
		if len(cmd) != 1 {
			return errors.Errorf("push command should have exactly 1 word, but got %d", len(cmd))
		}
		p.command = cmdArithmetic
		p.arg1 = cmd[0]
		p.arg2 = invalidArg2
	}
	return nil
}

// 現コマンドの種類を返す。算術コマンドはすべてC_ARITHMETICが返される
func (p *Parser) CommandType() string {
	switch p.cmds[p.current][0] {
	case push:
		return cmdPush
	case pop:
		return cmdPop
	case add, sub, neg, eq, gt, lt, and, or, not:
		return cmdArithmetic
	case label:
		return cmdLabel
	case goTo:
		return cmdGoto
	case ifGoTo:
		return cmdIf
	case function:
		return cmdFunc
	case call:
		return cmdCall
	case return_:
		return cmdReturn
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
