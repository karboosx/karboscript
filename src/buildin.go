package karboscript

import (
	"fmt"
)

type buildInFunction func(program *Program) error

var buildInFunctions = map[string]buildInFunction{
	"out": out,
}

func out(program *Program) error {

	arguments := getFunctionArguments(program)
	fmt.Println(arguments...)

	return nil
}
