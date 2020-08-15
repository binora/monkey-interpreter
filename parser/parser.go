package parser

import (
	"fmt"
	"interpreters/ast"
	"interpreters/lexer"
	"interpreters/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// We call nextToken twice to set curToken and peekToken.
	// We need to do this because we need to know peekToken ( the next token to be parsed) to make decisions while
	// parsing curToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

// parseLetStatement as of now parses simpler 'let' statements
// e.g. let x = 5;
func (p *Parser) parseLetStatement() *ast.LetStatement {
	// curToken is "let"
	letStatement := &ast.LetStatement{Token: p.curToken}

	// return nil if next token is not an identifier
	// otherwise set curToken to identifier i.e. "x"
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// create identifier and assign to letStatement.Name
	letStatement.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// next token should be assign "="
	// if yes, then move to next token
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// skip all expressions
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return letStatement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStatement := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return returnStatement
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.Type) {
	message := fmt.Sprintf("expected next token to be %s, got: %s", t, p.peekToken.Type)
	p.errors = append(p.errors, message)
}
