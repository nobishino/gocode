package hackasm

import (
	"fmt"
	"io"
	"strings"
)

func Assemble(src string) []string {
	var result []string

	p := NewParser()
	var instructions = p.ParseLines(src)
	for _, instruction := range instructions {
		result = append(result, instruction.Code())
	}
	return result
}

func AssembleRW(r io.Reader, w io.Writer) error {
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, r); err != nil {
		return err
	}

	lines := Assemble(buf.String())
	for _, line := range lines {
		_, err := fmt.Fprintln(w, line)
		if err != nil {
			return err
		}
	}
	return nil
}
