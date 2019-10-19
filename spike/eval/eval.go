package eval

import (
	"spike-interpreter-go/spike/eval/object"
	"spike-interpreter-go/spike/parser/ast"
)

func Eval(node ast.Node) (object.Object, error) {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.Integer:
		return &object.Integer{Value: node.Value}, nil
	}
	return nil, nil
}

func evalStatements(statements []ast.Statement) (object.Object, error) {
	for _, statement := range statements {
		return Eval(statement)
	}

	return nil, nil
}
