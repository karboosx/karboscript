package karboscript

import (
	"errors"
	"fmt"
)

type Call struct {
	returnPointer int
	returnType    *VarType
}

type Program struct {
	Opcodes               []*Opcode
	codePointer           *int
	running               *bool
	callstack             []Call
	functionArgsStack     []any
	functionArgumentCount *int
	scopes                []*Scope
	lastSubScope          *Scope
}

type Var struct {
	value   any
	varType VarType
}

type Scope struct {
	expresionStack []any
	variable       map[string]*Var
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
		variable:       map[string]*Var{},
		isFinal:        false,
	})
}

func (program *Program) subScope() *Scope {
	x := len(program.scopes) - 1

	value := program.scopes[x]

	program.scopes = program.scopes[0:x]
	return value
}

func (program *Program) getVariable(name string) *Var {
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
	callstack := []Call{}
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
			opcode := program.Opcodes[*program.codePointer-1]

			return errors.New(opcode.Position + ": " + err.Error())
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
			x := program.getVariable(name)
			if x == nil {
				return errors.New("Undeclared variable: " + name)
			}
			program.getScope(0).pushExp(x.value)
		}
		return nil
	}

	if opcode.Operation == "push_last_exp" {
		x, error := (*program).lastSubScope.popExp()
		if error != nil {
			return error
		}
		program.getScope(0).pushExp(x)
	}

	if opcode.Operation == "push_empty_arr" {
		program.getScope(0).pushExp([]any{})
		return nil
	}

	if opcode.Operation == "push_arr_exp" {
		newElement, error := (*program).lastSubScope.popExp()
		if error != nil {
			return error
		}

		arr, err1 := program.getScope(0).popExp()
		if err1 != nil {
			return err1
		}

		if arrayToAdd, ok := arr.([]any); ok {
			program.getScope(0).pushExp(append(arrayToAdd, newElement))

			return nil
		} else {
			return errors.New("variable is not array!")
		}
	}

	if opcode.Operation == "push_arr_call" {
		index, error := (*program).lastSubScope.popExp()
		if error != nil {
			return error
		}

		arr := program.getVariable(opcode.Arguments[0].(string))
		if arr == nil {
			return errors.New("Undeclared variable: " + opcode.Arguments[0].(string))
		}

		if array, ok := arr.value.([]any); ok {
			if index, ok := index.(int); ok {
				if index >= len(array) {
					return errors.New("Index out of range!")
				}

				program.getScope(0).pushExp(array[index])
			} else {
				return errors.New("Index is not integer!")
			}
		} else {
			return errors.New("variable is not array!")
		}
	}

	if opcode.Operation == "add_scope" {
		program.addScope()
		return nil
	}

	if opcode.Operation == "sub_scope" {
		program.lastSubScope = program.subScope()

		return nil
	}

	if opcode.Operation == "set_local_var_arg" {
		if varName, ok := opcode.Arguments[0].(string); ok {
			if name, ok := opcode.Arguments[1].(string); ok {
				varScopePosition := program.getVariableScopePosition(name)
				variable := Var{program.popFunctionArgument(), VarType{varName}}

				if err, ok := validateVariable(variable); !ok {
					return err
				}

				if varScopePosition > -1 {
					program.getScope(varScopePosition).variable[name] = &variable
				} else {
					program.getScope(0).variable[name] = &variable
				}
			}
		}

		return nil
	}

	if opcode.Operation == "set_local_var_exp" {

		if name, ok := opcode.Arguments[1].(string); ok {
			var varName = ""

			if varTypeFromOpcode, ok := opcode.Arguments[0].(string); ok && varTypeFromOpcode != "" {
				varName = varTypeFromOpcode
			} else {
				variableForType := program.getVariable(name)

				if variableForType != nil {
					varName = variableForType.varType.Value
				} else {
					return errors.New("Undeclared variable: " + name)
				}
			}

			if varName == "" {
				return errors.New("Broken variable: " + name)
			}

			varScopePosition := program.getVariableScopePosition(name)
			varValue, err := program.lastSubScope.popExp()

			if err != nil {
				return err
			}

			variable := Var{varValue, VarType{varName}}

			if err, ok := validateVariable(variable); !ok {
				return err
			}

			if varScopePosition > -1 {

				program.getScope(varScopePosition).variable[name] = &variable
			} else {
				program.getScope(0).variable[name] = &variable
			}
		}

		return nil
	}

	if opcode.Operation == "set_array_var_exp" {

		if name, ok := opcode.Arguments[0].(string); ok {

			variableForType := program.getVariable(name)

			if variableForType == nil {
				return errors.New("Undeclared variable: " + name)
			}

			expression, err := program.getScope(0).popExp()

			if err != nil {
				return err
			}

			index, err := program.getScope(0).popExp()

			if err != nil {
				return err
			}

			variable := program.getVariable(name)

			if variable, err := variable.value.([]any); err {
				if index, ok := index.(int); ok {
					if index >= len(variable) {
						return errors.New("Index out of range!")
					}

					variable[index] = expression
				} else {
					return errors.New("Index is not integer!")
				}
			} else {
				return errors.New("variable is not array!")
			}
		}

		return nil
	}

	if opcode.Operation == "push_bellow" {
		value, err := program.lastSubScope.popExp()
		if err != nil {
			return err
		}

		newCodePointer := program.callstack[len(program.callstack)-1]

		if newCodePointer.returnType != nil {
			if err, ok := validateReturnType(newCodePointer, value); !ok {
				return err
			}
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
		lastVal, err := program.lastSubScope.popExp()
		if err != nil {
			return err
		}

		if val, ok := lastVal.(bool); ok {
			if label, ok := opcode.Arguments[0].(string); ok {
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
		lastVal, err := program.lastSubScope.popExp()
		if err != nil {
			return err
		}

		if val, ok := lastVal.(bool); ok {
			if label, ok := opcode.Arguments[0].(string); ok {
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

	if opcode.Operation == "forinc_start" {
		if variable, ok := opcode.Arguments[0].(string); ok {
			val := program.getVariable(variable)

			if val == nil {
				return errors.New("forint use uninitalized variable!")
			}

			valEnd := program.getVariable(variable + "_end")

			if val == nil {
				return errors.New("forint use uninitalized variable!")
			}
			if a, ok := val.value.(int); ok {
				if b, ok := valEnd.value.(int); ok {
					if a == b {
						if label, ok := opcode.Arguments[1].(string); ok {
							*program.codePointer, err = findLabel(program, label)
							if err != nil {
								return err
							}
						}
					}
				}
			}

			return nil
		} else {
			return errors.New("??")
		}
	}

	if opcode.Operation == "forinc" {
		if variable, ok := opcode.Arguments[0].(string); ok {
			val := program.getVariable(variable)

			if val == nil {
				return errors.New("forint use uninitalized variable!")
			}

			valEnd := program.getVariable(variable + "_end")

			if val == nil {
				return errors.New("forint use uninitalized variable!")
			}

			if a, ok := val.value.(int); ok {
				if b, ok := valEnd.value.(int); ok {
					if a < b {
						val.value = a + 1
					}
					if a > b {
						val.value = a - 1
					}
				}
			}

			if label, ok := opcode.Arguments[1].(string); ok {
				*program.codePointer, err = findLabel(program, label)
				if err != nil {
					return err
				}
			}

			return nil
		} else {
			return errors.New("??")
		}
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

			if count, ok := opcode.Arguments[1].(int); ok {
				*program.functionArgumentCount = count
			} else {
				return errors.New("call_function needs to have number of arguments as second parameter")
			}

			var returnType *VarType

			if len(opcode.Arguments) == 3 {
				if varType, ok := opcode.Arguments[2].(string); ok {
					returnType = &VarType{varType}
				}
			}

			program.callstack = append(program.callstack, Call{*program.codePointer, returnType})
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
		x, error := (*program).lastSubScope.popExp()
		if error != nil {
			return error
		}
		program.functionArgsStack = append(program.functionArgsStack, x)
	}

	if opcode.Operation == "function_return" {
		program.subScope()
		//todo clear functions args from stack

		newCodePointer := program.callstack[len(program.callstack)-1]
		program.callstack = program.callstack[0 : len(program.callstack)-1]
		*program.codePointer = newCodePointer.returnPointer
	}

	return nil
}

func validateReturnType(newCodePointer Call, value any) (error, bool) {

	if newCodePointer.returnType == nil {
		return nil, true
	}

	if newCodePointer.returnType.Value == "string" {
		if _, ok := value.(string); ok {
			return nil, true
		} else {
			return errors.New("return value is not string!"), false
		}
	}
	if newCodePointer.returnType.Value == "int" {
		if _, ok := value.(int); ok {
			return nil, true
		} else {
			return errors.New("return value is not int!"), false
		}
	}
	if newCodePointer.returnType.Value == "bool" {
		if _, ok := value.(bool); ok {
			return nil, true
		} else {
			return errors.New("return value is not bool!"), false
		}
	}
	if newCodePointer.returnType.Value == "float" {
		if _, ok := value.(float64); ok {
			return nil, true
		} else {
			return errors.New("return value is not float!"), false
		}
	}
	if newCodePointer.returnType.Value == "array" {
		if _, ok := value.([]any); ok {
			return nil, true
		} else {
			return errors.New("return value is not array!"), false
		}
	}

	return errors.New("cant validate return value!"), false
}

func validateVariable(variable Var) (error, bool) {
	if variable.varType.Value == "string" {
		if _, ok := variable.value.(string); ok {
			return nil, true
		} else {
			return errors.New("variable is not string!"), false
		}
	}
	if variable.varType.Value == "int" {
		if _, ok := variable.value.(int); ok {
			return nil, true
		} else {
			return errors.New("variable is not int!"), false
		}
	}
	if variable.varType.Value == "bool" {
		if _, ok := variable.value.(bool); ok {
			return nil, true
		} else {
			return errors.New("variable is not bool!"), false
		}
	}
	if variable.varType.Value == "float" {
		if _, ok := variable.value.(float64); ok {
			return nil, true
		} else {
			return errors.New("variable is not float!"), false
		}
	}
	if variable.varType.Value == "array" {
		if _, ok := variable.value.([]any); ok {
			return nil, true
		} else {
			return errors.New("variable is not array!"), false
		}
	}

	return errors.New("cant validate variable!"), false
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
