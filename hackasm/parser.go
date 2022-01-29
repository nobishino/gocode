package hackasm

import (
	"strconv"
	"strings"
)

type Parser struct {
	variableSymbolOffset uint64            // u = unsigned(符号なし整数) // 0
	variableSymbols      map[string]uint64 // nil
}

func NewParser() *Parser {
	return &Parser{
		variableSymbolOffset: 16,
		variableSymbols:      map[string]uint64{},
	}
}

// ProcessLine:

func (p *Parser) Parse(line string) Instruction {
	switch {
	case line[0] == '@':
		return p.parseA(line[1:])
	default:
		return p.parseC(line)
	}
}

// ParseLines はアセンブリソースコードsrcからInstruction
func (p *Parser) ParseLines(src string) Instructions {
	lines := strings.Split(src, "\n")
	var result Instructions
	for _, line := range lines {
		line = trimLine(line)
		if shouldSkip(line) {
			continue
		}
		result = append(result, p.Parse(line))
	}
	return result
}

// accept line = xxx in @xxx
func (p *Parser) parseA(aValue string) Instruction {
	return Instruction{
		Kind:  "A",
		Value: uint16(p.calcAValue(aValue)),
	}
}

func (p *Parser) calcAValue(aValue string) uint64 {
	n, err := strconv.ParseUint(aValue, 10, 15)
	if err == nil {
		return n
	}
	if aValue[0] == 'R' { // case R0, R1, R2, ... R15
		n, err = strconv.ParseUint(aValue[1:], 10, 15) // n = 0, 1, 2, ... 15
		if err == nil {
			return n
		}
	}
	switch aValue {
	case "SP":
		return 0
	case "LCL":
		return 1
	case "ARG":
		return 2
	case "THIS":
		return 3
	case "THAT":
		return 4
	case "SCREEN":
		return 16384
	case "KBD":
		return 24576
	}
	if _, ok := p.variableSymbols[aValue]; !ok {
		p.variableSymbols[aValue] = p.variableSymbolOffset
		p.variableSymbolOffset++
	}
	return p.variableSymbols[aValue]
}

func (p *Parser) parseC(line string) Instruction {
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

func trimLine(line string) string {
	inlineCommentIdx := strings.Index(line, "//")
	if inlineCommentIdx != -1 {
		line = line[:inlineCommentIdx]
	}
	line = strings.Trim(line, " ")
	line = strings.Trim(line, "\t")
	line = strings.Trim(line, "\r")
	return line
}

func shouldSkip(line string) bool {
	return line == ""
}
