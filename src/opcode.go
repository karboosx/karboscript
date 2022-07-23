package karboscript

import (
	"errors"
	"strconv"
)

type OpCodes struct {
	List []*Opcode
}

type Opcode struct {
	Operation string
	Arguments []any
	Label     *string
}

type ParseError struct {
	Message string
}

func (m *ParseError) Error() string {
	return m.Message
}

func parseFunctionBody(stack *[]*Opcode, function *Function) error {
	parseBody(stack, function.Body)
	if (*stack)[len(*stack)-1].Operation != "function_return" {
		*stack = append(*stack, &Opcode{"function_return", []any{}, nil})
	}
	return nil
}

func parseBody(stack *[]*Opcode, statements []*Statement) error {
	for _, statement := range statements {
		if statement.FunctionCall != nil {
			parseFunctionCall(stack, statement.FunctionCall)
		}
		if statement.Expression != nil {
			parseExpresionWithNewScope(stack, statement.Expression)
		}
		if statement.ReturnStmt != nil {
			parseReturnStmt(stack, statement.ReturnStmt)
		}
		if statement.Assigment != nil {
			parseAssigment(stack, statement.Assigment)
		}
		if statement.If != nil {
			parseIf(stack, statement.If)
		}
	}

	return nil
}

func parseIf(stack *[]*Opcode, ifStmt *If) {
	parseExpresionWithNewScope(stack, &ifStmt.Condition)

	lenStack := len(*stack)
	label := "_if." + strconv.FormatInt(int64(lenStack), 16)

	*stack = append(*stack, &Opcode{"if", []any{"last_pop_exp", label}, nil})
	parseBody(stack, ifStmt.Body)

	*stack = append(*stack, &Opcode{"if_else", []any{}, &label})

}

func parseAssigment(stack *[]*Opcode, assigment *Assigment) {
	parseExpresionWithNewScope(stack, &assigment.Expression)
	*stack = append(*stack, &Opcode{"set_local_var_exp", []any{assigment.Variable.Value, "last_pop_exp"}, nil})

}

func parseFunctionCall(stack *[]*Opcode, functionCall *FunctionCall) {
	// todo check function declaration before making opcodes (like checking types of called function and numer of arguments)
	for _, argument := range functionCall.Arguments {
		parseExpresionWithNewScope(stack, argument)
		*stack = append(*stack, &Opcode{"push_function_arg", []any{"pop_exp"}, nil})
	}
	*stack = append(*stack, &Opcode{"call_function", []any{functionCall.FunctionName, len(functionCall.Arguments)}, nil})
}

func parseReturnStmt(stack *[]*Opcode, returnStmt *ReturnStmt) {
	parseExpresionWithNewScope(stack, &returnStmt.Expression)
	*stack = append(*stack, &Opcode{"push_bellow", []any{"last_pop_exp"}, nil})
	*stack = append(*stack, &Opcode{"function_return", []any{}, nil})
}

func parseExpresionWithNewScope(stack *[]*Opcode, expression *Expression) {
	*stack = append(*stack, &Opcode{"add_scope", []any{""}, nil})
	parseExpresion(stack, expression)
	*stack = append(*stack, &Opcode{"sub_scope", []any{""}, nil})
}

func parseExpresion(stack *[]*Opcode, expression *Expression) {
	parseComTerm(stack, expression.Left)
	parseRightComExpresion(stack, expression.Right)
}

func parseRightComExpresion(stack *[]*Opcode, opComTerm []*OpComTerm) {
	for _, opTerm := range opComTerm {
		parseComTerm(stack, opTerm.Term)
		*stack = append(*stack, &Opcode{"exp_call", []any{opTerm.Operator}, nil})
	}
}

func parseComTerm(stack *[]*Opcode, comTerm *ComTerm) {
	parseLeftTerm(stack, comTerm.Left)
	parseRightTerm(stack, comTerm.Right)
}

func parseRightTerm(stack *[]*Opcode, opTerm []*OpTerm) {
	parseOpTerm(stack, opTerm)
}

func parseOpTerm(stack *[]*Opcode, opTerms []*OpTerm) {
	for _, opTerm := range opTerms {
		parseTerm(stack, opTerm.Term)
		*stack = append(*stack, &Opcode{"exp_call", []any{opTerm.Operator}, nil})
	}
}

func parseLeftTerm(stack *[]*Opcode, term *Term) {
	parseTerm(stack, term)
}

func parseTerm(stack *[]*Opcode, term *Term) {
	parseFactor(stack, term.Left)
	parseOpFactor(stack, term.Right)
}

func parseOpFactor(stack *[]*Opcode, opFactors []*OpFactor) {
	for _, opFactor := range opFactors {
		parseFactor(stack, opFactor.Factor)
		*stack = append(*stack, &Opcode{"exp_call", []any{opFactor.Operator}, nil})
	}
}

func parseFactor(stack *[]*Opcode, factor *Factor) {
	if factor.Value != nil {
		if factor.Value.Float != nil {
			*stack = append(*stack, &Opcode{"push_exp", []any{factor.Value.Float.Value}, nil})
		} else if factor.Value.Integer != nil {
			*stack = append(*stack, &Opcode{"push_exp", []any{factor.Value.Integer.Value}, nil})
		} else if factor.Value.String != nil {
			*stack = append(*stack, &Opcode{"push_exp", []any{factor.Value.String.Value}, nil})
		} else if factor.Value.Boolean != nil {
			*stack = append(*stack, &Opcode{"push_exp", []any{factor.Value.Boolean.Value}, nil})
		}
	}
	if factor.FunctionCall != nil {
		parseFunctionCall(stack, factor.FunctionCall)
	}
	if factor.Variable != nil {
		*stack = append(*stack, &Opcode{"push_exp_var", []any{factor.Variable.Value}, nil})
	}
	if factor.Subexpression != nil {
		parseExpresion(stack, factor.Subexpression)
	}
}

func parseFunction(stack *[]*Opcode, function *Function) error {
	label := "_function." + function.Name

	for _, OpCode := range *stack {
		if OpCode.Label != nil && *OpCode.Label == label {
			return errors.New("function " + function.Name + " is already declared")
		}
	}

	*stack = append(*stack, &Opcode{"function", []any{}, &label})

	for _, argument := range function.Arguments {
		*stack = append(*stack, &Opcode{"set_local_var_arg", []any{argument.Name.Value, argument.Type}, nil})
	}

	err := parseFunctionBody(stack, function)
	if err != nil {
		return err
	}
	return nil
}

func GetOpcodes(code *Code) ([]*Opcode, error) {
	var opcodes []*Opcode

	for _, function := range code.Functions {
		err := parseFunction(&opcodes, function)
		if err != nil {
			return nil, err
		}
	}

	opcodes = append(opcodes, &Opcode{"call_function", []any{"main", 0}, nil}, &Opcode{"exit", []any{"main", 0}, nil})

	return opcodes, nil
}
