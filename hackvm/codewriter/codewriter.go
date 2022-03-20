package codewriter

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

type CodeWriter struct {
	out             io.Writer
	comparisonIndex int
	fileName        string
}

//出力ファイル/ストリームを開き書き込む準備を行う
func New(w io.Writer) *CodeWriter {
	return &CodeWriter{
		out: w,
	}
}

// CodeWriterに、あたらしいVMファイルの変換が開始したことを伝える
// .vm拡張子以外を渡すとpanicする
func (c *CodeWriter) SetFileName(name string) {
	if !strings.HasSuffix(name, ".vm") {
		panic(fmt.Sprintf("file name must have extension '.vm'. got: %q", name))
	}
	c.fileName = name[:len(name)-3]
}

//与えられた算術コマンドをアッセンブリーコードに変換し、それを書き込む
func (c *CodeWriter) WriteArithmetic(command string) error {
	var code string
	switch command {
	case "neg", "not":
		code = c.unaryArithmetic(command)
	case "add", "sub", "and", "or":
		code = c.binaryArithmetic(command)
	case "eq", "gt", "lt":
		code = c.comparison(command)
	default:
		return errors.Errorf("invalid command %q", command)
	}

	_, err := io.WriteString(c.out, code)
	if err != nil {
		return err
	}
	return nil
}

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

func (c *CodeWriter) codePushTemp(index int) string {
	const tempSegmentOffset = 5
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
	const tempSegmentOffset = 5
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

func (c *CodeWriter) unaryArithmetic(command string) string {
	comp := map[string]string{
		"neg": "-D",
		"not": "!D",
	}
	codeFmt := `// %[1]s
@SP
M=M-1
A=M
D=M
M=%[2]s
@SP
M=M+1
`
	return fmt.Sprintf(codeFmt, command, comp[command])
}

func (c *CodeWriter) binaryArithmetic(command string) string {
	comp := map[string]string{
		"add": "D+M",
		"sub": "M-D",
		"and": "D&M",
		"or":  "D|M",
	}
	codeFmt := `// %[1]s
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=%[2]s
@SP
M=M+1
`
	return fmt.Sprintf(codeFmt, command, comp[command])
}

func (c *CodeWriter) comparison(command string) string {
	jmpInst := map[string]string{
		"eq": "JEQ",
		"lt": "JLT",
		"gt": "JGT",
	}
	codeFmt := `// %[1]s
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=M-D
@%[1]s_%[2]d_true
D;%[3]s
(%[1]s_%[2]d_false)
@SP
A=M
M=0
@SP
M=M+1
@%[1]s_%[2]d_end
0;JMP
(%[1]s_%[2]d_true)
@SP
A=M
M=-1
@SP
M=M+1
(%[1]s_%[2]d_end)
`
	result := fmt.Sprintf(codeFmt, command, c.comparisonIndex, jmpInst[command])
	c.comparisonIndex++
	return result
}

// 出力ファイルを閉じる
func (c *CodeWriter) Close() error { return nil }
