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

func exec(srcPath, outPath string) {
	f, err := os.Open(srcPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	out, err := os.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	if err := hackasm.AssembleRW(f, out); err != nil {
		log.Fatalln(err)
	}

	log.Println("assembled into", outPath)
}

func main() {
	flag.Parse()
	srcPath := flag.Arg(0)
	exec(srcPath, outPath)
}
