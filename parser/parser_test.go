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
