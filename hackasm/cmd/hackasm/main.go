package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nobishino/gocode/hackasm"
)

var srcPath string
var outPath string

func init() {
	flag.StringVar(&srcPath, "src", "", "specify input file name")
	flag.StringVar(&outPath, "out", "out.hack", "specify output file name")
}

func exec(srcPath, outPath string) {
	f, err := os.Open(srcPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	lines := hackasm.Assemble(string(b))

	out, err := os.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	for _, line := range lines {
		_, err := fmt.Fprintln(out, line)
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Println("assembled into", outPath)
}

func main() {
	flag.Parse()
	exec(srcPath, outPath)
}
