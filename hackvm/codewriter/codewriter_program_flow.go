package codewriter

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// WriteLabel labelコマンドを行うアセンブリコードを書く
func (c *CodeWriter) WriteLabel(label string) error {
	if err := c.validateLabel(label); err != nil {
		return err
	}
	format := `// label %[1]s
(%[1]s)
`
	if _, err := fmt.Fprintf(c.out, format, c.currentFuncName+"$"+label); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// WriteGoto labelコマンドのアセンブリコードを書く
func (c *CodeWriter) WriteGoto(label string) error {
	if err := c.validateLabel(label); err != nil {
		return err
	}
	format := `// goto %[1]s
@%[1]s
0;JMP
`
	if _, err := fmt.Fprintf(c.out, format, c.generateLabel(label)); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// WriteIf labelコマンドのアセンブリコードを書く
func (c *CodeWriter) WriteIf(label string) error {
	if err := c.validateLabel(label); err != nil {
		return err
	}
	format := `// if-goto %[1]s
@SP
M=M-1
A=M
D=M
@%[1]s
D;JNE
`
	if _, err := fmt.Fprintf(c.out, format, c.generateLabel(label)); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *CodeWriter) validateLabel(label string) error {
	if label == "" {
		return errors.Errorf("invalid label %q", label)
	}
	if strings.TrimLeft(label, "0123456789") != label {
		return errors.Errorf("invalid label %q", label)
	}
	return nil
}

func (c *CodeWriter) generateLabel(label string) string {
	return c.currentFuncName + "$" + label
}
