package evaluator

import (
	"testing"

	"github.com/estevesnp/dsb/pkg/object"
)

func TestExpandMacro(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`
let infixExpression = macro() { quote(1 + 2); };
infixExpression();`,
			"(1 + 2)",
		},
		{
			`
let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); };
reverse(2 + 2, 10 - 5);`,
			"(10 - 5) - (2 + 2)",
		},
		{
			`
let unless = macro(condition, consequence, alternative) {
    quote(if(!unquote(condition)) {
        unquote(consequence)
    } else {
        unquote(alternative)
    });
};

unless(10 > 5, print("not greater"), print("greater"));`,
			`if (!(10 > 5)) { print("not greater") } else { print("greater") }`,
		},
	}

	for _, tt := range tests {
		expected := testParseProgram(tt.expected)
		program := testParseProgram(tt.input)

		env := object.NewEnvironment()

		DefineMacros(program, env)
		expanded := ExpandMacros(program, env)

		if expanded.String() != expected.String() {
			t.Errorf("not equal. want %q, got %q", expected.String(), expanded.String())
		}
	}
}
