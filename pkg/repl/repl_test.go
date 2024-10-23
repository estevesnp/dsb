package repl

import (
	"strings"
	"testing"
)

func TestStart(t *testing.T) {
	t.Skip()

	tests := []struct {
		input    string
		expected string
	}{
		{"let x = 42;", "let x = 42;"},
		{"return -1337;", "return (-1337);"},
		{"x * y / 2 + 3 * 8 - 123", "((((x * y) / 2) + (3 * 8)) - 123)"},
		{"true!=false", "true != false"},
		{" let        x = foo    <= 2", "let x = (foo <= 2);"},
	}

	for _, tt := range tests {

		out := &strings.Builder{}
		in := strings.NewReader(tt.input)

		Start(in, out)

		got := out.String()

		if ok := strings.Contains(got, tt.expected); !ok {
			t.Fatalf("expected to contain:\n%q;\ngot:\n%q", tt.expected, got)
		}
	}
}
