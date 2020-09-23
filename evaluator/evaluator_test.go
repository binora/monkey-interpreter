package evaluator

import (
	"interpreters/lexer"
	"interpreters/object"
	"interpreters/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			input:    "5",
			expected: 5,
		},
		{
			input:    "10",
			expected: 10,
		},
		{
			input:    "-5",
			expected: -5,
		},
		{
			input:    "-10",
			expected: -10,
		},
		{
			input:    "5 + 5 + 5 + 5 - 10",
			expected: 10,
		},
		{
			input:    "5 * 2 + 10 ",
			expected: 20,
		},
		{
			input:    "50 / 2 * 2 + 10",
			expected: 60,
		},
		{
			input:    "2 * ( 5 + 10)",
			expected: 30,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			input:    "true",
			expected: true,
		},
		{
			input:    "false",
			expected: false,
		},
		{
			input:    "true == true",
			expected: true,
		},
		{
			input:    "(1 < 2) == true",
			expected: true,
		},
		{
			input:    "(1 < 2) == false",
			expected: false,
		},
		{
			input:    "(1 > 2) == true",
			expected: false,
		},
		{
			input:    "(1 > 2) == false",
			expected: true,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!!true", true},
		{"!5", false},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			input:    "if (true) { 10 }",
			expected: 10,
		},
		{
			input:    "if (false) { 10 }",
			expected: nil,
		},
		{
			input:    "if (1) { 10 }",
			expected: 10,
		},
		{
			// evaluates to true
			input:    "if (0) { 10 }",
			expected: 10,
		},
		{
			input:    "if (1 > 2) { 10 } else { 20 }",
			expected: 20,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
			continue
		}
		testNullObject(t, evaluated)
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("expected null object, got: %T (%v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			input:    "return 10;",
			expected: 10,
		},
		{
			input:    "return 10*10",
			expected: 100,
		},
		{
			input:    "9; return 2*5; 9",
			expected: 10,
		},
		{
			input: `
if (10 > 1) {
	if (10 > 1) {
		return 10;
	}
	return 1;
}`,
			expected: 10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "5  + true",
			expected: "type mismatch: INTEGER + BOOLEAN",
		},
		{
			input:    "-true",
			expected: "unknown operator: -BOOLEAN",
		},
		{
			input:    "if ( 10 > 1) { true + false; }",
			expected: "unknown operator: BOOLEAN + BOOLEAN",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("No error returned but we need errors!, got: %v, want: %v", evaluated, tt.expected)
			continue
		}
		if tt.expected != errObj.Message {
			t.Errorf("error messages are different. got:%s, want: %s", errObj.Message, tt.expected)
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("obj is not object.Integer, got: %T", obj)
		return false
	}

	if result.Type() != object.INTEGER_OBJ {
		t.Errorf("result.Type() is not object.INTEGER_OBJ, got: %+v", result.Type())
		return false
	}

	if result.Value != expected {
		t.Errorf("result.Value is not %d, got: %d", expected, result.Value)
		return false
	}

	return true

}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("obj is not object.Boolean, got: %T", obj)
		return false
	}

	if result.Type() != object.BOOLEAN_OBJ {
		t.Errorf("result.Type() is not object.BOOLEAN_OBJ, got: %+v", result.Type())
		return false
	}

	if result.Value != expected {
		t.Errorf("result.Value is not %t, got: %t", expected, result.Value)
		return false
	}

	return true
}