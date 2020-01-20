package object


import (
	"fmt"
	"bytes"
	"strings"
	"../ast"
)

type ObjectType string

type Object interface {
	Type() 		ObjectType
	Inspect()	string
}

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ 	= "NULL" 
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ = "ERROR"
	FUNCTION_OBJ  = "FUNCTION"
	STRING_OBJ	 = "STRING"
)

// -----------------------
// Type Definitions
// ----------------------
type Integer struct {
	Value int64
}
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value)}
func (i *Integer) Type() ObjectType {return INTEGER_OBJ}

type Boolean struct{
	Value bool
}
func (b *Boolean) Inspect() string {return fmt.Sprintf("%t", b.Value)}
func (b *Boolean) Type() ObjectType {return BOOLEAN_OBJ}

type  Null struct {}
func (n *Null) Inspect() string {return "null"}
func (n *Null) Type() ObjectType {return NULL_OBJ}

type ReturnValue struct {
	Value Object
}
func (rv *ReturnValue) Inspect() string {return rv.Value.Inspect()}
func (rv *ReturnValue) Type()	ObjectType {return RETURN_VALUE_OBJ}

type Error struct {
	Message string
}

func (e *Error) Inspect() string {return "ERROR: " + e.Message}
func (e *Error) Type()	  ObjectType {return ERROR_OBJ}

type Function struct {
	Parameters []*ast.Identifier
	Body	   *ast.BlockStatement
	Env		   *Environment				// Contains local parameters inside the function
}

func (fn *Function) Type() ObjectType {return FUNCTION_OBJ}
func (fn *Function) Inspect() string  {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fn.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(fn.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType{return STRING_OBJ}
func (s *String) Inspect() string {return s.Value}