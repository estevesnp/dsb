package main

import (
	"fmt"
	"io"
	"os"

	"github.com/estevesnp/dsb/pkg/interpreter"
	"github.com/estevesnp/dsb/pkg/repl"
)

func main() {
	if len(os.Args) == 1 {
		startRepl()
		return
	}

	for _, arg := range os.Args[1:] {
		file, err := os.Open(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening file %q: %v", arg, err)
			os.Exit(1)
		}

		startInterpreter(file)
	}
}

func startRepl() {
	fmt.Print("this is the Domain Specific Bullshit REPL, have fun\n\n")
	repl.Start(os.Stdin, os.Stdout)
}

func startInterpreter(reader io.Reader) {
	err := interpreter.Start(reader)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
