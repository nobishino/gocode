package codewriter

import (
	"fmt"
	"io"
)

//C_PUSHまたはC_POPコマンドをアッセンブリーに変換しそれを書き込む
func (c *CodeWriter) WritePushPop(command string, segment string, index int) error {
	var code string
	switch command {
	case "C_POP":
		switch segment {
		case "static":
			code = c.codePopStatic(index)
		case "local", "argument", "this", "that":
			code = c.codePopLocal(segment, index)
		case "temp":
			code = c.codePopTemp(index)
		}
	case "C_PUSH":
		switch segment {
		case "constant":
			code = c.codePushConstant(index)
		case "static":
			code = c.codePushStatic(index)
		case "local", "argument", "this", "that":
			code = c.codePushLocal(segment, index)
		case "temp":
			code = c.codePushTemp(index)
		}
	}
	if code == "" {
		panic(fmt.Sprintf("undefined. command: %q, segment %q", command, segment))
	}
	_, err := io.WriteString(c.out, code)
	if err != nil {
		return err
	}

	return nil
}

func (c *CodeWriter) codePushConstant(index int) string {
	format := `// push constant %[1]d
@%[1]d
D=A
@SP
A=M
M=D
@SP
M=M+1
`
	return fmt.Sprintf(format, index)
}

func (c *CodeWriter) codePopStatic(index int) string {
	format := `// pop static %[1]d
@SP
M=M-1
A=M
D=M
@%[2]s.%[1]d
M=D
`
	return fmt.Sprintf(format, index, c.fileName)
}

func (c *CodeWriter) codePushStatic(index int) string {
	format := `// push static %[1]d
@%[2]s.%[1]d
D=M
@SP
A=M
M=D
@SP
M=M+1
`
	return fmt.Sprintf(format, index, c.fileName)
}

var segmentToSymbol = map[string]string{
	"local":    "LCL",
	"argument": "ARG",
	"this":     "THIS",
	"that":     "THAT",
}

func (c *CodeWriter) codePopLocal(segment string, index int) string {
	baseAddrSymbol, ok := segmentToSymbol[segment]
	if !ok {
		panic(fmt.Sprintf("invalid segment name %q", segment))
	}
	format := `// pop %[2]s %[1]d
@%[3]s
D=M
@%[1]d
D=D+A
@R13 // general register
M=D
@SP
M=M-1
A=M
D=M
@R13 // general register
A=M
M=D
`
	return fmt.Sprintf(format, index, segment, baseAddrSymbol)
}

func (c *CodeWriter) codePushLocal(segment string, index int) string {
	baseAddrSymbol, ok := segmentToSymbol[segment]
	if !ok {
		panic(fmt.Sprintf("invalid segment name %q", segment))
	}
	format := `// push %[2]s %[1]d
@%[1]d
D=A
@%[3]s
D=D+M
A=D
D=M
@SP
A=M
M=D
@SP
M=M+1
`
	return fmt.Sprintf(format, index, segment, baseAddrSymbol)
}

const tempSegmentOffset = 5

func (c *CodeWriter) codePushTemp(index int) string {
	format := `// push temp %d
@R%d
D=M
@SP
A=M
M=D
@SP
M=M+1
`
	return fmt.Sprintf(format, index, index+tempSegmentOffset)
}

func (c *CodeWriter) codePopTemp(index int) string {
	format := `// pop temp %d
@SP
M=M-1
A=M
D=M
@R%d
M=D
`
	return fmt.Sprintf(format, index, index+tempSegmentOffset)
}
