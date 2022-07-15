package main

import (
	"os"

	"github.com/alecthomas/participle/v2"

	"github.com/alecthomas/participle/v2/lexer"
)

type Code struct {
	Functions []*Function `@@*`
}

type Function struct {
	Name      string       `"function" @Symbol "("`
	Arguments []*Argument  `@@*")"`
	Body      []*Statement `"{" @@* "}"`
}

type Statement struct {
	Declaration  *Declaration  `(@@ `
	Expression   *Expression   `| @@`
	FunctionCall *FunctionCall `| @@`
	ReturnStmt   *ReturnStmt   `| @@) ";"`
}

type Declaration struct {
	Variable   string     `"var" @Symbol`
	Expression Expression `"=" @@`
}

type ReturnStmt struct {
	Expression Expression `"return" @@`
}

type FunctionCall struct {
	FunctionName string        `@Symbol "("`
	Arguments    []*Expression ` [@@ ("," @@)*] ")"`
}

type Argument struct {
	Type string `@Type`
	Name string `@Symbol`
}

type Value struct {
	String  *String  `@@`
	Integer *Integer `| @@`
	Float   *Float   `| @@`
	Boolean *Boolean `| @@`
}

type String struct {
	Value string `"\"" @Symbol "\""`
}

type Integer struct {
	Value int `@Integer`
}

type Float struct {
	Value float64 `@Float`
}

type Boolean struct {
	Value int `@Boolean`
}

type Expression struct {
	Value        *Value        `(@@`
	FunctionCall *FunctionCall `| @@)`
}

var LexerRules = []lexer.SimpleRule{
	{"Type", `string|int|boolean|float`},
	{"Boolean", `true|false`},
	{"Comment", `(?:#|//)[^\n]*\n?`},
	{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
	{"Symbol", `[a-zA-Z^\n]\w*`},
	{"Float", `\d+\.\d+`},
	{"Integer", `[0-9^\.]+`},
	{"Whitespace", `[ \t\n\r]+`},
}

var (
	karboScriptLexer = lexer.MustSimple(LexerRules)
	Parser           = participle.MustBuild[Code](
		participle.Lexer(karboScriptLexer),
		participle.Elide("Comment", "Whitespace"),
		participle.UseLookahead(4),
	)
)

func parse(file string) (*Code, error) {
	r, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	ast, err := Parser.Parse(file, r)

	r.Close()

	if err != nil {
		return nil, err
	}

	return ast, nil
}

// func ggg() {
// 	ctx := kong.Parse(&cli)
// 	if cli.EBNF {
// 		fmt.Println(parser.String())
// 		ctx.Exit(0)
// 	}
// 	for _, file := range cli.Files {
// 		r, err := os.Open(file)
// 		ctx.FatalIfErrorf(err)
// 		if cli.Tokens {
// 			lexer, err := parser.Lexer().Lex(file, r)

// 			run := true
// 			for run == true {
// 				n, err := lexer.Next()

// 				if n.Value == "" {
// 					run = false
// 				}

// 				if len(strings.TrimSpace(n.Value)) == 0 {
// 					continue
// 				}
// 				repr.Print(LexerRules[-n.Type-2].Name)

// 				repr.Println(n.Value)

// 				if err != nil {
// 					run = false
// 				}
// 			}

// 			ctx.FatalIfErrorf(err)
// 		}

// 		ast, err := parser.Parse(file, r)
// 		r.Close()
// 		repr.Println(ast)
// 		ctx.FatalIfErrorf(err)
// 	}
// }
