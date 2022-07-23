package karboscript

import (
	"os"

	"github.com/alecthomas/participle/v2"

	"github.com/alecthomas/participle/v2/lexer"
)

type Code struct {
	Functions []*Function `@@*`
}

type Function struct {
	Name      string       `"function" @Ident "("`
	Arguments []*Argument  ` [@@ ("," @@)*] ")"`
	Body      []*Statement `"{" @@* "}"`
}

type Statement struct {
	If           *If           `(@@ `
	For          *For          `| @@ `
	While        *While        `| @@ ) | `
	Assigment    *Assigment    `(@@ `
	FunctionCall *FunctionCall `| @@`
	Expression   *Expression   `| @@`
	ReturnStmt   *ReturnStmt   `| @@) ";"`
}

type Assigment struct {
	Variable   Variable   `@@`
	Expression Expression `"=" @@`
}

type ReturnStmt struct {
	Expression Expression `"return" @@`
}

type FunctionCall struct {
	FunctionName string        `@Ident "("`
	Arguments    []*Expression ` [@@ ("," @@)*] ")"`
}

type If struct {
	Condition Expression   `"if" @@`
	Body      []*Statement `"{" @@* "}"`
}

type While struct {
	Condition Expression   `"while" @@`
	Body      []*Statement `"{" @@* "}"`
}
type For struct {
	Init      Statement    `"for" "("? @@`
	Condition Expression   `@@ ";"`
	Increment Statement    `@@ ")"?`
	Body      []*Statement `"{" @@* "}"`
}

type Argument struct {
	Variable Variable `@@`
}

type Value struct {
	Integer *Integer `@@`
	Float   *Float   `| @@`
	Boolean *Boolean `| @@`
	String  *String  `| @@`
}

type String struct {
	Value string `@String`
}

type Integer struct {
	Value int `@Int`
}

type Float struct {
	Value float64 `@Float`
}

type Boolean struct {
	Value string `@("true"|"false")`
}

type Variable struct {
	Value string `"$" @Ident`
}

type Factor struct {
	FunctionCall  *FunctionCall `(@@`
	Value         *Value        `| @@`
	Subexpression *Expression   `| "(" @@ ")"`
	Variable      *Variable     `| @@)`
}

type OpFactor struct {
	Operator string `@("*" | "/")`
	Factor   *Factor  `@@`
}

type Term struct {
	Left  *Factor     `@@`
	Right []*OpFactor `@@*`
}

type OpTerm struct {
	Operator string `@("+" | "-")`
	Term     *Term    `@@`
}

type ComTerm struct {
	Left  *Term     `@@`
	Right []*OpTerm `@@*`
}

type OpComTerm struct {
	Operator string `@("=""=" | "!""=" | ">" | "<" | ">""=" | "<""=")`
	Term     *ComTerm `@@`
}

type Expression struct {
	Left  *ComTerm     `@@`
	Right []*OpComTerm `@@*`
}

var (
	karboScriptLexer = lexer.NewTextScannerLexer(nil)

	Parser = participle.MustBuild[Code](
		participle.Lexer(karboScriptLexer),
		participle.Elide("Comment"),
		participle.UseLookahead(2),
	)
)

func Parse(file string) (*Code, error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	ast, err := Parser.Parse(file, r)
	if err != nil {
		return nil, err
	}

	err = r.Close()
	if err != nil {
		return nil, err
	}

	return ast, nil
}

func GetTokens(file string) (lexer.Lexer, map[string]lexer.TokenType, error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}

	ast, err := karboScriptLexer.Lex(file, r)
	if err != nil {
		return nil, nil, err
	}

	symbols := karboScriptLexer.Symbols()

	if err != nil {
		return nil, nil, err
	}

	return ast, symbols, nil
}

func ParseString(code string) (*Code, error) {
	ast, err := Parser.ParseString("", code)
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
