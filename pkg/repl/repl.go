package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/estevesnp/dsb/pkg/evaluator"
	"github.com/estevesnp/dsb/pkg/lexer"
	"github.com/estevesnp/dsb/pkg/object"
	"github.com/estevesnp/dsb/pkg/parser"
)

const PROMPT = ">>  "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)

		if scanned := scanner.Scan(); !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		evaluated := evaluator.Eval(expanded, env)

		if evaluated == nil {
			continue
		}

		fmt.Fprintf(out, "%s\n", evaluated.Inspect())
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "\t%s\n", msg)
	}
}
