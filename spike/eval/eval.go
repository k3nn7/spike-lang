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
	case *ast.Boolean:
		return evalBoolean(node)
	case *ast.PrefixExpression:
		right, _ := Eval(node.Right)
		return evalPrefixExpression(right, node.Operator)
	}
	return nil, nil
}

func evalStatements(statements []ast.Statement) (object.Object, error) {
	for _, statement := range statements {
		return Eval(statement)
	}

	return nil, nil
}

func evalBoolean(node *ast.Boolean) (object.Object, error) {
	if node.Value {
		return &object.True, nil
	}

	return &object.False, nil
}

func evalPrefixExpression(right object.Object, operator string) (object.Object, error) {
	switch operator {
	case "!":
		return evalBangOperator(right)
	default:
		return nil, nil
	}
}

func evalBangOperator(right object.Object) (object.Object, error) {
	switch right {
	case &object.True:
		return &object.False, nil
	case &object.False:
		return &object.True, nil
	default:
		return nil, nil
	}
}
