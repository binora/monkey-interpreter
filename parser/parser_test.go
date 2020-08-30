package parser

import (
	"fmt"
	"interpreters/ast"
	"interpreters/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	l := lexer.New(input)
	p := New(l)

	program := p.parseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got: %d", len(program.Statements))
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

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("It's supposed to be let but got: %s", statement.TokenLiteral())
		return false
	}

	letStmt, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("unable to cast statement to *ast.LetStatement, got %T", statement)
		return false
	}

	// test for identifier
	if letStmt.Name.Value != name {
		t.Errorf("letStatement identifier value different from expected, got: %s, expected: %s", letStmt.Name.Value, name)
		return false
	}

	// Why are we doing this check again ?
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStatement identifier tokenLiteral different from expected, got: %s, expected: %s", letStmt.Name.TokenLiteral(), name)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 2134234;
`

	l := lexer.New(input)
	p := New(l)

	program := p.parseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("expected statements to be 3, got: %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStatement, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("could cast statement to ReturnStatement, got: %T", stmt)
			continue
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("expected literal to be return, but got: %s", returnStatement.TokenLiteral())
		}

	}

}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.parseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statement, got: %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an expression statement, got: %T", program.Statements[0])
	}

	identifier, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("program.Statemet[0].Expression is not an identifier, got: %T", stmt.Expression)
	}

	if identifier.TokenLiteral() != "foobar" {
		t.Fatalf("identifier.TokenLiteral() != %s, got: %s", "foobar", identifier.TokenLiteral())
	}

	if identifier.Value != "foobar" {
		t.Fatalf("identifier.Value != %s, got: %s", "foobar", identifier.Value)
	}

}

func TestBooleanExpression(t *testing.T) {
	input := "false"
	l := lexer.New(input)
	p := New(l)

	program := p.parseProgram()
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statement, got: %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an expression statement, got: %T", program.Statements[0])
	}

	boolean, ok := stmt.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.Boolean, got: %T", stmt.Expression)
	}

	if boolean.TokenLiteral() != "false" {
		t.Fatalf("boolean.TokenLiteral() is not %s, got: %s", "false", boolean.TokenLiteral())
	}

	if boolean.Value != false {
		t.Fatalf("boolean.Value is not %t, got: %t", false, boolean.Value)
	}
}

func TestIntegerLiterals(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.parseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statement, got: %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an expression statement, got: %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("program.Statemet[0].Expression is not an IntegerLiteral, got: %T", stmt.Expression)
	}

	if literal.TokenLiteral() != "5" {
		t.Fatalf("literal.TokenLiteral() != %s, got: %s", "5", literal.TokenLiteral())
	}

	if literal.Value != 5 {
		t.Fatalf("literal.Value != %d, got: %d", 5, literal.Value)
	}

}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true", "!", true},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.parseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("unexpected number of program statements, got: %d, want: %d ", len(program.Statements), 1)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("unable to cast statement to Expression type, got: %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("unable to cast statement to prefix expression, got: %T", stmt)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s, got: %s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}

	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not *ast.IntegerLiteral, got: %T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integer value not equal to %d, got: %d", value, integ.Value)
		return false
	}

	if il.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("il.TokenLiteral() not %d, got: %s", value, il.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp is not an *ast.Identifier, got: %T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value is not %s, got: %s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() is not %s, got: %s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp is not *ast.Boolean, got: %T", exp)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value is not %t, got: %t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral() is not %t, got: %s", value, boolean.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
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
	t.Errorf("type of exp not handled, got: %T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	infixExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not *ast.Expression, got: %T", exp)
		return false
	}

	if !testLiteralExpression(t, infixExp.Left, left) {
		return false
	}

	if infixExp.Operator != operator {
		t.Errorf("infixExp.Operator is not %s, got:%s", operator, infixExp.Operator)
		return false
	}

	if !testLiteralExpression(t, infixExp.Right, right) {
		return false
	}

	return true
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.parseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("unexpected number of program statements, got: %d, want: %d", len(program.Statements), 1)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not an *ast.ExpressionStatement, got: %T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			t.Fatalf("unable to parse stmt.Expression")
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
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
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
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"true",
			"true",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
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
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.parseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, actual=%q", tt.expected, actual)
		}

	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)

	program := p.parseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements is not %d, got %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an if expression, got: %T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not an *ast.IfExpression, got: %T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		t.Fatalf("exp.Condition is not an InfixExpression, got: %T", exp.Condition)
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("exp.Consequence.Statements is not %d, got: %d", 1, len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not an ast.ExpressionStatement , got: %T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		t.Fatalf("consequence.Expression is not an identifier")
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Alternative.Statements[0] is not an ast.ExpressionStatement, got: %T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		t.Fatalf("alternative.Expression is not an identifier, got: %+v", alternative.Expression)
	}

}

func checkParserErrors(t *testing.T, p *Parser) {
	e := p.Errors()
	if len(e) == 0 {
		return
	}
	t.Errorf("parser had %d errors", len(e))
	for _, message := range e {
		t.Errorf("parser error: %q", message)
	}

	t.FailNow()
}
