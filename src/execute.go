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
	scopes                []*Scope
	lastScope             *Scope
}

type Scope struct {
	expresionStack []any
	variable       map[string]any
	isFinal        bool
}

func (program *Program) getScope(depth int) *Scope {
	x := len(program.scopes) - 1 - depth
	return program.scopes[x]
}

func (scope *Scope) popExp() (any, error) {
	x := len(scope.expresionStack) - 1

	if x < 0 {
		return 0, errors.New("No value on expresion stack!")
	}

	value := scope.expresionStack[x]

	scope.expresionStack = scope.expresionStack[0:x]
	return value, nil
}

func (scope *Scope) pushExp(value any) {
	scope.expresionStack = append(scope.expresionStack, value)
}

func (program *Program) addScope() {
	program.scopes = append(program.scopes, &Scope{
		expresionStack: []any{},
		variable:       map[string]any{},
		isFinal:        false,
	})
}

func (program *Program) subScope() *Scope {
	x := len(program.scopes) - 1

	value := program.scopes[x]

	program.scopes = program.scopes[0:x]
	return value
}

func (program *Program) getVariable(name string) any {
	meetFinal := false
	for i := 0; i < len(program.scopes)-1; i++ {
		if meetFinal {
			return nil
		}

		if program.getScope(i).isFinal {
			meetFinal = true
		}

		if _, ok := program.getScope(i).variable[name]; ok {
			return program.getScope(i).variable[name]
		}
	}

	return nil
}

func (program *Program) getVariableScopePosition(name string) int {
	meetFinal := false
	for i := 0; i < len(program.scopes)-1; i++ {
		if meetFinal {
			return -1
		}

		if program.getScope(i).isFinal {
			meetFinal = true
		}

		if _, ok := program.getScope(i).variable[name]; ok {
			return i
		}
	}

	return -1
}

func Execute(stack *[]*Opcode) error {
	killSwitch := 100000
	codePointer := len(*stack) - 2
	running := true
	callstack := []int{}
	functionArgumentCount := 0

	program := Program{
		*stack, &codePointer, &running, callstack, []any{}, &functionArgumentCount, []*Scope{}, nil,
	}
	program.addScope()

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

	//fmt.Println(opcode)

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
		program.getScope(0).pushExp(opcode.Arguments[0])
		return nil
	}
	if opcode.Operation == "push_exp_var" {
		if name, ok := opcode.Arguments[0].(string); ok {
			program.getScope(0).pushExp(program.getVariable(name))
		}
		return nil
	}

	if opcode.Operation == "add_scope" {
		program.addScope()
		return nil
	}

	if opcode.Operation == "sub_scope" {
		program.lastScope = program.subScope()

		return nil
	}

	if opcode.Operation == "set_local_var_arg" {
		if name, ok := opcode.Arguments[0].(string); ok {
			varScopePosition := program.getVariableScopePosition(name)

			if varScopePosition > -1 {
				program.getScope(varScopePosition).variable[name] = program.popFunctionArgument()
			} else {
				program.getScope(0).variable[name] = program.popFunctionArgument()
			}
		}

		return nil
	}

	if opcode.Operation == "set_local_var_exp" {
		if name, ok := opcode.Arguments[0].(string); ok {
			varScopePosition := program.getVariableScopePosition(name)

			if varScopePosition > -1 {
				program.getScope(varScopePosition).variable[name], err = program.lastScope.popExp()
			} else {
				program.getScope(0).variable[name], err = program.lastScope.popExp()
			}

			if err != nil {
				return err
			}
		}

		return nil
	}

	if opcode.Operation == "push_bellow" {
		value, err := program.lastScope.popExp()
		if err != nil {
			return err
		}
		program.getScope(1).pushExp(value)
		return nil
	}

	if opcode.Operation == "exp_call" {
		error := mathOperation(program, opcode)
		if error != nil {
			return error
		}

		return nil
	}

	if opcode.Operation == "if" {
		lastVal, err := program.lastScope.popExp()
		if err != nil {
			return err
		}

		if val, ok := lastVal.(bool); ok {
			if label, ok := opcode.Arguments[1].(string); ok {
				if !val {
					*program.codePointer, err = findLabel(program, label)
				}
				if err != nil {
					return err
				}
			}
		} else {
			return errors.New("Condition must return bool")
		}
		return nil
	}
	if opcode.Operation == "jmp" {

		if label, ok := opcode.Arguments[0].(string); ok {
			*program.codePointer, err = findLabel(program, label)

			if err != nil {
				return err
			}
		}

		return nil
	}

	if opcode.Operation == "while" || opcode.Operation == "for" {
		lastVal, err := program.lastScope.popExp()
		if err != nil {
			return err
		}

		if val, ok := lastVal.(bool); ok {
			if label, ok := opcode.Arguments[1].(string); ok {
				if !val {
					*program.codePointer, err = findLabel(program, label)
					if err != nil {
						return err
					}
				}
			}
		} else {
			return errors.New("Condition must return bool")
		}
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
			program.addScope()
			program.getScope(0).isFinal = true
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
			x, error := (*program).lastScope.popExp()
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
		program.subScope()
		//todo clear functions args from stack
		newCodePointer := program.callstack[len(program.callstack)-1]
		program.callstack = program.callstack[0 : len(program.callstack)-1]
		*program.codePointer = newCodePointer
	}

	return nil
}

func mathOperation(program *Program, opcode *Opcode) error {
	operation := fmt.Sprintf("%v", opcode.Arguments[0])

	if operation == "*" || operation == "/" || operation == "+" || operation == "-" {
		val1, err1 := program.getScope(0).popExp()
		if err1 != nil {
			return err1
		}
		val2, err2 := program.getScope(0).popExp()
		if err2 != nil {
			return err2
		}
		if val1, ok := val1.(int); ok {
			if val2, ok := val2.(int); ok {
				switch operation {
				case "*":
					program.getScope(0).pushExp(val2 * val1)
				case "/":
					if val1 == 0 {
						return errors.New("Division by 0!")
					}
					program.getScope(0).pushExp(val2 / val1)
				case "+":
					program.getScope(0).pushExp(val2 + val1)
				case "-":
					program.getScope(0).pushExp(val2 - val1)
				}

				return nil
			}
		}

		return errors.New("Can't perform math operation!")
	}

	if operation == "!=" || operation == "==" || operation == ">" || operation == ">=" || operation == "<=" || operation == "<" {
		val1, err1 := program.getScope(0).popExp()
		if err1 != nil {
			return err1
		}
		val2, err2 := program.getScope(0).popExp()
		if err2 != nil {
			return err2
		}
		if val1, ok := val1.(int); ok {
			if val2, ok := val2.(int); ok {
				switch operation {
				case "==":
					program.getScope(0).pushExp(val2 == val1)
				case "!=":
					program.getScope(0).pushExp(val2 != val1)
				case ">":
					program.getScope(0).pushExp(val2 > val1)
				case ">=":
					program.getScope(0).pushExp(val2 >= val1)
				case "<":
					program.getScope(0).pushExp(val2 < val1)
				case "<=":
					program.getScope(0).pushExp(val2 <= val1)
				}

				return nil
			}
		}

	}
	return errors.New("Wrong operation!")
}

func getFunctionArguments(program *Program) []any {
	x := len(program.functionArgsStack) - *program.functionArgumentCount
	x1 := len(program.functionArgsStack)
	arguments := program.functionArgsStack[x:x1]

	program.functionArgsStack = program.functionArgsStack[0 : len(program.functionArgsStack)-*program.functionArgumentCount]
	*program.functionArgumentCount = 0

	return arguments
}

func (program *Program) popFunctionArgument() any {
	x := len(program.functionArgsStack) - 1
	argument := program.functionArgsStack[x]

	program.functionArgsStack = program.functionArgsStack[0 : len(program.functionArgsStack)-1]
	*program.functionArgumentCount = 0

	return argument
}

func findLabel(program *Program, label string) (int, error) {
	for i, opcode := range program.Opcodes {
		if opcode.Label != nil && *opcode.Label == label {
			return i, nil
		}
	}

	return 0, errors.New("can't find label: " + label)
}
