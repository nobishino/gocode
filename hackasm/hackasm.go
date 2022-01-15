package hackasm

import (
	"fmt"
	"strconv"
	"strings"
)

type Instruction struct {
	Kind  string // A or C
	Value uint16 // only for A-instruction
	Dest  string
	Comp  string
	Jump  string
}

func (i Instruction) Code() string {
	if i.Kind == "A" {
		return fmt.Sprintf("0%015b", i.Value)
	}
	return fmt.Sprintf("111%s%03b%s", compMap[i.Comp], i.destEncode(), "000")
}

// destの部分をエンコードする
// 戻り値は3bit整数
func (i Instruction) destEncode() uint8 {
	var result uint8 // = 0
	if strings.Contains(i.Dest, "A") {
		result += 1 << 2
	}
	if strings.Contains(i.Dest, "D") {
		result += 1 << 1
	}
	if strings.Contains(i.Dest, "M") {
		result += 1
	}
	return result
}

var compMap = map[string]string{
	"D+A": "0000010",
	"D&M": "1000000",
}

// ProcessLine:

func Parse(line string) Instruction {
	switch {
	case line[0] == '@':
		return parseA(line)
	default:
		return parseC(line)
	}
}

func parseA(line string) Instruction {
	v, err := strconv.ParseUint(line[1:], 10, 15)
	if err != nil {
		panic(err)
	}
	return Instruction{
		Kind:  "A",
		Value: uint16(v),
	}
}

func parseC(line string) Instruction {
	result := Instruction{
		Kind: "C",
	}
	var head int

	if assign := strings.Index(line, "="); assign != -1 {
		result.Dest = line[:assign] // = の手前までとりたいのでスライシング
		head = assign + 1
	}

	if semicolon := strings.Index(line, ";"); semicolon != -1 {
		result.Comp = line[head:semicolon]
		result.Jump = line[semicolon+1:]
		return result
	}

	result.Comp = line[head:]

	return result
}
