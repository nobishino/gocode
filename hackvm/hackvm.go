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
	if len(os.Args) < 2 {
		return errors.New("need exactly 1 argument")
	}
	srcPath := os.Args[1]
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
