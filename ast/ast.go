package ast

import (
	"../token"
)

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program node is the root node of each AST Tree
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// LetStatement ast Tree building
type LetStatement struct {
	Name *Identifier
	Token token.Token
	Value Expression
}

// Implement Statement Interface
func (ls *LetStatement) statementNode() {}

// Implement Node Interface
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// Let Identifier to implement interface Expression
// Although identifiers don't always return a value,
// we can save conditions in this way
type Identifier struct {
	Token token.Token
	Value string
}

func (id *Identifier) expressionNode() {}
func (id *Identifier) TokenLiteral() string {
	return id.Token.Literal
}

// Return Statement
type ReturnStatement struct {
	Token token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

