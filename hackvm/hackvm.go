package hackvm

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nobishino/gocode/hackvm/codewriter"
	"github.com/nobishino/gocode/hackvm/parser"
	"github.com/pkg/errors"
)

func Exec() int {
	if err := exec(); err != nil {
		log.Println(err)
		return 1
	}
	return 0
}

func exec() error {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		return errors.New("need at least 1 argument")
	}
	arg := args[0]
	outFile := strings.TrimSuffix(filepath.Base(arg), filepath.Ext(arg)) + ".asm"

	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()
	cw := codewriter.New(out)

	src, err := os.Stat(arg)
	if err != nil {
		return errors.WithStack(err)
	}

	if src.IsDir() {
		dir, err := os.ReadDir(arg)
		if err != nil {
			return errors.WithStack(err)
		}

		// Initialization
		cw.Init()

		for _, src := range dir {
			if filepath.Ext(src.Name()) != "vm" {
				continue
			}
			if err := func() error {
				src, err := os.Open(filepath.Join(arg, src.Name()))
				if err != nil {
					return err
				}
				defer src.Close()

				if err := Translate(cw, src, src.Name()); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return err
			}
		}
		return nil
	}

	vmFile, err := os.Open(arg)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := Translate(cw, vmFile, vmFile.Name()); err != nil {
		return err
	}

	return nil
}

func Translate(cw *codewriter.CodeWriter, r io.Reader, fileName string) error {

	p := parser.New(r)
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
		case "C_FUNCTION":
			if err := cw.WriteFunction(p.Arg1(), p.Arg2()); err != nil {
				return err
			}
		case "C_CALL":
			if err := cw.WriteCall(p.Arg1(), p.Arg2()); err != nil {
				return err
			}
		case "C_RETURN":
			if err := cw.WriteReturn(); err != nil {
				return err
			}
		}
	}
	return nil
}
