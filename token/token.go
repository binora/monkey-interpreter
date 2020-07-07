package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// IDENT stands for Identifier type
	// E.g. foobar
	IDENT = "IDENT"

	// INT stands for Integer type
	INT = "INT"

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "fn"
	LET      = "LET"
)

// TokenType is a string and serves the
// purpose of differentiating between various token types
type Type string

// Token is created for each character of the input string parsed by the lexer
type Token struct {
	Type    Type
	Literal string
}

var keywords = map[string]Type{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
