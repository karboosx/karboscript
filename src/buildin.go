package main

import (
	"fmt"
)

type buildInFunction func(program *Program) error

var buildInFunctions = map[string]buildInFunction{
	"out": out,
}

func out(program *Program) error {
	fmt.Print(*program.expresionOutout)
	return nil
}
