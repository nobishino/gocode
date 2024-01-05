package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
	f, err := os.Create("out.log")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	for i := 0; ; i++ {
		if i%10000 == 0 {
			log.Printf("trial %d\n", i)
		}
		if err := ex(f); err != nil {
			log.Fatalln(err)
		}
	}
}

func ex(w io.Writer) error {
	c := exec.Command("samevariable")
	c.Stdout = w
	c.Stderr = w
	return c.Run()
}
