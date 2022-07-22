package main

import (
	"fmt"
	"strconv"

	karboscript "karboScript/src"

	"github.com/alecthomas/kong"
	"github.com/alecthomas/repr"
)

var cli struct {
	EBNF   bool   `help:"Display DBNF."`
	Opcode bool   `help:"Display DBNF."`
	File   string `arg:"" optional:"" type:"existingfile" help:"GraphQL schema files to parse."`
}

var ctx kong.Context

func main() {

	ctx := kong.Parse(&cli)

	if cli.EBNF {
		fmt.Println(karboscript.Parser.String())
		ctx.Exit(0)
	}

	ast, err := karboscript.Parse(cli.File)

	opcodes, err := karboscript.GetOpcodes(ast)
	ctx.FatalIfErrorf(err)

	if cli.Opcode {
		var str string

		for _, opcode := range opcodes {
			str = ""
			if opcode.Label != nil {
				str = str + *opcode.Label + ": "
			}
			str = str + opcode.Operation

			if len(opcode.Arguments) > 0 {
				str = str + " ("
				for _, argument := range opcode.Arguments {

					if argstr, ok := argument.(string); ok {
						str = str + " " + argstr
					} else if argint, ok := argument.(int); ok {
						str = str + " " + strconv.FormatInt(int64(argint), 10)
					} else if argfloat, ok := argument.(float64); ok {
						str = str + " " + strconv.FormatFloat(argfloat, 'f', 0, 6)
					} else if argbool, ok := argument.(bool); ok {
						str = str + " " + strconv.FormatBool(argbool)
					}

				}
				str = str + " )"
			}
			repr.Println(str)
		}

		ctx.Exit(0)
	}

	err = karboscript.Execute(&opcodes)
	ctx.FatalIfErrorf(err)
}
