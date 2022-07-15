package main

import "errors"

type Program struct {
	Opcodes           []*Opcode
	codePointer       *int
	running           *bool
	callstack         []int
	functionArgsStack []any
	expresionOutout   *any
}

func execute(stack *[]*Opcode) {
	killSwitch := 1000
	codePointer := len(*stack) - 1
	running := true
	callstack := []int{}

	program := Program{
		*stack, &codePointer, &running, callstack, []any{}, nil,
	}

	for *program.running {
		killSwitch--

		if killSwitch < 0 {
			running = false
		}

		err := executeOpcode(&program)
		if err != nil {
			return
		}
	}
}

func getNextOpcode(program *Program) (*Opcode, error) {
	*program.codePointer++

	if *program.codePointer > len((*program).Opcodes) {
		return nil, errors.New("accessing outside of the program")
	}

	return program.Opcodes[*program.codePointer-1], nil
}

func executeOpcode(program *Program) error {
	opcode, err := getNextOpcode(program)

	if err != nil {
		return err
	}

	if opcode.Operation == "exit" {
		*program.running = false
		return nil
	}

	if opcode.Operation == "set_expresion_output" {
		program.expresionOutout = &opcode.Arguments[0]
		return nil
	}

	if opcode.Operation == "call_function" {
		if functionName, ok := opcode.Arguments[0].(string); ok {
			if val, ok := buildInFunctions[functionName]; ok {
				val(program)
				program.functionArgsStack = []any{}
				return nil
			}

			program.callstack = append(program.callstack, *program.codePointer)
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
		program.functionArgsStack = append(program.functionArgsStack, opcode.Arguments...)
	}

	if opcode.Operation == "function_return" {
		program.functionArgsStack = []any{}
		newCodePointer := program.callstack[len(program.callstack)-1]
		program.callstack = program.callstack[0 : len(program.callstack)-1]
		*program.codePointer = newCodePointer
	}

	return nil
}

func findLabel(program *Program, label string) (int, error) {
	for i, opcode := range program.Opcodes {
		if opcode.Label != nil && *opcode.Label == label {
			return i, nil
		}
	}

	return 0, errors.New("can't find label: " + label)
}
