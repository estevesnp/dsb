package interpreter

import (
	"fmt"
	"io"
	"strings"

	"github.com/estevesnp/dsb/pkg/evaluator"
	"github.com/estevesnp/dsb/pkg/lexer"
	"github.com/estevesnp/dsb/pkg/object"
	"github.com/estevesnp/dsb/pkg/parser"
)

func Start(reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("error loading program: %w", err)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()

	if errs := p.Errors(); len(errs) != 0 {
		return fmt.Errorf("error parsing program:\n\t%s", strings.Join(errs, "\n\t"))
	}

	env := object.NewEnvironment()

	res := evaluator.Eval(program, env)
	if err, ok := res.(*object.Error); ok {
		return fmt.Errorf("error evaluating the program: %s", err.Message)
	}

	return nil
}
