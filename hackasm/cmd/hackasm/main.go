package main

import (
	"flag"
	"log"
	"os"

	"github.com/nobishino/gocode/hackasm"
)

var outPath string

func init() {
	flag.StringVar(&outPath, "o", "out.hack", "specify output file name")
}

func exec(srcPath, outPath string) int {
	f, err := os.Open(srcPath)
	if err != nil {
		log.Println(err)
		return 1
	}
	defer f.Close()

	out, err := os.Create(outPath)
	if err != nil {
		log.Println(err)
		return 1
	}
	defer out.Close()

	if err := hackasm.Assemble(f, out); err != nil {
		log.Println(err)
		return 1
	}

	log.Println("assembled into", outPath)
	return 0
}

func main() {
	flag.Parse()
	srcPath := flag.Arg(0)
	stCode := exec(srcPath, outPath)
	os.Exit(stCode)
}
