package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/estevesnp/dsb/pkg/lexer"
	"github.com/estevesnp/dsb/pkg/token"
)

const PROMPT = ">>  "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)

		if scanned := scanner.Scan(); !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintln(out, formatToken(tok))
		}
		fmt.Fprintln(out)
	}
}

func formatToken(tok token.Token) string {
	return fmt.Sprintf("[ Type: %s | Literal: %s ]", tok.Type, tok.Literal)
}
