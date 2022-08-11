package karboscript

import (
	"errors"
	"strconv"
	"strings"
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

type ParsedCode struct {
	functions   map[string]Function
	stack       *[]*Opcode
	parsedError error
}

func (parsed *ParsedCode) append(opcode *Opcode) {
	*(*parsed).stack = append(*(*parsed).stack, opcode)
}

func parseFunctionBody(parsed *ParsedCode, function *Function) error {
	parseBody(parsed, function.Body)
	if (*(*parsed).stack)[len(*(*parsed).stack)-1].Operation != "function_return" {
		parsed.append(&Opcode{"function_return", []any{}, nil})
	}
	return nil
}

func parseBody(parsed *ParsedCode, statements []*Statement) error {
	for _, statement := range statements {
		parseStatement(parsed, statement)
	}

	return nil
}

func parseStatement(parsed *ParsedCode, statement *Statement) {
	if statement.FunctionCall != nil {
		parseFunctionCall(parsed, statement.FunctionCall)
	}
	if statement.Expression != nil {
		parseExpresionWithNewScope(parsed, statement.Expression)
	}
	if statement.ReturnStmt != nil {
		parseReturnStmt(parsed, statement.ReturnStmt)
	}
	if statement.Assigment != nil {
		parseAssigment(parsed, statement.Assigment)
	}
	if statement.If != nil {
		parseIf(parsed, statement.If)
	}
	if statement.While != nil {
		parseWhile(parsed, statement.While)
	}
	if statement.For != nil {
		parseFor(parsed, statement.For)
	}
	if statement.ForInc != nil {
		parseForInc(parsed, statement.ForInc)
	}
}

func parseWhile(parsed *ParsedCode, while *While) {
	labelBeforeExpresion := newLabel(parsed, "while")
	parsed.append(&Opcode{"while_start", []any{}, &labelBeforeExpresion})

	parseExpresionWithNewScope(parsed, &while.Condition)

	label := newLabel(parsed, "while")

	parsed.append(&Opcode{"while", []any{"last_pop_exp", label}, nil})
	parseBody(parsed, while.Body)
	parsed.append(&Opcode{"jmp", []any{labelBeforeExpresion}, nil})

	parsed.append(&Opcode{"while_else", []any{}, &label})
}

func parseFor(parsed *ParsedCode, forStmt *For) {
	parseStatement(parsed, &forStmt.Init)

	labelBeforeExpresion := newLabel(parsed, "for")
	parsed.append(&Opcode{"for_start", []any{}, &labelBeforeExpresion})

	parseExpresionWithNewScope(parsed, &forStmt.Condition)

	label := newLabel(parsed, "for")

	parsed.append(&Opcode{"for", []any{"last_pop_exp", label}, nil})
	parseBody(parsed, forStmt.Body)
	parseStatement(parsed, &forStmt.Increment)

	parsed.append(&Opcode{"jmp", []any{labelBeforeExpresion}, nil})

	parsed.append(&Opcode{"for_end", []any{}, &label})
}

func parseForInc(parsed *ParsedCode, forStmt *ForInc) {
	parseExpresionWithNewScope(parsed, &forStmt.ExpressionA)

	parsed.append(&Opcode{"set_local_var_exp", []any{"int", forStmt.Variable.Value, "last_pop_exp"}, nil})
	
	parseExpresionWithNewScope(parsed, &forStmt.ExpressionB)

	parsed.append(&Opcode{"set_local_var_exp", []any{"int", forStmt.Variable.Value+"_end", "last_pop_exp"}, nil})

	incLabelStart := newLabel(parsed, "forinc")
	incLabelEnd := newLabel(parsed, "forinc_e")
	parsed.append(&Opcode{"forinc_start", []any{forStmt.Variable.Value, incLabelEnd}, &incLabelStart})

	parseBody(parsed, forStmt.Body)

	parsed.append(&Opcode{"forinc", []any{forStmt.Variable.Value, incLabelStart}, nil})

	parsed.append(&Opcode{"forinc_end", []any{}, &incLabelEnd})
}

func newLabel(parsed *ParsedCode, labelType string) string {
	lenStack := len(*(*parsed).stack)
	label := "_" + labelType + "." + strconv.FormatInt(int64(lenStack), 16)
	return label
}

func parseIf(parsed *ParsedCode, ifStmt *If) {
	parseExpresionWithNewScope(parsed, &ifStmt.Condition)

	lenStack := len(*(*parsed).stack)
	label := "_if." + strconv.FormatInt(int64(lenStack), 16)

	parsed.append(&Opcode{"if", []any{"last_pop_exp", label}, nil})
	parseBody(parsed, ifStmt.Body)

	parsed.append(&Opcode{"if_else", []any{}, &label})

}

func parseAssigment(parsed *ParsedCode, assigment *Assigment) {
	parseExpresionWithNewScope(parsed, &assigment.Expression)
	parsed.append(&Opcode{"set_local_var_exp", []any{assigment.VarType.Value, assigment.Variable.Value, "last_pop_exp"}, nil})
}

func parseFunctionCall(parsed *ParsedCode, functionCall *FunctionCall) {
	// todo check function declaration before making opcodes (like checking types of called function and numer of arguments)
	for _, argument := range functionCall.Arguments {
		parseExpresionWithNewScope(parsed, argument)
		parsed.append(&Opcode{"push_function_arg", []any{"last_pop_exp"}, nil})
	}

	if function, ok := parsed.functions[functionCall.FunctionName]; ok {
		if function.ReturnType != nil {
			parsed.append(&Opcode{"call_function", []any{functionCall.FunctionName, len(functionCall.Arguments), function.ReturnType.Value}, nil})
			return;
		} else {
			parsed.append(&Opcode{"call_function", []any{functionCall.FunctionName, len(functionCall.Arguments)}, nil})
			return;
		}
	}
	
	if _, ok := buildInFunctions[functionCall.FunctionName]; ok {
		parsed.append(&Opcode{"call_function", []any{functionCall.FunctionName, len(functionCall.Arguments)}, nil})
		
	}else {
		parsed.parsedError = errors.New("Can't find " + functionCall.FunctionName + " function!")
	}
}

func parseReturnStmt(parsed *ParsedCode, returnStmt *ReturnStmt) {
	parseExpresionWithNewScope(parsed, &returnStmt.Expression)
	parsed.append(&Opcode{"push_bellow", []any{"last_pop_exp"}, nil})
	parsed.append(&Opcode{"function_return", []any{}, nil})
}

func parseExpresionWithNewScope(parsed *ParsedCode, expression *Expression) {
	parsed.append(&Opcode{"add_scope", []any{}, nil})
	parseExpresion(parsed, expression)
	parsed.append(&Opcode{"sub_scope", []any{}, nil})
}

func parseExpresion(parsed *ParsedCode, expression *Expression) {
	parseComTerm(parsed, expression.Left)
	parseRightComExpresion(parsed, expression.Right)
}

func parseRightComExpresion(parsed *ParsedCode, opComTerm []*OpComTerm) {
	for _, opTerm := range opComTerm {
		parseComTerm(parsed, opTerm.Term)
		parsed.append(&Opcode{"exp_call", []any{opTerm.Operator}, nil})
	}
}

func parseComTerm(parsed *ParsedCode, comTerm *ComTerm) {
	parseLeftTerm(parsed, comTerm.Left)
	parseRightTerm(parsed, comTerm.Right)
}

func parseRightTerm(parsed *ParsedCode, opTerm []*OpTerm) {
	parseOpTerm(parsed, opTerm)
}

func parseOpTerm(parsed *ParsedCode, opTerms []*OpTerm) {
	for _, opTerm := range opTerms {
		parseTerm(parsed, opTerm.Term)
		parsed.append(&Opcode{"exp_call", []any{opTerm.Operator}, nil})
	}
}

func parseLeftTerm(parsed *ParsedCode, term *Term) {
	parseTerm(parsed, term)
}

func parseTerm(parsed *ParsedCode, term *Term) {
	parseFactor(parsed, term.Left)
	parseOpFactor(parsed, term.Right)
}

func parseOpFactor(parsed *ParsedCode, opFactors []*OpFactor) {
	for _, opFactor := range opFactors {
		parseFactor(parsed, opFactor.Factor)
		parsed.append(&Opcode{"exp_call", []any{opFactor.Operator}, nil})
	}
}

func parseFactor(parsed *ParsedCode, factor *Factor) {
	if factor.Value != nil {
		if factor.Value.Float != nil {
			parsed.append(&Opcode{"push_exp", []any{factor.Value.Float.Value}, nil})
		} else if factor.Value.Integer != nil {
			parsed.append(&Opcode{"push_exp", []any{factor.Value.Integer.Value}, nil})
		} else if factor.Value.String != nil {
			stripSlash := strings.ReplaceAll(factor.Value.String.Value, "\\\"", "\"")
			parsed.append(&Opcode{"push_exp", []any{stripSlash[1 : len(stripSlash)-1]}, nil})
		} else if factor.Value.Boolean != nil {
			parsed.append(&Opcode{"push_exp", []any{factor.Value.Boolean.Value}, nil})
		}
	}
	if factor.FunctionCall != nil {
		parseFunctionCall(parsed, factor.FunctionCall)
	}
	if factor.Variable != nil {
		parsed.append(&Opcode{"push_exp_var", []any{factor.Variable.Value}, nil})
	}
	if factor.Subexpression != nil {
		parseExpresion(parsed, factor.Subexpression)
	}
}

func parseFunction(parsed *ParsedCode, function *Function) error {
	label := "_function." + function.Name

	for _, OpCode := range *(*parsed).stack {
		if OpCode.Label != nil && *OpCode.Label == label {
			return errors.New("function " + function.Name + " is already declared")
		}
	}

	*(*parsed).stack = append(*(*parsed).stack, &Opcode{"function", []any{}, &label})

	for _, argument := range function.Arguments {
		*(*parsed).stack = append(*(*parsed).stack, &Opcode{"set_local_var_arg", []any{argument.VarType.Value, argument.Variable.Value}, nil})
	}

	err := parseFunctionBody(parsed, function)
	if err != nil {
		return err
	}
	return nil
}

func registerFunction(parsed *ParsedCode, function *Function) {
	parsed.functions[function.Name] = *function
}

func GetOpcodes(code *Code) ([]*Opcode, error) {
	parsed := ParsedCode{map[string]Function{}, &[]*Opcode{}, nil}

	var opcodes []*Opcode

	for _, function := range code.Functions {
		registerFunction(&parsed, function)
	}

	for _, function := range code.Functions {
		err := parseFunction(&parsed, function)
		if err != nil {
			return nil, err
		}
	}

	if (parsed.parsedError != nil) {
		return opcodes, parsed.parsedError;
	}
	opcodes = append(*parsed.stack, &Opcode{"call_function", []any{"main", 0}, nil}, &Opcode{"exit", []any{"main", 0}, nil})

	return opcodes, nil
}
