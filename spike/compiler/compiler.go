package compiler

import (
	"fmt"
	"spike-interpreter-go/spike/code"
	"spike-interpreter-go/spike/eval/object"
	"spike-interpreter-go/spike/parser/ast"

	"github.com/pkg/errors"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (compiler *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, statement := range node.Statements {
			err := compiler.Compile(statement)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := compiler.Compile(node.Expression)
		if err != nil {
			return err
		}
		compiler.emit(code.OpPop)

	case *ast.InfixExpression:
		if node.Operator == "<" {
			err := compiler.Compile(node.Right)
			if err != nil {
				return err
			}

			err = compiler.Compile(node.Left)
			if err != nil {
				return err
			}

			compiler.emit(code.OpGreaterThan)

			return nil
		}

		err := compiler.Compile(node.Left)
		if err != nil {
			return err
		}

		err = compiler.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			compiler.emit(code.OpAdd)
		case "-":
			compiler.emit(code.OpSub)
		case "*":
			compiler.emit(code.OpMul)
		case "/":
			compiler.emit(code.OpDiv)
		case "==":
			compiler.emit(code.OpEqual)
		case "!=":
			compiler.emit(code.OpNotEqual)
		case ">":
			compiler.emit(code.OpGreaterThan)
		default:
			return fmt.Errorf("unknown operator: %s", node.Operator)
		}

	case *ast.PrefixExpression:
		err := compiler.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			compiler.emit(code.OpBang)
		case "-":
			compiler.emit(code.OpMinus)
		default:
			return errors.Errorf("invalid prefix operator: %s", node.Operator)
		}

	case *ast.Integer:
		integer := &object.Integer{Value: node.Value}
		compiler.emit(code.OpConstant, compiler.addConstant(integer))

	case *ast.Boolean:
		if node.Value {
			compiler.emit(code.OpTrue)
		} else {
			compiler.emit(code.OpFalse)
		}
	}

	return nil
}

func (compiler *Compiler) addConstant(obj object.Object) int {
	compiler.constants = append(compiler.constants, obj)
	return len(compiler.constants) - 1
}

func (compiler *Compiler) emit(opcode code.Opcode, operands ...int) int {
	instruction, _ := code.Make(opcode, operands...)

	newInstructionIndex := len(compiler.instructions)
	compiler.instructions = append(compiler.instructions, instruction...)
	return newInstructionIndex
}

func (compiler *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: compiler.instructions,
		Constants:    compiler.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
