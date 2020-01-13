package token

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"		// a token that we don't know
	EOF = "EOF"				// end of file . 

	// identifiers + literals
	IDENT = "IDENT"		// variables, function names
	INT = "INT"			// 123456

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT = "<"
	GT = ">"

	// delimiter
	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// keywords
	FUNCTION = "FUNCTION"
	LET = "LET"
)

var keywords = map[string] TokenType{
	"fn" : FUNCTION,
	"let" : LET,
}

func LookUpIndent(indent string) TokenType {
	if tok, ok := keywords[indent]; ok {
		return tok
	}
	return IDENT
}
