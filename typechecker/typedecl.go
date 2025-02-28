package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
	"walrus/utils"
)

func checkTypeDeclaration(node ast.TypeDeclStmt, env *TypeEnvironment) ValueTypeInterface {

	typeName := node.UDType

	fmt.Printf("declaring type %s\n", node.UDTypeName)

	var val ValueTypeInterface

	switch t := typeName.(type) {
	case ast.StructType:
		val = checkStructTypeDecl(node.UDTypeName, t, env)
	case ast.InterfaceType:
		val = checkInterfaceTypeDecl(node.UDTypeName, t, env)
	default:
		val = evaluateTypeName(typeName, env)
	}

	typeVal := UserDefined{
		DataType: USER_DEFINED_TYPE,
		TypeName: node.UDTypeName,
		TypeDef:  val,
	}

	err := DeclareType(node.UDTypeName, typeVal)
	if err != nil {
		fmt.Printf("cannot declare type %s: %s\n", node.UDTypeName, err.Error())
		errgen.AddError(env.filePath, node.Start.Line, node.End.Line, node.Start.Column, node.End.Column, err.Error())
	}

	utils.CYAN.Print("Type ")
	utils.PURPLE.Print(node.UDTypeName)
	fmt.Println(" declared")

	return val
}
