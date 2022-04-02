package codewriter

import (
	"fmt"

	"github.com/pkg/errors"
)

func (c *CodeWriter) WriteFunction(funcName string, numLocal int) error {
	format := `// function %[1]s %[2]d
(function_Test_%[1]s)
`
	initializeDRegister := `@0
D=A
`
	push := `@SP
A=M
M=D
@SP
M=M+1
`
	if _, err := fmt.Fprintf(c.out, format, funcName, numLocal); err != nil {
		return errors.WithStack(err)
	}
	if numLocal == 0 {
		return nil
	}
	if _, err := fmt.Fprint(c.out, initializeDRegister); err != nil {
		return errors.WithStack(err)
	}
	for i := 0; i < numLocal; i++ {
		if _, err := fmt.Fprint(c.out, push); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (c *CodeWriter) WriteReturn() error {
	code := `// return
// FRAME = LCL
@LCL
D=M
@R13
M=D
// RET = *(FRAME-5)
@5
A=D-A
D=M
@R14
M=D
// *ARG = pop()
@SP
M=M-1
@SP
A=M
D=M
@ARG
A=M
M=D
// SP = ARG+1
@ARG
D=M+1
@SP
M=D
// THAT = *(FRAME-1)
@1
A=D-A
D=M
@THAT
M=D
// THIS = *(FRAME-2)
@2
A=D-A
D=M
@THIS
M=D
// ARG = *(FRAME-3)
@3
A=D-A
D=M
@ARG
M=D
// LCL = *(FRAME-4)
@4
A=D-A
D=M
@LCL
M=D
// goto RET
@R14
A=M
0;JMP
`
	if _, err := fmt.Fprint(c.out, code); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
