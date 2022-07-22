package karboscript

import (
	"errors"
	"fmt"
)

type Program struct {
	Opcodes               []*Opcode
	codePointer           *int
	running               *bool
	callstack             []int
	functionArgsStack     []any
	functionArgumentCount *int
	expresionStack        []any
	returnExp             *any
}

func (program *Program) popExp() (any, error) {
	x := len(program.expresionStack) - 1

	if x < 0 {
		return 0, errors.New("No value on expresion stack!")
	}

	value := program.expresionStack[x]

	program.expresionStack = program.expresionStack[0:x]
	return value, nil
}

func (program *Program) pushExp(value any) {
	program.expresionStack = append(program.expresionStack, value)
}

func Execute(stack *[]*Opcode) error {
	killSwitch := 1000
	codePointer := len(*stack) - 2
	running := true
	callstack := []int{}
	functionArgumentCount := 0

	program := Program{
		*stack, &codePointer, &running, callstack, []any{}, &functionArgumentCount, []any{}, nil,
	}

	for *program.running {
		killSwitch--

		if killSwitch < 0 {
			running = false
		}

		err := executeOpcode(&program)
		if err != nil {
			return err
		}
	}

	return nil
}

func getNextOpcode(program *Program) (*Opcode, error) {
	*program.codePointer++

	if *program.codePointer > len((*program).Opcodes) {
		return nil, nil
	}

	return program.Opcodes[*program.codePointer-1], nil
}

func executeOpcode(program *Program) error {
	opcode, err := getNextOpcode(program)

	if opcode == nil {
		return nil
	}

	if err != nil {
		return err
	}

	if opcode.Operation == "exit" {
		*program.running = false
		return nil
	}

	if opcode.Operation == "push_exp" {
		program.pushExp(opcode.Arguments[0])
		return nil
	}
	if opcode.Operation == "set_return" {
		value, err := program.popExp()
		if err != nil {
			return err
		}
		program.returnExp = &value
		return nil
	}

	if opcode.Operation == "exp_call" {
		error := mathOperation(program, opcode)
		if error != nil {
			return error
		}

		return nil
	}

	if opcode.Operation == "call_function" {
		program.returnExp = nil
		if functionName, ok := opcode.Arguments[0].(string); ok {
			if val, ok := buildInFunctions[functionName]; ok {
				if count, ok := opcode.Arguments[1].(int); ok {
					*program.functionArgumentCount = count
				} else {
					return errors.New("call_function needs to have number of arguments as second parameter")
				}

				err := val(program)
				if err != nil {
					return err
				}
				return nil
			}

			program.callstack = append(program.callstack, *program.codePointer)
			if count, ok := opcode.Arguments[1].(int); ok {
				*program.functionArgumentCount = count
			} else {
				return errors.New("call_function needs to have number of arguments as second parameter")
			}

			*program.codePointer, err = findLabel(program, "_function."+functionName)
			if err != nil {
				return err
			}
			return nil
		} else {
			return errors.New("call_function opcode has wrong argument")
		}

	}

	if opcode.Operation == "push_function_arg" {
		if opcode.Arguments[0] == "pop_exp" {
			x, error := (*program).popExp()
			if error != nil {
				return error
			}
			program.functionArgsStack = append(program.functionArgsStack, x)
			//program.functionArgsStack = append(program.functionArgsStack, opcode.Arguments...)

		} else {
			program.functionArgsStack = append(program.functionArgsStack, opcode.Arguments...)
		}
	}

	if opcode.Operation == "function_return" {
		//todo clear functions args from stack
		newCodePointer := program.callstack[len(program.callstack)-1]
		program.callstack = program.callstack[0 : len(program.callstack)-1]
		*program.codePointer = newCodePointer

		if (program.returnExp != nil){
			program.pushExp(*program.returnExp)
		}
	}

	return nil
}

func mathOperation(program *Program, opcode *Opcode) error {
	operation := fmt.Sprintf("%v", opcode.Arguments[0])

	if operation == "math_op_*" || operation == "math_op_/" || operation == "math_op_+" || operation == "math_op_-" {
		val1, err1 := program.popExp()
		if err1 != nil {
			return err1
		}
		val2, err2 := program.popExp()
		if err2 != nil {
			return err2
		}
		if val1, ok := val1.(int); ok {
			if val2, ok := val2.(int); ok {
				switch operation {
				case "math_op_*":
					program.pushExp(val1 * val2)
				case "math_op_/":
					if val2 == 0 {
						return errors.New("Division by 0!")
					}
					program.pushExp(val1 / val2)
				case "math_op_+":
					program.pushExp(val1 + val2)
				case "math_op_-":
					program.pushExp(val1 - val2)
				}

				return nil
			}
		}

		return errors.New("Can't perform math operation!")
	}

	return errors.New("Wrong operation!")
}

func getFunctionArguments(program *Program) []any {
	arguments := program.functionArgsStack[0:*program.functionArgumentCount]

	program.functionArgsStack = program.functionArgsStack[0 : len(program.functionArgsStack)-*program.functionArgumentCount]
	*program.functionArgumentCount = 0

	return arguments
}

func findLabel(program *Program, label string) (int, error) {
	for i, opcode := range program.Opcodes {
		if opcode.Label != nil && *opcode.Label == label {
			return i, nil
		}
	}

	return 0, errors.New("can't find label: " + label)
}
