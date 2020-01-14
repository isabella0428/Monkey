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
	EQ     	 = "=="
    NOT_EQ 	 = "!="

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
	IF = "IF"
	TRUE = "TRUE"
	FALSE = "FALSE"
	ELSE = "ELSE"
	RETURN = "RETURN"
)

var keywords = map[string] TokenType{
	"fn" : FUNCTION,
	"let" : LET,
	"true": TRUE,
	"false": FALSE,
	"else" : ELSE,
	"return" : RETURN,
	"if" : IF,
}

func LookUpIndent(indent string) TokenType {
	if tok, ok := keywords[indent]; ok {
		return tok
	}
	return IDENT
}
