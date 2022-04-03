package codewriter

import (
	"fmt"

	"github.com/pkg/errors"
)

func (c *CodeWriter) WriteFunction(funcName string, numLocal int) error {
	c.currentFuncName = funcName // label, goto, if-gotoコマンドのために必要
	format := `// function %[1]s %[2]d
(function_%[1]s)
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
@R13
D=M
@1
A=D-A
D=M
@THAT
M=D
// THIS = *(FRAME-2)
@R13
D=M
@2
A=D-A
D=M
@THIS
M=D
// ARG = *(FRAME-3)
@R13
D=M
@3
A=D-A
D=M
@ARG
M=D
// LCL = *(FRAME-4)
@R13
D=M
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

func (c *CodeWriter) WriteCall(funcName string, argCount int) error {
	if _, err := fmt.Fprint(c.out, c.codeCall(funcName, argCount)); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *CodeWriter) codeCall(funcName string, argCount int) string {
	// %[1]s = funcName
	// %[2]s = argCount
	// %[3]s = current file name
	// %[4]d = globally unique key for "call" VM Command
	// %[5]d = argCount+5
	format := `// call %[1]s %[2]d
// push return-address 
@return_address_%[4]d
D=A
@SP
A=M
M=D
@SP
M=M+1
// push LCL
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1
// push ARG
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
// push THIS
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
// push THAT
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1
// ARG = SP-n-5
@%[5]d
D=A
@SP
D=M-D
@ARG
M=D
// LCL = SP
@SP
D=M
@LCL
M=D
// goto %[1]s
@function_%[1]s
0;JMP
(return_address_%[4]d)
`
	const frameHeight = 5
	callCount := c.callCount // globally unique key
	c.callCount++
	return fmt.Sprintf(format, funcName, argCount, c.fileName, callCount, argCount+frameHeight)
}
