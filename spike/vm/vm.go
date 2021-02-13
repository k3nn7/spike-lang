package vm

import (
	"encoding/binary"
	"spike-interpreter-go/spike/code"
	"spike-interpreter-go/spike/compiler"
	"spike-interpreter-go/spike/object"

	"github.com/pkg/errors"
)

const (
	StackSize   = 2048
	MaxFrames   = 1024
	GlobalsSize = 65536
)

var (
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
	Null  = &object.Null{}
)

type VM struct {
	constants []object.Object
	globals   []object.Object

	stack []object.Object
	sp    int

	frames      []*Frame
	framesIndex int
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	return &VM{
		constants:   bytecode.Constants,
		stack:       make([]object.Object, StackSize),
		globals:     make([]object.Object, GlobalsSize),
		sp:          0,
		frames:      frames,
		framesIndex: 1,
	}
}

func NewWithGlobalStore(bytecode *compiler.Bytecode, globals []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = globals
	return vm
}

func (vm *VM) Run() error {
	var ip int
	var instructions code.Instructions
	var op code.Opcode

	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++

		ip = vm.currentFrame().ip
		instructions = vm.currentFrame().Instructions()
		op = code.Opcode(instructions[ip])

		switch op {
		case code.OpConstant:
			index := binary.BigEndian.Uint16(instructions[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.constants[index])
			if err != nil {
				return err

			}

		case code.OpAdd:
			err := vm.executePlusOperation()
			if err != nil {
				return err
			}

		case code.OpSub, code.OpMul, code.OpDiv:
			err := vm.executeBinaryIntegerOperation(op)
			if err != nil {
				return err
			}

		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}

		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}

		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()

		case code.OpBang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}

		case code.OpMinus:
			err := vm.executeMinusOperator()
			if err != nil {
				return err
			}

		case code.OpJump:
			jumpIndex := binary.BigEndian.Uint16(instructions[ip+1:])
			vm.currentFrame().ip = int(jumpIndex) - 1

		case code.OpJumpNotTrue:
			jumpIndex := binary.BigEndian.Uint16(instructions[ip+1:])
			vm.currentFrame().ip += 2

			condition := vm.pop().(*object.Boolean).Value
			if !condition {
				vm.currentFrame().ip = int(jumpIndex) - 1
			}

		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}

		case code.OpSetGlobal:
			globalIndex := binary.BigEndian.Uint16(instructions[ip+1:])
			vm.currentFrame().ip += 2

			vm.globals[globalIndex] = vm.pop()

		case code.OpGetGlobal:
			globalIndex := binary.BigEndian.Uint16(instructions[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}

		case code.OpArray:
			elementsCount := int(binary.BigEndian.Uint16(instructions[ip+1:]))
			vm.currentFrame().ip += 2

			elements := make([]object.Object, elementsCount)
			for i := 0; i < elementsCount; i++ {
				elements[i] = vm.stack[vm.sp-elementsCount+i]
			}

			vm.sp -= elementsCount

			array := &object.Array{Elements: elements}
			err := vm.push(array)
			if err != nil {
				return err
			}

		case code.OpHash:
			elementsCount := int(binary.BigEndian.Uint16(instructions[ip+1:]))
			vm.currentFrame().ip += 2

			pairs := make(map[object.HashKey]object.HashPair)

			for i := 0; i < elementsCount; i += 2 {
				key := vm.stack[vm.sp-elementsCount+i].(object.Hashable)
				value := vm.stack[vm.sp-elementsCount+i+1]

				pairs[key.GetHashKey()] = object.HashPair{
					Key:   key.(object.Object),
					Value: value,
				}
			}

			hash := &object.Hash{Pairs: pairs}
			err := vm.push(hash)
			if err != nil {
				return err
			}

		case code.OpIndex:
			index := vm.pop()
			array := vm.pop()

			switch array := array.(type) {
			case *object.Array:
				index, ok := index.(*object.Integer)
				if !ok {
					return errors.Errorf("Array index must be an integer, got: %s", index.Type())
				}

				if index.Value < 0 || index.Value >= int64(len(array.Elements)) {
					err := vm.push(Null)
					if err != nil {
						return err
					}
				} else {
					err := vm.push(array.Elements[index.Value])
					if err != nil {
						return err
					}
				}
			case *object.Hash:
				hashKey, ok := index.(object.Hashable)
				if !ok {
					return errors.Errorf("Object of type %s can not be used as a hash key", index.Type())
				}

				value, err := array.Get(hashKey)
				if err != nil {
					err = vm.push(Null)
					if err != nil {
						return err
					}
				} else {
					err = vm.push(value)
					if err != nil {
						return err
					}
				}
			}

		case code.OpCall:
			argumentsCount := int(instructions[ip+1])
			vm.currentFrame().ip++
			fn := vm.stack[vm.sp-1-argumentsCount]

			switch fn := fn.(type) {
			case *object.CompiledFunction:
				if fn.ParametersCount != argumentsCount {
					return errors.Errorf(
						"mismatched number of function call arguments. Expected %d, got %d",
						fn.ParametersCount,
						argumentsCount,
					)
				}

				frame := NewFrame(fn, vm.sp-argumentsCount)
				vm.pushFrame(frame)
				vm.sp = frame.basePointer + fn.LocalsCount

			case *object.BuiltinFunction:
				args := vm.stack[vm.sp-argumentsCount : vm.sp]

				result, err := fn.Function(args...)
				if err != nil {
					return err
				}
				err = vm.push(result)
				if err != nil {
					return err
				}

			default:
				return errors.Errorf("Calling non-function %T", fn)
			}

		case code.OpReturnValue:
			returnValue := vm.pop()

			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(returnValue)
			if err != nil {
				return err
			}

		case code.OpReturn:
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(Null)
			if err != nil {
				return err
			}

		case code.OpSetLocal:
			index := int(instructions[ip+1])
			vm.currentFrame().ip++

			vm.stack[vm.currentFrame().basePointer+index] = vm.pop()

		case code.OpGetLocal:
			index := int(instructions[ip+1])
			vm.currentFrame().ip++

			value := vm.stack[vm.currentFrame().basePointer+index]
			err := vm.push(value)
			if err != nil {
				return err
			}

		case code.OpGetBuiltin:
			index := int(instructions[ip+1])
			vm.currentFrame().ip++

			definition := object.Builtins[index]

			err := vm.push(definition)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (vm *VM) executePlusOperation() error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.IntegerType && right.Type() == object.IntegerType {
		leftValue := left.(*object.Integer).Value
		rightValue := right.(*object.Integer).Value

		result := &object.Integer{Value: leftValue + rightValue}
		return vm.push(result)
	} else if left.Type() == object.StringType && right.Type() == object.StringType {
		leftValue := left.(*object.String).Value
		rightValue := right.(*object.String).Value

		result := &object.String{Value: leftValue + rightValue}
		return vm.push(result)
	}

	return nil
}

func (vm *VM) executeBinaryIntegerOperation(opcode code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64
	switch opcode {
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	}
	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if right.Type() != left.Type() {
		return errors.Errorf("both operands must have same type, had: %s and %s", left.Type(), right.Type())
	}

	if right.Type() == object.IntegerType {
		return vm.executeIntegerComparison(left, right, op)
	}

	if right.Type() == object.BooleanType {
		return vm.executeBooleanComparison(left, right, op)
	}

	return errors.Errorf("unable to compare variables of type %s and %s", left.Type(), right.Type())
}

func (vm *VM) executeIntegerComparison(left object.Object, right object.Object, op code.Opcode) error {
	leftInt := left.(*object.Integer).Value
	rightInt := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBoolean(leftInt == rightInt))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBoolean(leftInt != rightInt))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBoolean(leftInt > rightInt))
	}

	return errors.Errorf("unexpected operation: %d", op)
}

func (vm *VM) executeBooleanComparison(left object.Object, right object.Object, op code.Opcode) error {
	leftBool := left.(*object.Boolean).Value
	rightBool := right.(*object.Boolean).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBoolean(leftBool == rightBool))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBoolean(leftBool != rightBool))
	}

	return errors.Errorf("unexpected operation: %d", op)
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	default:
		return errors.Errorf("invalid operand for bang prefix operator: %#v", operand)
	}
}

func (vm *VM) executeMinusOperator() error {
	value := vm.pop().(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
}

func nativeBoolToBoolean(nativeBool bool) object.Object {
	if nativeBool {
		return True
	} else {
		return False
	}
}

func (vm *VM) LastPoppedStackElement() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return errors.New("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushFrame(frame *Frame) {
	vm.frames[vm.framesIndex] = frame
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}

func (vm *VM) pop() object.Object {
	result := vm.stack[vm.sp-1]
	vm.sp--
	return result
}
