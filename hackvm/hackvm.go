package hackvm

import (
	"errors"
	"flag"
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
	srcPath := args[0] // TODO handle multiple source files
	f, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer f.Close()
	p := parser.New(f)

	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	cw := codewriter.New(out)
	for p.HasMoreCommands() {
		p.Advance()
		switch p.CommandType() {
		case "C_ARITHMETIC":
			if err := cw.WriteArithmetic(p.Arg1()); err != nil {
				return err
			}
		case "C_PUSH":
			if err := cw.WritePushPop("C_PUSH", p.Arg1(), p.Arg2()); err != nil {
				return err
			}
		}
	}
	return nil
}
