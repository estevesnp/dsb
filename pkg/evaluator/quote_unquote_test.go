package evaluator

import (
	"testing"

	"github.com/estevesnp/dsb/pkg/object"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"quote(5)",
			"5",
		},
		{
			"quote(5 + 8)",
			"(5 + 8)",
		},
		{
			"quote(foobar)",
			"foobar",
		},
		{
			"quote(foobar + barfoo)",
			"(foobar + barfoo)",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		quote, ok := evaluated.(*object.Quote)
		if !ok {
			t.Fatalf("expected *object.Quote, got %T (%+v)", evaluated, evaluated)
		}

		if quote.Node == nil {
			t.Fatalf("qote.Node is nil")
		}

		if got := quote.Node.String(); got != tt.expected {
			t.Errorf("not equal. got %q, want %q", got, tt.expected)
		}
	}
}

func TestQuoteUnquote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"quote(unquote(4))",
			"4",
		},
		{
			"quote(unquote(4 + 4))",
			"8",
		},
		{
			"quote(8 + unquote(4 + 4))",
			"(8 + 8)",
		},
		{
			"quote(unquote(4 + 4) + 8)",
			"(8 + 8)",
		},
		{
			"let foobar = 8; quote(foobar)",
			"foobar",
		},
		{
			"let foobar = 8; quote(unquote(foobar))",
			"8",
		},
		{
			"quote(unquote(true))",
			"true",
		},
		{
			"quote(unquote(true == false))",
			"false",
		},
		{
			`quote(unquote("foobar"))`,
			"foobar",
		},
		{
			`quote(unquote("foo" + "bar"))`,
			"foobar",
		},
		{
			"quote(unquote(null))",
			"null",
		},
		{
			"quote(unquote([1, 2, 3]))",
			"[1, 2, 3]",
		},
		{
			"quote(unquote({ 1: true }))",
			"{1:true}",
		},
		{
			"quote(unquote(fn(x) { return x * 2; }))",
			"fn(x) return (x * 2);",
		},
		{
			"quote(unquote(quote(4 + 4)))",
			"(4 + 4)",
		},
		{
			`let quotedInfixExpression = quote(4 + 4);
		quote(unquote(4 + 4) + unquote(quotedInfixExpression))`,
			"(8 + (4 + 4))",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		quote, ok := evaluated.(*object.Quote)
		if !ok {
			t.Fatalf("expected *object.Quote, got %T (%+v)", evaluated, evaluated)
		}

		if quote.Node == nil {
			t.Fatalf("quote.Node is nil with input %q", tt.input)
		}

		if got := quote.Node.String(); got != tt.expected {
			t.Errorf("not equal. got %q, want %q", got, tt.expected)
		}
	}
}
