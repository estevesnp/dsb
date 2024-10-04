package repl

import (
	"strings"
	"testing"
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
