package parser

import (
	"fmt"
	"testing"

	"github.com/estevesnp/dsb/pkg/ast"
	"github.com/estevesnp/dsb/pkg/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if n := len(program.Statements); n != 1 {
			t.Fatalf("program.Statements doesn't have 1 statement. got %d", n)
		}

		stmt := program.Statements[0]

		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value

		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
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
		t.Errorf("letStmt.Name.Value not '%s'. got %q", name, letStmt.Name.Value)
		return false
	}

	if literal := letStmt.Name.TokenLiteral(); literal != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got %s", name, literal)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue any
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foo;", "foo"},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		if n := len(program.Statements); n != 1 {
			t.Fatalf("program.Statements doesn't have 1 statement. got %d", n)
		}

		stmt := program.Statements[0]

		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got %T", stmt)
			continue
		}

		if literal := returnStmt.TokenLiteral(); literal != "return" {
			t.Errorf("returnStmt.TokenLiteral() not 'return', got %q", literal)
			continue
		}

		if !testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
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

func TestStringLiteralExpressions(t *testing.T) {
	input := `"hello world";`

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

	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.StringLiteral. got %T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("ident.Value not %q. got %q", "hello world", literal.Value)
	}

	if tokLiteral := literal.TokenLiteral(); tokLiteral != "hello world" {
		t.Errorf("ident.TokenLiteral() not %q. got %q", "hello world", tokLiteral)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

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

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.ArrayLiteral. got %T", stmt.Expression)
	}

	if n := len(array.Elements); n != 3 {
		t.Fatalf("len(array.Elements) not 3. got %d", n)
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingMapLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

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

	mapLiteral, ok := stmt.Expression.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.MapLiteral. got %T", stmt.Expression)
	}

	if n := len(mapLiteral.Pairs); n != 3 {
		t.Errorf("mapLiteral.Pairs has wrong length. got %d", n)
	}

	for key, value := range mapLiteral.Pairs {
		keyLiteral, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not *ast.StringLiteral. got %T", key)
			continue
		}

		expectedValue := expected[keyLiteral.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingMapLiteralsIntegerKeys(t *testing.T) {
	input := `{1: 1, 2: 2, 3: 3}`

	expected := map[int64]int64{
		1: 1,
		2: 2,
		3: 3,
	}

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

	mapLiteral, ok := stmt.Expression.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.MapLiteral. got %T", stmt.Expression)
	}

	if n := len(mapLiteral.Pairs); n != 3 {
		t.Errorf("mapLiteral.Pairs has wrong length. got %d", n)
	}

	for key, value := range mapLiteral.Pairs {
		keyLiteral, ok := key.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("key is not *ast.IntegerLiteral. got %T", key)
			continue
		}

		expectedValue := expected[keyLiteral.Value]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingMapLiteralsBooleanKeys(t *testing.T) {
	input := `{true: 1, false: 2}`

	expected := map[bool]int64{
		true:  1,
		false: 2,
	}

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

	mapLiteral, ok := stmt.Expression.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.MapLiteral. got %T", stmt.Expression)
	}

	if n := len(mapLiteral.Pairs); n != 2 {
		t.Errorf("mapLiteral.Pairs has wrong length. got %d", n)
	}

	for key, value := range mapLiteral.Pairs {
		keyLiteral, ok := key.(*ast.Boolean)
		if !ok {
			t.Errorf("key is not *ast.Boolean. got %T", key)
			continue
		}

		expectedValue := expected[keyLiteral.Value]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingMapLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

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

	mapLiteral, ok := stmt.Expression.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.MapLiteral. got %T", stmt.Expression)
	}

	if n := len(mapLiteral.Pairs); n != 3 {
		t.Errorf("mapLiteral.Pairs has wrong length. got %d", n)
	}

	for key, value := range mapLiteral.Pairs {
		keyLiteral, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not *ast.StringLiteral. got %T", key)
			continue
		}

		testFunc, ok := tests[keyLiteral.String()]
		if !ok {
			t.Errorf("No test function for key %q found", keyLiteral.String())
			continue
		}

		testFunc(value)
	}
}

func TestParsingEmptyMapLiteral(t *testing.T) {
	input := "{}"

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

	mapLiteral, ok := stmt.Expression.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.MapLiteral. got %T", stmt.Expression)
	}

	if n := len(mapLiteral.Pairs); n != 0 {
		t.Errorf("mapLiteral.Pairs has wrong length. got %d", n)
	}
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

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

	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IndexExpression. got %T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"

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

	literal, ok := stmt.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.Boolean. got %T", stmt.Expression)
	}

	if literal.Value != true {
		t.Errorf("ident.Value not %t. got %t", true, literal.Value)
	}

	if tokLiteral := literal.TokenLiteral(); tokLiteral != "true" {
		t.Errorf("ident.TokenLiteral() not %s. got %s", "true", tokLiteral)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		if n := len(program.Statements); n != 1 {
			t.Fatalf("program.Statements doesn't have 1 statement. got %d", n)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not *ast.PrefixExpression. got %T", stmt)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s. got %s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"5 >= 5", 5, ">=", 5},
		{"5 <= 5", 5, "<=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements, got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %t", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 >= 4 * 2",
			"((3 + 4) >= (4 * 2))",
		},
		{
			"3 - 4 <= 4 / 2",
			"((3 - 4) <= (4 / 2))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if n := len(program.Statements); n != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got %d", 1, n)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %t", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got %t", program.Statements[0])
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if n := len(exp.Consequence.Statements); n != 1 {
		t.Fatalf("consequence is not 1 statement. got %d", n)
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got %T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got %+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if n := len(program.Statements); n != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got %d", 1, n)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %t", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got %t", program.Statements[0])
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if n := len(exp.Consequence.Statements); n != 1 {
		t.Fatalf("consequence is not 1 statement. got %d", n)
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence.Statements[0] is not ast.ExpressionStatement. got %T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if n := len(exp.Alternative.Statements); n != 1 {
		t.Fatalf("consequence is not 1 statement. got %d", n)
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative.Statements[0] is not ast.ExpressionStatement. got %T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if n := len(program.Statements); n != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got %d", 1, n)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %t", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.FunctionLiteral. got %t", stmt.Expression)
	}

	if n := len(function.Parameters); n != 2 {
		t.Fatalf("function literal parameters wrong. want %d, got %d", 2, n)
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if n := len(function.Body.Statements); n != 1 {
		t.Fatalf("function.Body.Statements wrong. want %d, got %d", 1, n)
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not ast.ExpressionStatement. got %t", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if nFunc, nExpected := len(function.Parameters), len(tt.expectedParams); nFunc != nExpected {
			t.Errorf("length parameters wrong. want %d, got %d", nExpected, nFunc)
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if n := len(program.Statements); n != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got %d", 1, n)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got %t", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.CallExpression. got %T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if n := len(exp.Arguments); n != 3 {
		t.Fatalf("wrong length of arguments. wanted %d, got %d", 3, n)
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
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

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got %T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got %d", value, integ.Value)
		return false
	}

	if tok := integ.TokenLiteral(); tok != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral() not %d. got %s", value, tok)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got %T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got %s", value, ident.Value)
		return false
	}

	if tok := ident.TokenLiteral(); tok != value {
		t.Errorf("ident.TokenLiteral() not %s. got %s", value, tok)
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("type of exp not handled. got %T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got %T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not %s. got %q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp is not ast.Boolean. got %T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got %t", value, bo.Value)
		return false
	}

	if tok := bo.TokenLiteral(); tok != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral() not %t. got %s", value, tok)
		return false
	}

	return true
}
