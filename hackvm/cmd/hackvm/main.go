package main

import (
	"os"

	"github.com/nobishino/gocode/hackvm"
)

func main() {
	os.Exit(hackvm.Exec())
}
