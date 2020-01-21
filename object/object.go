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
	BUILTIN_OBJ = "BUILTIN"
	ARRAY_OBJ = "ARRAY"
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

func (s *String) Type() ObjectType {return STRING_OBJ}
func (s *String) Inspect() string {return s.Value}

type Builtin struct {
	Fn BuiltinFunction
}

func (bi *Builtin) Type() ObjectType {return BUILTIN_OBJ}
func (bi *Builtin) Inspect() string {return "builtin function"}

type BuiltinFunction func(args ...Object) Object

type Array struct {
	Elements	[]Object
}

func (a *Array) Type() ObjectType {return ARRAY_OBJ}
func (a *Array) Inspect()	string {
	var out bytes.Buffer

	elements := []string{}
	for _, tt := range a.Elements {
		elements = append(elements, tt.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}