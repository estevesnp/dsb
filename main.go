package main

import (
	"fmt"
	"io"
	"os"

	"github.com/estevesnp/dsb/pkg/interpreter"
	"github.com/estevesnp/dsb/pkg/repl"
)

func main() {
	switch {

	case len(os.Args) == 1 && inputIsPiped():
		startInterpreter(os.Stdin)

	case len(os.Args) == 1:
		startREPL()

	default:
		for _, arg := range os.Args[1:] {
			processFile(arg)
		}
	}
}

func startREPL() {
	fmt.Print("this is the Domain Specific Bullshit REPL, have fun\n\n")
	repl.Start(os.Stdin, os.Stdout)
}

func startInterpreter(reader io.Reader) {
	if err := interpreter.Start(reader); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func inputIsPiped() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error checking stdin: %v\n", err)
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) == 0
}

func processFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file %q: %v", fileName, err)
	}
	defer file.Close()

	startInterpreter(file)
}
