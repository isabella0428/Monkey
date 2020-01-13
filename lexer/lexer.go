package lexer

import "../token"

type Lexer struct {
	input 			string
	position 		int 		// current position in input(points to current char)
	readPosition 	int 		// current reading position in input(after current char)
	ch 				byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Read the character and move to the next one
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'	// Monkey only supports integers
}

// Convert current character into token and move to the next one
func (l *Lexer) nextToken() token.Token {
	var tok token.Token
	l.skipWhiteSpace()

	switch l.ch {
		case '=':
			tok = newToken(token.ASSIGN, l.ch)
		case '+':
			tok = newToken(token.PLUS, l.ch)
		case '-':
			tok = newToken(token.MINUS, l.ch)
		case '!':
			tok = newToken(token.BANG, l.ch)
		case '*':
			tok = newToken(token.ASTERISK, l.ch)
		case '/':
			tok = newToken(token.SLASH, l.ch)
		case '<':
			tok = newToken(token.LT, l.ch)
		case '>':
			tok = newToken(token.GT, l.ch)
		case ';':
			tok = newToken(token.SEMICOLON, l.ch)
		case '(':
			tok = newToken(token.LPAREN, l.ch)
		case ')':
			tok = newToken(token.RPAREN, l.ch)
		case '{':
			tok = newToken(token.LBRACE, l.ch)
		case '}':
			tok = newToken(token.RBRACE, l.ch)
		case ',':
			tok = newToken(token.COMMA, l.ch)
		case 0:
			tok.Literal = ""
			tok.Type = token.EOF
		default:
			if isLetter(l.ch) {
				tok.Literal = l.readIdentifier()
				tok.Type = token.LookUpIndent(tok.Literal)
				return tok
			} else if isDigit(l.ch) {
				tok.Literal = l.readNumber()
				// fmt.Printf("%q", tok.Literal)
				tok.Type = token.INT
				return tok
			} else {
				tok = newToken(token.ILLEGAL, l.ch)
			}
	}
	l.readChar()
	return tok
}


func (l *Lexer)skipWhiteSpace() {
	for (l.ch == '\t') || (l.ch == '\r') || (l.ch == '\n') || (l.ch == ' ') {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token{
	return token.Token{Type:tokenType, Literal:string(ch)}
}
