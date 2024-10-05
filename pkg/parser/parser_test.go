package parser

import (
	"testing"

	"github.com/estevesnp/dsb/pkg/ast"
	"github.com/estevesnp/dsb/pkg/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
    `
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if n := len(program.Statements); n != 3 {
		t.Fatalf("program.Statements doesn't have 3 statements, got %d", n)
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 1337420;
    `
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if n := len(program.Statements); n != 3 {
		t.Fatalf("program.Statements doesn't have 3 statements, got %d", n)
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got %T", stmt)
			continue
		}

		if literal := returnStmt.TokenLiteral(); literal != "return" {
			t.Errorf("returnStmt.TokenLiteral() not 'return', got %q", literal)
			continue
		}
	}
}

func TestIdentifierExpressions(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if n := len(program.Statements); n != 1 {
		t.Fatalf("program.Statements doesn't have 1 statement, got %d", n)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("stmt.Expression is not *ast.Identifier. got %T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got %s", "foobar", ident.Value)
	}

	if literal := ident.TokenLiteral(); literal != "foobar" {
		t.Errorf("ident.TokenLiteral() not %s. got %s", "foobar", literal)
	}
}

func TestIntegerLiteralExpressions(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if n := len(program.Statements); n != 1 {
		t.Fatalf("program.Statements doesn't have 1 statement, got %d", n)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IntegerLiteral. got %T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("ident.Value not %d. got %d", 5, literal.Value)
	}

	if tokLiteral := literal.TokenLiteral(); tokLiteral != "5" {
		t.Errorf("ident.TokenLiteral() not %s. got %s", "5", tokLiteral)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	n := len(errors)

	if n == 0 {
		return
	}

	t.Errorf("parser has %d errors", n)

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if literal := s.TokenLiteral(); literal != "let" {
		t.Errorf("s.TokenLiteral() not 'let', got %s", literal)
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s', got %q", name, letStmt.Name.Value)
		return false
	}

	if literal := letStmt.Name.TokenLiteral(); literal != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s', got %s", name, literal)
		return false
	}

	return true
}
