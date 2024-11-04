package evaluator

import (
	"fmt"

	"github.com/estevesnp/dsb/pkg/object"
)

func Print(args ...object.Object) object.Object {
	arguments := make([]any, len(args))

	for idx, arg := range args {
		arguments[idx] = arg.Inspect()
	}

	fmt.Println(arguments...)

	return NULL
}

func TypeOf(args ...object.Object) object.Object {
	if err := validateLength(1, args); err != nil {
		return err
	}
	return &object.String{Value: string(args[0].Type())}
}

func Len(args ...object.Object) object.Object {
	if err := validateLength(1, args); err != nil {
		return err
	}

	switch arg := args[0].(type) {
	case *object.String:
		length := len([]rune(arg.Value))
		return &object.Integer{Value: int64(length)}
	case *object.Array:
		length := len(arg.Elements)
		return &object.Integer{Value: int64(length)}
	default:
		return notSupported("len", args[0])
	}
}

func First(args ...object.Object) object.Object {
	if err := validateLength(1, args); err != nil {
		return err
	}

	arr, ok := args[0].(*object.Array)
	if !ok {
		return notSupported("first", args[0])
	}

	if len(arr.Elements) == 0 {
		return NULL
	}

	return arr.Elements[0]
}

func Last(args ...object.Object) object.Object {
	if err := validateLength(1, args); err != nil {
		return err
	}

	arr, ok := args[0].(*object.Array)
	if !ok {
		return notSupported("last", args[0])
	}

	length := len(arr.Elements)

	if length == 0 {
		return NULL
	}

	return arr.Elements[length-1]
}

func Tail(args ...object.Object) object.Object {
	if err := validateLength(1, args); err != nil {
		return err
	}

	arr, ok := args[0].(*object.Array)
	if !ok {
		return notSupported("tail", args[0])
	}

	length := len(arr.Elements)

	if length == 0 {
		return NULL
	}

	newElems := make([]object.Object, length-1)
	copy(newElems, arr.Elements[1:length])

	return &object.Array{Elements: newElems}
}

func Push(args ...object.Object) object.Object {
	if n := len(args); n < 2 {
		return newError("wrong number of arguments: expected at least 2, got %d", n)
	}

	arr, ok := args[0].(*object.Array)
	if !ok {
		return notSupported("push", args[0])
	}

	length := len(arr.Elements)

	newElems := make([]object.Object, length, length+len(args)-1)
	copy(newElems, arr.Elements)

	newElems = append(newElems, args[1:]...)

	return &object.Array{Elements: newElems}
}

func notSupported(name string, obj object.Object) *object.Error {
	return newError("argument to `%s` not supported, got %s", name, obj.Type())
}

func validateLength(length int, args []object.Object) *object.Error {
	if n := len(args); n != length {
		return newError("wrong number of arguments: expected 1, got %d", n)
	}
	return nil
}
