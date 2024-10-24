package evaluator

import (
	"errors"
	"testing"

	"github.com/estevesnp/dsb/pkg/lexer"
	"github.com/estevesnp/dsb/pkg/object"
	"github.com/estevesnp/dsb/pkg/parser"
)

func TestEvalNull(t *testing.T) {
	tests := []string{
		"null",
		"fn(){}()",
		"let func = fn(){}; func()",
	}
	for _, input := range tests {
		testNullObject(t, testEval(input))
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"let x = 10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got %d, want %d", result.Value, expected)
		return false
	}

	return true
}

func TestIntegerCache(t *testing.T) {
	tests := []struct {
		input        string
		expectedSame bool
	}{
		{"-128", true},
		{"128", true},
		{"0", true},
		{"1", true},
		{"-1", true},
		{"5", true},
		{"-5", true},
		{"-129", false},
		{"129", false},
		{"2024", false},
	}

	for _, tt := range tests {
		firstInt := testEval(tt.input).(*object.Integer)
		secondInt := testEval(tt.input).(*object.Integer)

		if isSame := firstInt == secondInt; isSame != tt.expectedSame {
			t.Errorf("expected object equality with value %d to be %t, was %t",
				firstInt.Value, tt.expectedSame, isSame)
		}
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	expected := "Hello World!"

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got %T (%+v)", evaluated, evaluated)
	}

	if str.Value != expected {
		t.Errorf("String is not %q, got %q", expected, str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	expected := "Hello World!"

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got %T (%+v)", evaluated, evaluated)
	}

	if str.Value != expected {
		t.Errorf("String is not %q, got %q", expected, str.Value)
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	arr, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got %T (%+v)", evaluated, evaluated)
	}

	if n := len(arr.Elements); n != 3 {
		t.Fatalf("Array has wrong number of elements. got %d", n)
	}

	testIntegerObject(t, arr.Elements[0], 1)
	testIntegerObject(t, arr.Elements[1], 4)
	testIntegerObject(t, arr.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i]",
			1,
		},

		{
			"[1, 2, 3][1 + 1]",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[1];",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestMapLiterals(t *testing.T) {
	input := `let two = "two";
{
    "one": 10 - 9,
    two: 1 + 1,
    "thr" + "ee": 6 / 2,
    4: 4,
    true: 5,
    false: 6
}`

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Map)
	if !ok {
		t.Fatalf("Eval didn't return Map. got %T (%+v)", evaluated, evaluated)
	}

	if n := len(result.Pairs); n != len(expected) {
		t.Fatalf("Map has wrong number of pairs. got %d", n)
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Error("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestMapIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got %T", obj)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"1 <= 2", true},
		{"1 >= 2", false},
		{"1 <= 1", true},
		{"1 >= 1", true},
		{`"a" < "b"`, true},
		{`"a" > "b"`, false},
		{`"a" < "a"`, false},
		{`"a" > "a"`, false},
		{`"foo" == "foo"`, true},
		{`"foo" == "bar"`, false},
		{`"a" <= "b"`, true},
		{`"a" >= "b"`, false},
		{`"a" <= "a"`, true},
		{`"a" >= "a"`, true},
		{"1 + 1 == 2", true},
		{" 4 * (1 - 4) >= -13", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"null == null", true},
		{"null != null", false},
		{"null == 5", false},
		{"null != 5", true},
		{"null == fn(){}()", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got %t, want %t", result.Value, expected)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
if (10 > 1) {
    if (10 > 1) {
        return 10;
    }

    return 1;
}           `,
			10,
		},
		{
			`
let func = fn() { return 5; };
func();
return 10;
0;
            `,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unkown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unkown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unkown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unkown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
if (10 > 1) {
    if (10 > 1) {
        return true + false;
    }

    return 1;
}`,
			"unkown operator: BOOLEAN + BOOLEAN",
		},
		{
			`"Hello" - "World"`,
			"unkown operator: STRING - STRING",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			"let func = fn(x) {}; func()",
			"wrong number of arguments: expected 1, got 0",
		},
		{
			"let func = fn() {}; func(0)",
			"wrong number of arguments: expected 0, got 1",
		},
		{
			`{fn(x) { x }: "bar"}`,
			"unusable as hash key: FUNCTION",
		},
		{
			`{"foo": "bar"}[fn(x) { x }];`,
			"unusable as hash key: FUNCTION",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errorObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned, got %T (%+v)", evaluated, evaluated)
			continue
		}

		if errorObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected %q, got %q", tt.expectedMessage, errorObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not a Function. got %T (%+v)", evaluated, evaluated)
	}

	if n := len(fn.Parameters); n != 1 {
		t.Fatalf("function has wrong parameters. got %d (%+v)", n, fn.Parameters)
	}

	if p := fn.Parameters[0].String(); p != "x" {
		t.Fatalf("parameter is not 'x'. got %q", p)
	}

	expectedBody := "(x + 2)"

	if b := fn.Body.String(); b != expectedBody {
		t.Fatalf("body is not %q. got %q", expectedBody, b)
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { return x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { return x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { return x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
    fn(y) { x + y };
}

let addTwo = newAdder(2);
addTwo(2);`

	testIntegerObject(t, testEval(input), 4)
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"print()", nil},
		{"print(1)", nil},
		{"print(1, 2, 3)", nil},

		{"typeOf(1)", "INTEGER"},
		{`typeOf("")`, "STRING"},
		{"typeOf(true)", "BOOLEAN"},
		{"typeOf(null)", "NULL"},
		{"typeOf([])", "ARRAY"},
		{"typeOf({})", "MAP"},
		{"typeOf()", errors.New("wrong number of arguments: expected 1, got 0")},

		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{"len([])", 0},
		{"len([1, 2, 3])", 3},
		{`len(1)`, errors.New("argument to `len` not supported, got INTEGER")},
		{`len()`, errors.New("wrong number of arguments: expected 1, got 0")},
		{`len("one", "two")`, errors.New("wrong number of arguments: expected 1, got 2")},

		{"first([])", nil},
		{"first([1])", 1},
		{"first([1, 2, 3])", 1},
		{"first(0)", errors.New("argument to `first` not supported, got INTEGER")},
		{"first()", errors.New("wrong number of arguments: expected 1, got 0")},
		{"first([], [])", errors.New("wrong number of arguments: expected 1, got 2")},

		{"last([])", nil},
		{"last([1])", 1},
		{"last([1, 2, 3])", 3},
		{"last(0)", errors.New("argument to `last` not supported, got INTEGER")},
		{"last()", errors.New("wrong number of arguments: expected 1, got 0")},
		{"last([], [])", errors.New("wrong number of arguments: expected 1, got 2")},

		{"tail([])", nil},
		{"tail([1])", []int{}},
		{"tail([1, 2, 3])", []int{2, 3}},
		{"tail(0)", errors.New("argument to `tail` not supported, got INTEGER")},
		{"tail()", errors.New("wrong number of arguments: expected 1, got 0")},
		{"tail([], [])", errors.New("wrong number of arguments: expected 1, got 2")},

		{"push([], 1)", []int{1}},
		{"push([1, 2], 3)", []int{1, 2, 3}},
		{"push([1, 2], 3, 4, 5)", []int{1, 2, 3, 4, 5}},
		{"push(0, 0)", errors.New("argument to `push` not supported, got INTEGER")},
		{"push()", errors.New("wrong number of arguments: expected at least 2, got 0")},
		{"push([])", errors.New("wrong number of arguments: expected at least 2, got 1")},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {

		case int:
			testIntegerObject(t, evaluated, int64(expected))

		case string:
			stringObj, ok := evaluated.(*object.String)
			if !ok {
				t.Errorf("object is not String. got %T (%+v)", evaluated, evaluated)
				continue
			}

			if stringObj.Value != expected {
				t.Errorf("wrong string value. want %q, got %q", expected, stringObj.Value)
			}

		case []int:
			arrayObj, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("object is not Array. got %T (%+v)", evaluated, evaluated)
				continue
			}

			if len(arrayObj.Elements) != len(expected) {
				t.Errorf("Array has wrong number of elements. wanted %d, got %d", len(expected), len(arrayObj.Elements))
				continue
			}

			for idx, exp := range expected {
				testIntegerObject(t, arrayObj.Elements[idx], int64(exp))
			}

		case error:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got %T (%+v)", evaluated, evaluated)
				continue
			}

			if errObj.Message != expected.Error() {
				t.Errorf("wrong error message. expected %q, got %q", expected.Error(), errObj.Message)
			}

		case nil:
			testNullObject(t, evaluated)
		}

	}
}
