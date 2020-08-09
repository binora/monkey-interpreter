package ast

import "interpreters/token"

/*
In the monkey language, We differentiate between a statement and an expression
An expression is any that returns a value e.g. 5, 1+1, function declaration
A statement doesn't return a value e.g. let a = 5;
*/

type Node interface {
	TokenLiteral() string
}

// Statement interface wraps the Node interface and also has a dummy StatementNode method to
// distinguish between expression and statement
type Statement interface {
	Node
	StatementNode()
}

// Expression interface also wraps the Node interface and has the dummy ExpressionNode() method
type Expression interface {
	Node
	ExpressionNode()
}

// Program is really just a bunch of statements where each statement is a node
// in the AST we are trying to build
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// LetStatement is the ast node which holds various fields to represent different
// parts of the 'let' statement.
// e.g. let x = 5.
// "let" is represented by token.Token
// "x" is the Identifier struct ( We consider it as an Expression node for convenience...something I don't understand)
// 5 is represented by the Expression node which makes it possible to represent anything after "=" that returns a value
type LetStatement struct {
	// "x" token.IDENT
	Name *Identifier
	// "let" token.LET
	Token token.Token
	Value Expression
}

func (l *LetStatement) StatementNode() {}

func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) ExpressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
