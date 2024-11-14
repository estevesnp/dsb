package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/estevesnp/dsb/pkg/ast"
)

type ObjectType string

const (
	NULL_OBJ         = "NULL"
	ERROR_OBJ        = "ERROR"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	BUILTIN_OBJ      = "BUILTIN"
	FUNCTION_OBJ     = "FUNCTION"
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	ARRAY_OBJ        = "ARRAY"
	MAP_OBJ          = "MAP"
	QUOTE_OBJ        = "QUOTE"
	MACRO_OBJ        = "MACRO"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}

var stringHashCache = map[string]HashKey{}

type BuiltinFunction func(args ...Object) Object

// Null
type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

func (n *Null) Inspect() string {
	return "null"
}

// Error
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("ERROR: %s", e.Message)
}

// Integer
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// Boolean
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

// String
type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) HashKey() HashKey {
	if key, ok := stringHashCache[s.Value]; ok {
		return key
	}

	h := fnv.New64a()
	h.Write([]byte(s.Value))

	key := HashKey{Type: s.Type(), Value: h.Sum64()}
	stringHashCache[s.Value] = key

	return key
}

// Array
type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType {
	return ARRAY_OBJ
}

func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := make([]string, len(ao.Elements))
	for idx, el := range ao.Elements {
		elements[idx] = el.Inspect()
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// Map
type Map struct {
	Pairs map[HashKey]HashPair
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

func (m *Map) Type() ObjectType {
	return MAP_OBJ
}

func (m *Map) Inspect() string {
	var out bytes.Buffer

	pairs := make([]string, 0, len(m.Pairs))
	for _, pair := range m.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// ReturnValue
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

// Function
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := make([]string, 0, len(f.Parameters))
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// Builtin
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}

// Quote
type Quote struct {
	Node ast.Node
}

func (q *Quote) Type() ObjectType {
	return QUOTE_OBJ
}

func (q *Quote) Inspect() string {
	return fmt.Sprintf("QUOTE(%s)", q.Node.String())
}

// Macro
type Macro struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (m *Macro) Type() ObjectType {
	return MACRO_OBJ
}

func (m *Macro) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range m.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("macro(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(m.Body.String())
	out.WriteString("\n}")

	return out.String()
}
