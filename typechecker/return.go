package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkReturnStmt(returnNode ast.ReturnStmt, env *TypeEnvironment) ValueTypeInterface {
	//check if the function is declared
	if env.scopeType != FUNCTION_SCOPE {

		errgen.AddError(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, "Return statement must be inside a function")
	}

	//check if the return type matches the function return type
	returnType := nodeType(returnNode.Value, env)

	fnReturns := getFunctionReturnValue(env, returnNode)

	err := matchTypes(fnReturns, returnType)
	if err != nil {

		errgen.AddError(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, fmt.Sprintf("cannot return '%s' from this scope. function '%s' expects return type '%s'", valueTypeInterfaceToString(returnType), env.scopeName, valueTypeInterfaceToString(fnReturns)))
	}

	return ReturnType{
		DataType:   RETURN_TYPE,
		Expression: returnType,
	}
}
