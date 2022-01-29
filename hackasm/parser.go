package hackasm

import (
	"strconv"
	"strings"
)

type Parser struct{}

var variableSymbolOffset uint64 = 16
var variableSymbols = map[string]uint64{}

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
	if aValue[0] == 'R' {
		n, err = strconv.ParseUint(aValue[1:], 10, 15)
		if err == nil {
			return n
		}
	}
	if _, ok := variableSymbols[aValue]; !ok {
		variableSymbols[aValue] = variableSymbolOffset
		variableSymbolOffset++
	}
	return variableSymbols[aValue]
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
	return line
}

func shouldSkip(line string) bool {
	return line == ""
}
