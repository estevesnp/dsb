package repl

import (
	"strings"
	"testing"

	"github.com/estevesnp/dsb/pkg/token"
)

func TestStart(t *testing.T) {
	out := &strings.Builder{}
	in := strings.NewReader("let five == 5;")

	expected := `[ Type: LET | Literal: let ]
[ Type: IDENT | Literal: five ]
[ Type: == | Literal: == ]
[ Type: INT | Literal: 5 ]
[ Type: ; | Literal: ; ]`

	Start(in, out)

	got := out.String()

	if ok := strings.Contains(got, expected); !ok {
		t.Fatalf("expected to contain:\n%q;\ngot:\n%q", expected, got)
	}
}

func TestFormatToken(t *testing.T) {
	tests := []struct {
		input    token.Token
		expected string
	}{
		{
			token.Token{Type: token.IF, Literal: "if"},
			"[ Type: IF | Literal: if ]",
		},
		{
			token.Token{Type: token.COMMA, Literal: ","},
			"[ Type: , | Literal: , ]",
		},
		{
			token.Token{Type: token.IDENT, Literal: "cena"},
			"[ Type: IDENT | Literal: cena ]",
		},
		{
			token.Token{Type: token.RETURN, Literal: "return"},
			"[ Type: RETURN | Literal: return ]",
		},
	}

	for i, tt := range tests {
		if got := formatToken(tt.input); got != tt.expected {
			t.Fatalf("tests[%d]: ", i)
		}
	}
}
