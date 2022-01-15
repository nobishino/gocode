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
	return fmt.Sprintf("111%s%03b%03b", compMap[i.Comp], i.destEncode(), i.decodeJump())
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
	// a = 0
	"0":   "0101010",
	"1":   "0111111",
	"-1":  "0111010",
	"D":   "0001100",
	"A":   "0110000",
	"!D":  "0001101",
	"!A":  "0110001",
	"-D":  "0001111",
	"-A":  "0110011",
	"D+1": "0011111",
	"A+1": "0110111",
	"D+A": "0000010",
	"A+D": "0000010",
	"D-A": "0010011",
	"A-D": "0000111",
	"D&A": "0000000",
	"A&D": "0000000",
	"D|A": "0010101",
	"A|D": "0010101",
	// a = 1
	"M":   "1110000",
	"!M":  "1110001",
	"-M":  "1110011",
	"M+1": "1110111",
	"D+M": "1000010",
	"M+D": "1000010",
	"D-M": "1010011",
	"M-D": "1000111",
	"D&M": "1000000",
	"M&D": "1000000",
	"D|M": "1010101",
	"M|D": "1010101",
}

// Jump部分を3bit整数にエンコードする
func (i Instruction) decodeJump() uint8 {
	switch i.Jump {
	case "":
		return 0b000
	case "JGT":
		return 0b001
	case "JEQ":
		return 0b010
	case "JGE":
		return 0b011
	case "JLT":
		return 0b100
	case "JNE":
		return 0b101
	case "JLE":
		return 0b110
	case "JMP":
		return 0b111
	default:
		return 0
	}
}

// ProcessLine:

func Parse(line string) Instruction {
	switch {
	case line[0] == '@':
		return parseA(line[1:])
	default:
		return parseC(line)
	}
}

// accept line = xxx in @xxx
func parseA(aValue string) Instruction {
	return Instruction{
		Kind:  "A",
		Value: uint16(calcAValue(aValue)),
	}
}

var variableSymbolOffset uint64 = 16
var variableSymbols = map[string]uint64{}

func calcAValue(aValue string) uint64 {
	n, err := strconv.ParseUint(aValue, 10, 15)
	if err == nil {
		return n
	}
	if aValue[0] == 'R' {
		n, err = strconv.ParseUint(aValue[1:], 10, 15)
		if err != nil {
			panic(err)
		}
		return n
	}
	if _, ok := variableSymbols[aValue]; !ok {
		variableSymbols[aValue] = variableSymbolOffset
		variableSymbolOffset++
	}
	return variableSymbols[aValue]
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

type Instructions []Instruction

// ParseLines はアセンブリソースコードsrcからInstruction
func ParseLines(src string) Instructions {
	lines := strings.Split(src, "\n")
	var result Instructions
	for _, line := range lines {
		line = trimLine(line)
		if shouldSkip(line) {
			continue
		}
		result = append(result, Parse(line))
	}
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

func Assemble(src string) []string {
	var result []string
	var instructions = ParseLines(src)
	for _, instruction := range instructions {
		result = append(result, instruction.Code())
	}
	return result
}
