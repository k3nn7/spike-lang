package compiler

import (
	"fmt"
	"spike-interpreter-go/spike/code"
	"spike-interpreter-go/spike/eval/object"
	"spike-interpreter-go/spike/parser/ast"

	"github.com/pkg/errors"
)

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
	symbolTable  *SymbolTable

	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
		symbolTable:  NewSymbolTable(),
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

	case *ast.BlockStatement:
		for _, statement := range node.Statements {
			err := compiler.Compile(statement)
			if err != nil {
				return err
			}
		}

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

	case *ast.IfExpression:
		err := compiler.Compile(node.Condition)
		if err != nil {
			return err
		}

		jumpNotTrueIndex := compiler.emit(code.OpJumpNotTrue, -1)

		err = compiler.Compile(node.Then)
		if err != nil {
			return err
		}

		if compiler.lastInstruction.Opcode == code.OpPop {
			compiler.removeLastInstruction()
		}

		if node.Else == nil {
			jumpIndex := compiler.emit(code.OpJump, -1)
			afterJumpIndex := len(compiler.instructions)
			compiler.emit(code.OpNull)
			afterNullIndex := len(compiler.instructions)

			compiler.changeOperand(jumpIndex, afterNullIndex)
			compiler.changeOperand(jumpNotTrueIndex, afterJumpIndex)
		} else {
			jumpIndex := compiler.emit(code.OpJump, -1)

			afterThenIndex := len(compiler.instructions)
			compiler.changeOperand(jumpNotTrueIndex, afterThenIndex)

			err := compiler.Compile(node.Else)
			if err != nil {
				return err
			}

			if compiler.lastInstruction.Opcode == code.OpPop {
				compiler.removeLastInstruction()
			}

			afterElseIndex := len(compiler.instructions)
			compiler.changeOperand(jumpIndex, afterElseIndex)
		}

	case *ast.LetStatement:
		err := compiler.Compile(node.Value)
		if err != nil {
			return err
		}

		symbol := compiler.symbolTable.Define(node.Name.Value)
		compiler.emit(code.OpSetGlobal, symbol.Index)

	case *ast.Identifier:
		symbol, ok := compiler.symbolTable.Resolve(node.Value)
		if !ok {
			return errors.Errorf("unable to resolve identifier: %s", node.Value)
		}

		compiler.emit(code.OpGetGlobal, symbol.Index)
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

	compiler.previousInstruction = compiler.lastInstruction
	compiler.lastInstruction = EmittedInstruction{
		Opcode:   opcode,
		Position: newInstructionIndex,
	}

	return newInstructionIndex
}

func (compiler *Compiler) removeLastInstruction() {
	compiler.instructions = compiler.instructions[:compiler.lastInstruction.Position]
	compiler.lastInstruction = compiler.previousInstruction
}

func (compiler *Compiler) changeOperand(instructionIndex, operand int) {
	opcode := code.Opcode(compiler.instructions[instructionIndex])
	newInstruction, _ := code.Make(opcode, operand)

	compiler.replaceInstruction(instructionIndex, newInstruction)
}

func (compiler *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: compiler.instructions,
		Constants:    compiler.constants,
	}
}

func (compiler *Compiler) replaceInstruction(instructionIndex int, instruction []byte) {
	for i := 0; i < len(instruction); i++ {
		compiler.instructions[instructionIndex+i] = instruction[i]
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
