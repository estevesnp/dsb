package evaluator

import (
	"fmt"

	"github.com/estevesnp/dsb/pkg/ast"
	"github.com/estevesnp/dsb/pkg/object"
	"github.com/estevesnp/dsb/pkg/token"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{Node: node}
}

func evalUnquoteCalls(quoted ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}

		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if len(call.Arguments) != 1 {
			return node
		}

		unquoted := Eval(call.Arguments[0], env)
		return convertObjectToASTNode(unquoted)
	})
}

func isUnquoteCall(node ast.Node) bool {
	callExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}

	return callExpression.Function.TokenLiteral() == "unquote"
}

func convertObjectToASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {

	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}

	case *object.Boolean:
		var t token.Token
		if obj.Value {
			t = token.Token{Type: token.TRUE, Literal: "true"}
		} else {
			t = token.Token{Type: token.FALSE, Literal: "false"}
		}
		return &ast.Boolean{Token: t, Value: obj.Value}

	case *object.String:
		t := token.Token{
			Type:    token.STRING,
			Literal: obj.Value,
		}
		return &ast.StringLiteral{Token: t, Value: obj.Value}

	case *object.Array:
		t := token.Token{
			Type:    token.LBRACKET,
			Literal: "[",
		}

		elems := make([]ast.Expression, 0, len(obj.Elements))

		for _, elem := range obj.Elements {
			n := convertObjectToASTNode(elem)
			exp, ok := n.(ast.Expression)
			if !ok {
				continue
			}
			elems = append(elems, exp)
		}

		return &ast.ArrayLiteral{Token: t, Elements: elems}

	case *object.Map:
		t := token.Token{
			Type:    token.LBRACE,
			Literal: "{",
		}

		pairs := make(map[ast.Expression]ast.Expression, len(obj.Pairs))
		for _, value := range obj.Pairs {
			keyNode := convertObjectToASTNode(value.Key)
			valueNode := convertObjectToASTNode(value.Value)

			keyExpression, ok := keyNode.(ast.Expression)
			if !ok {
				continue
			}

			valueExpression, ok := valueNode.(ast.Expression)
			if !ok {
				continue
			}

			pairs[keyExpression] = valueExpression
		}

		return &ast.MapLiteral{Token: t, Pairs: pairs}

	case *object.Function:
		t := token.Token{
			Type:    token.FUNCTION,
			Literal: "fn",
		}
		return &ast.FunctionLiteral{
			Token:      t,
			Parameters: obj.Parameters,
			Body:       obj.Body,
		}

	case *object.Quote:
		return obj.Node

	default:
		t := token.Token{
			Type:    token.NULL,
			Literal: "null",
		}
		return &ast.NullLiteral{Token: t}
	}
}
