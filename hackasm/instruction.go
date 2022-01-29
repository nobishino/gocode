package hackasm

import (
	"fmt"
	"strings"
)

type Instruction struct {
	Kind  string // A or C
	Value uint16 // only for A-instruction
	Dest  string
	Comp  string
	Jump  string
}

func (i Instruction) String() string {
	switch i.Kind {
	case "A":
		return fmt.Sprintf("{%s: %d}", i.Kind, i.Value)
	case "C":
		return fmt.Sprintf("{%s: %s=%s;%s}", i.Kind, i.Dest, i.Comp, i.Jump)
	default:
		panic("")
	}
}

func (i Instruction) Code() string {
	result := func() string {
		if i.Kind == "A" {
			return fmt.Sprintf("0%015b", i.Value)
		}
		return fmt.Sprintf("111%s%03b%03b", compMap[i.Comp], i.destEncode(), i.decodeJump())
	}()
	return result
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

type Instructions []Instruction
