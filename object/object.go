package object


import (
	"fmt"
	"bytes"
	"strings"
	"hash/fnv"
	"../ast"
)

type ObjectType string

type Object interface {
	Type() 		ObjectType
	Inspect()	string
}

type Hashable interface {
	HashKey() HashKey
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
	HASH_OBJ = "HASH"
	HASHKEY_OBJ = "HASHKEY"
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

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType {return HASH_OBJ}
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs{
		pairs = append(pairs, fmt.Sprintf("%s: %s",
               pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type HashKey struct {
	ObjectType 	ObjectType
	Value 		uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value{
		value = 1
	} else {
		value = 0
	}
	return HashKey{ObjectType:b.Type(), Value:value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{ObjectType:i.Type(), Value:uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{ObjectType: s.Type(), Value: h.Sum64()}
}

func (hk HashKey) Type() ObjectType{return hk.ObjectType}
func (hk HashKey) Inspect() string {return string(hk.Value)}

type HashPair struct {
	Key		Object
	Value	Object
}
