package evaluator

import (
	"fmt"

	"github.com/estevesnp/dsb/pkg/object"
)

func Len(args ...object.Object) object.Object {
	if n := len(args); n != 1 {
		return newError("wrong number of arguments: expected 1, got %d", n)
	}

	switch arg := args[0].(type) {
	case *object.String:
		length := len([]rune(arg.Value))
		return &object.Integer{Value: int64(length)}
	default:
		return notSupported("len", args[0])
	}
}

func Print(args ...object.Object) object.Object {
	arguments := make([]any, 0, len(args))

	for _, arg := range args {
		switch arg := arg.(type) {
		case *object.Integer:
			arguments = append(arguments, arg.Value)
		case *object.String:
			arguments = append(arguments, arg.Value)
		case *object.Boolean:
			arguments = append(arguments, arg.Value)
		default:
			return notSupported("print", arg)
		}
	}

	fmt.Println(arguments...)

	return NULL
}

func notSupported(name string, obj object.Object) *object.Error {
	return newError("argument to `%s` not supported, got %s", name, obj.Type())
}
