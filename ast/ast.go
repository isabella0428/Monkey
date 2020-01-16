package ast

import (
	"bytes"
	"../token"
)

type Node interface {
	TokenLiteral() string
	String() string
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

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
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

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

// Implement interface Expression
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

func (id *Identifier) String() string {
	return id.Value
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

func (rs *ReturnStatement) String() string{
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

// Expression Statement
type ExpressionStatement struct {
	Token token.Token	// The first token of the expression
	Expression Expression
}

// So that we can also add Expression Statement to Statements[]
func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	var out bytes.Buffer

	if es.Expression != nil {
		out.WriteString(es.Expression.String())
	}

	return out.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string  {return il.Token.Literal}
func (il *IntegerLiteral) String()	string	     {return il.Token.Literal}

type PrefixExpression struct {
	Token token.Token 			// The prefix token, e.g. !
	Operator string
    Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {return pe.Token.Literal}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token 		token.Token
	Left 		Expression
	Operator 	string
	Right 		Expression
}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) String() string{
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" ")
	out.WriteString(ie.Operator)
	out.WriteString(" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) TokenLiteral() string{ return b.Token.Literal }
func (b *Boolean) String() string {return b.Token.Literal}
func (b *Boolean) expressionNode() {}
