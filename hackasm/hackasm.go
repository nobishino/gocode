package hackasm

import "strings"

type Instruction struct {
	Kind  string // A or C
	Value uint16 // only for A-instruction
	Dest  string
	Comp  string
	Jump  string
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
	return Instruction{
		Kind:  "A",
		Value: 15,
	}
}

func parseC(line string) Instruction {
	result := Instruction{
		Kind: "C",
	}
	var head int

	if assign := strings.Index(line, "="); assign != -1 {
		result.Dest = line[:assign]
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
