package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nobishino/gocode/hackasm"
)

var outPath string
var srcPath string

func init() {
	flag.StringVar(&outPath, "o", "out.hack", "specify output file name")
	flag.StringVar(&srcPath, "src", "", "specify output file name")
	flag.Parse()
}

func main() {
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
	log.Println(string(b))
	log.Println(lines)

	out, err := os.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	for _, line := range lines {
		log.Println(line)
		_, err := fmt.Fprintln(out, line)
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Println("assembled into", outPath)
}
