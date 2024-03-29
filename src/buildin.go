package karboscript

import (
	"bufio"
	"fmt"
	"os"
)

type buildInFunction func(program *Program) error

var buildInFunctions = map[string]buildInFunction{
	"out":      out,
	"readLine": readLine,
	"readInt":  readInt,
}

func out(program *Program) error {
	arguments := getFunctionArguments(program)
	fmt.Println(arguments...)

	return nil
}

func readLine(program *Program) error {
	getFunctionArguments(program)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')

	if err != nil {
		return nil
	}

	program.getScope(0).pushExp(text)
	return nil
}

func readInt(program *Program) error {
	getFunctionArguments(program)
	var out int

	fmt.Scanf("%d", &out)

	program.getScope(0).pushExp(out)
	return nil
}
