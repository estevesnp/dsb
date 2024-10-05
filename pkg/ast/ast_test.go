package ast

import (
	"testing"

	"github.com/estevesnp/dsb/pkg/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "dsb"},
					Value: "dsb",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "dsb"},
					Value: "otherDsb",
				},
			},
		},
	}
	expected := "let dsb = otherDsb;"

	if got := program.String(); got != expected {
		t.Errorf("program.String() wrong. got %q, expected %q", got, expected)
	}
}
