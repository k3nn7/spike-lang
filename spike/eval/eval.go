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
	case *ast.InfixExpression:
		left, _ := Eval(node.Left)
		right, _ := Eval(node.Right)

		return evalInfixExpression(left, right, node.Operator)
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
	case "-":
		return evalMinusOperator(right)
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

func evalMinusOperator(right object.Object) (object.Object, error) {
	switch rightObject := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -rightObject.Value}, nil
	default:
		return nil, nil
	}
}

func evalInfixExpression(left, right object.Object, operator string) (object.Object, error) {
	switch operator {
	case "+":
		return evalPlusInfixOperator(left, right)
	case "-":
		return evalMinusInfixOperator(left, right)
	case "*":
		return evalAsteriskInfixOperator(left, right)
	case "/":
		return evalAsteriskSlashOperator(left, right)
	case "==":
		equal, err := left.Equal(right)
		return nativeBoolToBoolean(equal), err
	case "!=":
		equal, err := left.Equal(right)
		return nativeBoolToBoolean(!equal), err
	case "<":
		leftComparable := left.(object.Comparable)
		rightComparable := right.(object.Comparable)
		result, err := leftComparable.Compare(rightComparable)
		return nativeBoolToBoolean(result == object.LT), err
	case ">":
		leftComparable := left.(object.Comparable)
		rightComparable := right.(object.Comparable)
		result, err := leftComparable.Compare(rightComparable)
		return nativeBoolToBoolean(result == object.GT), err
	case "<=":
		leftComparable := left.(object.Comparable)
		rightComparable := right.(object.Comparable)
		result, err := leftComparable.Compare(rightComparable)
		return nativeBoolToBoolean(result == object.LT || result == object.EQ), err
	case ">=":
		leftComparable := left.(object.Comparable)
		rightComparable := right.(object.Comparable)
		result, err := leftComparable.Compare(rightComparable)
		return nativeBoolToBoolean(result == object.GT || result == object.EQ), err
	case "||":
		leftBool := left.(*object.Boolean)
		rightBool := right.(*object.Boolean)
		return nativeBoolToBoolean(leftBool.Value || rightBool.Value), nil
	case "&&":
		leftBool := left.(*object.Boolean)
		rightBool := right.(*object.Boolean)
		return nativeBoolToBoolean(leftBool.Value && rightBool.Value), nil

	default:
		return nil, nil
	}
}

func evalPlusInfixOperator(left, right object.Object) (object.Object, error) {
	if left.Type() == object.IntegerType && right.Type() == object.IntegerType {
		newValue := left.(*object.Integer).Value + right.(*object.Integer).Value
		return &object.Integer{Value: newValue}, nil
	}

	return nil, nil
}

func evalMinusInfixOperator(left, right object.Object) (object.Object, error) {
	if left.Type() == object.IntegerType && right.Type() == object.IntegerType {
		newValue := left.(*object.Integer).Value - right.(*object.Integer).Value
		return &object.Integer{Value: newValue}, nil
	}

	return nil, nil
}

func evalAsteriskInfixOperator(left, right object.Object) (object.Object, error) {
	if left.Type() == object.IntegerType && right.Type() == object.IntegerType {
		newValue := left.(*object.Integer).Value * right.(*object.Integer).Value
		return &object.Integer{Value: newValue}, nil
	}

	return nil, nil
}

func evalAsteriskSlashOperator(left, right object.Object) (object.Object, error) {
	if left.Type() == object.IntegerType && right.Type() == object.IntegerType {
		newValue := left.(*object.Integer).Value / right.(*object.Integer).Value
		return &object.Integer{Value: newValue}, nil
	}

	return nil, nil
}

func nativeBoolToBoolean(b bool) *object.Boolean {
	if b {
		return &object.True
	}

	return &object.False
}
