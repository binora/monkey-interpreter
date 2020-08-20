package ast

import (
	"bytes"
	"interpreters/token"
)

/*
In the monkey language, We differentiate between a statement and an expression
An expression is anything that returns a value e.g. 5, 1+1, function declaration
A statement doesn't return a value e.g. let a = 5;
*/

type Node interface {
	TokenLiteral() string
	String() string
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

func (p *Program) String() string {
	var out bytes.Buffer
	for _, statement := range p.Statements {
		out.WriteString(statement.String())
	}
	return out.String()
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

func (l *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String())
	out.WriteString(" = ")

	// Until we know how to evaluate expressions
	if l.Value != nil {
		out.WriteString(l.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) ExpressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

// ReturnStatement represents the 'return' statement in the monkey language
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (r *ReturnStatement) StatementNode() {}

func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral() + " ")

	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement represents an expression
// e.g. x + 10; is a valid expression as well as a statement accepted in the monkey language
type ExpressionStatement struct {
	// first token of the expression
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStatement) StatementNode() {}
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) ExpressionNode() {}
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}
