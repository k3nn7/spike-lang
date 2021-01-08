package eval

import (
	"spike-interpreter-go/spike/eval/object"
	"spike-interpreter-go/spike/parser/ast"

	"github.com/pkg/errors"
)

func Eval(node ast.Node, environment *object.Environment) (object.Object, error) {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, environment)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, environment)
	case *ast.Integer:
		return &object.Integer{Value: node.Value}, nil
	case *ast.Boolean:
		return evalBoolean(node)
	case *ast.Array:
		array := &object.Array{
			Elements: make([]object.Object, 0, len(node.Elements)),
		}

		for _, element := range node.Elements {
			evaluatedElement, err := Eval(element, environment)
			if err != nil {
				return nil, err
			}
			array.Elements = append(array.Elements, evaluatedElement)
		}

		return array, nil

	case *ast.Hash:
		hash := &object.Hash{
			Pairs: make(map[object.HashKey]object.HashPair, len(node.Pairs)),
		}

		for key, value := range node.Pairs {
			evaluatedKey, err := Eval(key, environment)
			if err != nil {
				return nil, err
			}
			evalutedValue, err := Eval(value, environment)
			if err != nil {
				return nil, err
			}

			hashable, isHashable := evaluatedKey.(object.Hashable)
			if !isHashable {
				return nil, errors.Errorf("%s does not implement Hashable", evaluatedKey.Type())
			}

			hash.Pairs[hashable.GetHashKey()] = object.HashPair{
				Key:   evaluatedKey,
				Value: evalutedValue,
			}
		}

		return hash, nil

	case *ast.PrefixExpression:
		right, err := Eval(node.Right, environment)
		if err != nil {
			return nil, err
		}
		return evalPrefixExpression(right, node.Operator)
	case *ast.InfixExpression:
		left, err := Eval(node.Left, environment)
		if err != nil {
			return nil, err
		}
		right, err := Eval(node.Right, environment)
		if err != nil {
			return nil, err
		}

		return evalInfixExpression(left, right, node.Operator)
	case *ast.IfExpression:
		condition, _ := Eval(node.Condition, environment)
		if condition.Equal(&object.True) {
			return Eval(node.Then, environment)
		} else {
			return Eval(node.Else, environment)
		}
	case *ast.BlockStatement:
		return evalStatements(node.Statements, environment)
	case *ast.ReturnStatement:
		result, _ := Eval(node.Result, environment)
		return &object.Return{Value: result}, nil
	case *ast.LetStatement:
		result, _ := Eval(node.Value, environment)
		environment.Set(node.Name.Value, result)
	case *ast.Identifier:
		return evalIdentifier(node.Value, environment)
	case *ast.FunctionExpression:
		return &object.Function{
			Parameters:  node.Parameters,
			Body:        node.Body,
			Environment: environment,
		}, nil
	case *ast.CallExpression:
		function, _ := Eval(node.Function, environment)
		arguments, _ := evalExpressions(node.Arguments, environment)
		return applyFunction(function, arguments)
	case *ast.String:
		return &object.String{Value: node.Value}, nil
	case *ast.IndexExpression:
		evaluatedArray, err := Eval(node.Array, environment)
		if err != nil {
			return nil, err
		}
		evaluatedIndex, err := Eval(node.Index, environment)
		if err != nil {
			return nil, err
		}

		switch evaluatedArray.(type) {
		case *object.Array:
			arrayObject := evaluatedArray.(*object.Array)
			integerObject, ok := evaluatedIndex.(*object.Integer)
			if !ok {
				return nil, errors.New("only integer can be used as index")
			}

			return arrayObject.Elements[integerObject.Value], nil
		case *object.Hash:
			hashObject := evaluatedArray.(*object.Hash)
			hashable, ok := evaluatedIndex.(object.Hashable)
			if !ok {
				return nil, errors.New("hash key must implement Hashable interface")
			}

			return hashObject.Get(hashable)
		default:
			return nil, errors.New("index can be used only on array")
		}
	default:
		return nil, errors.Errorf("Trying to evaluate unknown node: %T: %#v", node, node)
	}
	return nil, nil
}

func applyFunction(function object.Object, arguments []object.Object) (object.Object, error) {
	if builtinFunction, ok := function.(*object.BuiltinFunction); ok {
		return builtinFunction.Function(arguments...)
	}

	functionObject, ok := function.(*object.Function)
	if !ok {
		return nil, nil
	}

	extendedEnvironment := object.ExtendEnvironment(functionObject.Environment)
	for i, identifier := range functionObject.Parameters {
		extendedEnvironment.Set(identifier.Value, arguments[i])
	}

	result, err := Eval(functionObject.Body, extendedEnvironment)
	if err != nil {
		return nil, err
	}

	if returnValue, ok := result.(*object.Return); ok {
		return returnValue.Value, nil
	}

	return result, nil
}

func evalProgram(program *ast.Program, environment *object.Environment) (object.Object, error) {
	var result object.Object
	var err error
	for _, statement := range program.Statements {
		result, err = Eval(statement, environment)
		if err != nil {
			return nil, err
		}

		if returnValue, ok := result.(*object.Return); ok {
			return returnValue.Value, nil
		}
	}

	return result, err
}

func evalStatements(statements []ast.Statement, environment *object.Environment) (object.Object, error) {
	var result object.Object
	var err error
	for _, statement := range statements {
		result, err = Eval(statement, environment)
		if err != nil {
			return nil, err
		}

		if _, ok := result.(*object.Return); ok {
			return result, nil
		}
	}

	return result, err
}

func evalExpressions(expressions []ast.Expression, environment *object.Environment) ([]object.Object, error) {
	result := make([]object.Object, 0)

	for _, expression := range expressions {
		evaluated, _ := Eval(expression, environment)
		result = append(result, evaluated)
	}

	return result, nil
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
		return nil, errors.Errorf("type mismatch: !%s", right.Type())
	}
}

func evalMinusOperator(right object.Object) (object.Object, error) {
	switch rightObject := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -rightObject.Value}, nil
	default:
		return nil, errors.Errorf("type mismatch: -%s", right.Type())
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
		equal := left.Equal(right)
		return nativeBoolToBoolean(equal), nil
	case "!=":
		equal := left.Equal(right)
		return nativeBoolToBoolean(!equal), nil
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

	return nil, errors.Errorf("type mismatch: %s + %s", left.Type(), right.Type())
}

func evalMinusInfixOperator(left, right object.Object) (object.Object, error) {
	if left.Type() == object.IntegerType && right.Type() == object.IntegerType {
		newValue := left.(*object.Integer).Value - right.(*object.Integer).Value
		return &object.Integer{Value: newValue}, nil
	}

	return nil, errors.Errorf("type mismatch: %s - %s", left.Type(), right.Type())
}

func evalAsteriskInfixOperator(left, right object.Object) (object.Object, error) {
	if left.Type() == object.IntegerType && right.Type() == object.IntegerType {
		newValue := left.(*object.Integer).Value * right.(*object.Integer).Value
		return &object.Integer{Value: newValue}, nil
	}

	return nil, errors.Errorf("type mismatch: %s * %s", left.Type(), right.Type())
}

func evalAsteriskSlashOperator(left, right object.Object) (object.Object, error) {
	if left.Type() == object.IntegerType && right.Type() == object.IntegerType {
		newValue := left.(*object.Integer).Value / right.(*object.Integer).Value
		return &object.Integer{Value: newValue}, nil
	}

	return nil, errors.Errorf("type mismatch: %s / %s", left.Type(), right.Type())
}

func nativeBoolToBoolean(b bool) *object.Boolean {
	if b {
		return &object.True
	}

	return &object.False
}

func evalIdentifier(name string, environment *object.Environment) (object.Object, error) {
	variable, err := environment.Get(name)
	if err == nil {
		return variable, nil
	}

	if builtin, ok := builtins[name]; ok {
		return builtin, nil
	}

	return nil, err
}
