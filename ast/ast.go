package ast

import (
	"bytes"
	"interpreters/token"
	"strings"
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

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) ExpressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) ExpressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) ExpressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type IfExpression struct {
	// if token
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) ExpressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) StatementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, statement := range bs.Statements {
		out.WriteString(statement.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) ExpressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(fl.TokenLiteral())
	out.WriteString(" (")

	var params []string
	for _, param := range fl.Parameters {
		params = append(params, param.String())
	}
	out.WriteString(strings.Join(params, ","))

	out.WriteString(" ) ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) ExpressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Value }

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) ExpressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ce.Function.String())

	var args []string
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString("(")
	out.WriteString(strings.Join(args, ","))
	out.WriteString(")")

	return out.String()
}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) ExpressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("[")
	var expressions []string
	for _, exp := range al.Elements {
		expressions = append(expressions, exp.String())
	}
	out.WriteString(strings.Join(expressions, ","))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) ExpressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")
	out.WriteString(")")

	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (h *HashLiteral) ExpressionNode()      {}
func (h *HashLiteral) TokenLiteral() string { return h.Token.Literal }
func (h *HashLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	pairs := []string{}
	for key, value := range h.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString(strings.Join(pairs, ","))
	out.WriteString("}")
	return out.String()
}
