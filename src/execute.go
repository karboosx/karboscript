package karboscript

import "errors"

type Program struct {
	Opcodes               []*Opcode
	codePointer           *int
	running               *bool
	callstack             []int
	functionArgsStack     []any
	expresionOutput       *any
	functionArgumentCount *int
}

func Execute(stack *[]*Opcode) error {
	killSwitch := 1000
	codePointer := len(*stack) - 2
	running := true
	callstack := []int{}
	functionArgumentCount := 0

	program := Program{
		*stack, &codePointer, &running, callstack, []any{}, nil, &functionArgumentCount,
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

	if opcode.Operation == "set_expresion_output" {
		program.expresionOutput = &opcode.Arguments[0]
		return nil
	}

	if opcode.Operation == "call_function" {
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
		if opcode.Arguments[0] == "expresion_output" {
			program.functionArgsStack = append(program.functionArgsStack, *program.expresionOutput)
			//program.functionArgsStack = append(program.functionArgsStack, opcode.Arguments...)

		} else {
			program.functionArgsStack = append(program.functionArgsStack, opcode.Arguments...)
		}
	}

	if opcode.Operation == "function_return" {
		//program.functionArgsStack = []any{}
		newCodePointer := program.callstack[len(program.callstack)-1]
		program.callstack = program.callstack[0 : len(program.callstack)-1]
		*program.codePointer = newCodePointer
	}

	return nil
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
