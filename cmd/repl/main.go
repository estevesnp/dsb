package main

import (
	"fmt"
	"os"

	"github.com/estevesnp/dsb/pkg/repl"
)

func main() {
	fmt.Print("this is the Domain Specific Bullshit REPL, have fun\n\n")

	repl.Start(os.Stdin, os.Stdout)
}
