package evaluator

import (
	"github.com/estevesnp/dsb/pkg/ast"
	"github.com/estevesnp/dsb/pkg/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}

	for i, statement := range program.Statements {
		if isMacroDefinition(statement) {
			addMacro(statement, env)
			definitions = append(definitions, i)
		}
	}

	for i := len(definitions) - 1; i >= 0; i -= 1 {
		definitionIdx := definitions[i]
		program.Statements = append(
			program.Statements[:definitionIdx],
			program.Statements[definitionIdx+1:]...)
	}
}

func isMacroDefinition(node ast.Statement) bool {
	letSatement, ok := node.(*ast.LetStatement)
	if !ok {
		return false
	}

	_, ok = letSatement.Value.(*ast.MacroLiteral)
	return ok
}

func addMacro(stmt ast.Statement, env *object.Environment) {
	letStatement, _ := stmt.(*ast.LetStatement)
	macroLiteral, _ := letStatement.Value.(*ast.MacroLiteral)

	macro := &object.Macro{
		Parameters: macroLiteral.Parameters,
		Env:        env,
		Body:       macroLiteral.Body,
	}

	env.Set(letStatement.Name.Value, macro)
}
