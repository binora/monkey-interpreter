package parser

import (
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
