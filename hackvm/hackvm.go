package hackvm

import (
	"errors"
	"flag"
	"io"
	"log"
	"os"

	"github.com/nobishino/gocode/hackvm/codewriter"
	"github.com/nobishino/gocode/hackvm/parser"
)

func Exec() int {
	if err := exec(); err != nil {
		log.Println(err)
		return 1
	}
	return 0
}

var outFile string

func init() {
	flag.StringVar(&outFile, "out", "out.asm", "specify output assembly file name")
}

func exec() error {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		return errors.New("need at least 1 argument")
	}

	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()
	for _, srcPath := range args {
		srcPath := srcPath
		if err := func() error {
			src, err := os.Open(srcPath)
			if err != nil {
				return err
			}
			defer src.Close()

			if err := Translate(out, src, srcPath); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			return err
		}
	}

	return nil
}

func Translate(w io.Writer, r io.Reader, fileName string) error {
	p := parser.New(r)
	cw := codewriter.New(w)
	cw.SetFileName(fileName)

	for p.HasMoreCommands() {
		if err := p.Advance(); err != nil {
			return err
		}
		switch p.CommandType() {
		case "C_ARITHMETIC":
			if err := cw.WriteArithmetic(p.Arg1()); err != nil {
				return err
			}
		case "C_PUSH", "C_POP":
			if err := cw.WritePushPop(p.CommandType(), p.Arg1(), p.Arg2()); err != nil {
				return err
			}
		case "C_LABEL":
			if err := cw.WriteLabel(p.Arg1()); err != nil {
				return err
			}
		case "C_GOTO":
			if err := cw.WriteGoto(p.Arg1()); err != nil {
				return err
			}
		case "C_IF":
			if err := cw.WriteIf(p.Arg1()); err != nil {
				return err
			}
		}
	}
	return nil
}
