package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkVariableAssignment(node ast.VarAssignmentExpr, env *TypeEnvironment) ValueTypeInterface {

	Assignee := node.Assignee
	valueToAssign := node.Value

	//varToAssign := node.Identifier
	expected := CheckAST(Assignee, env)
	provided := CheckAST(valueToAssign, env)

	MatchTypes(expected, provided, env.filePath, valueToAssign.StartPos().Line, valueToAssign.StartPos().Column, valueToAssign.EndPos().Column)

	var varName string

	if !IsLValue(Assignee) {
		errgen.MakeError(env.filePath, Assignee.StartPos().Line, Assignee.StartPos().Column, Assignee.EndPos().Column, "invalid assignment expression. the assignee must be a lvalue").AddHint("lvalue is something that has a memory address\nFor example: variables.\nso you cannot assign values something which does not exist in memory as an independent identifier.", errgen.TEXT_HINT).Display()
	}

	switch assignee := Assignee.(type) {
	case ast.IdentifierExpr:
		varName = assignee.Name
	case ast.ArrayIndexAccess:
		return nil
	case ast.PropertyExpr:
		return nil
	default:
		panic("cannot assign to this type")
	}

	//get the stored type
	scope, err := env.ResolveVar(varName)
	if err != nil {
		errgen.MakeError(env.filePath, valueToAssign.StartPos().Line, valueToAssign.StartPos().Column, valueToAssign.EndPos().Column, err.Error()).Display()
	}

	//if constant
	if scope.constants[varName] {
		errgen.MakeError(env.filePath, valueToAssign.StartPos().Line, valueToAssign.StartPos().Column, valueToAssign.EndPos().Column, fmt.Sprintf("'%s' is constant", varName)).AddHint("cannot assign value to constant variables", errgen.TEXT_HINT).Display()
	}
	scope.variables[varName] = provided
	return nil
}

func checkVariableDeclaration(node ast.VarDeclStmt, env *TypeEnvironment) ValueTypeInterface {

	varToDecl := node.Variable

	var expectedTypeInterface ValueTypeInterface

	if node.ExplicitType != nil {
		expectedTypeInterface = handleExplicitType(node, env)
	} else {
		typ := CheckAST(node.Value, env)
		expectedTypeInterface = typ
	}

	if node.IsAssigned && node.ExplicitType != nil {
		providedValue := CheckAST(node.Value, env)
		MatchTypes(expectedTypeInterface, providedValue, env.filePath, node.Value.StartPos().Line, node.Value.StartPos().Column, node.Value.EndPos().Column)
	}

	err := env.DeclareVar(varToDecl.Name, expectedTypeInterface, node.IsConst)
	if err != nil {
		errgen.MakeError(env.filePath, node.Variable.StartPos().Line, node.Variable.StartPos().Column, node.Variable.EndPos().Column, err.Error()).Display()
	}
	return nil
}
