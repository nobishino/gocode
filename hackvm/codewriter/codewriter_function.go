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
