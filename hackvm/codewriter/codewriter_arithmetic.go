package codewriter

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

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
