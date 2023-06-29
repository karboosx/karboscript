package karboscript

import (
	"os"
	"text/scanner"

	"github.com/alecthomas/participle/v2"

	"github.com/alecthomas/participle/v2/lexer"
)

type Code struct {
	Functions []*Function `@@*`
}

type Function struct {
	Pos lexer.Position

	Name       string       `"function" @Ident "("`
	Arguments  []*Argument  ` [@@ ("," @@)*] ")"`
	ReturnType *VarType     `@@?`
	Body       []*Statement `"{" @@* "}"`
}

type Statement struct {
	Pos lexer.Position

	If                *If                `(@@ `
	For               *For               `| @@ `
	ForInc            *ForInc            `| @@ `
	While             *While             `| @@ ) | `
	ReturnStmt        *ReturnStmt        `( @@ `
	ArrayAssigment    *ArrayAssigment    `| @@ `
	Assigment         *Assigment         `| @@ `
	FunctionCall      *FunctionCall      `| @@`
	Expression        *Expression        `| @@) ";"`
}

type VarType struct {
	Value string `@("array" | "string" | "int" | "float" | "bool")`
}

type Assigment struct {
	Pos lexer.Position

	VarType    VarType    `@@?`
	Variable   Variable   `@@`
	Expression Expression `"=" @@`
}

type ArrayAssigment struct {
	Pos lexer.Position

	Variable   Variable   `@@ "["`
	Index      Expression `@@ "]"`
	Expression Expression `"=" @@`
}

type ReturnStmt struct {
	Pos lexer.Position

	Expression Expression `"return" @@`
}

type FunctionCall struct {
	Pos lexer.Position

	FunctionName string        `@Ident "("`
	Arguments    []*Expression ` [@@ ("," @@)*] ")"`
}

type If struct {
	Pos lexer.Position

	Condition Expression   `"if" @@`
	Body      []*Statement `"{" @@* "}"`
}

type While struct {
	Pos lexer.Position

	Condition Expression   `"while" @@`
	Body      []*Statement `"{" @@* "}"`
}
type For struct {
	Pos lexer.Position

	Init      Statement    `"for" "("? @@`
	Condition Expression   `@@ ";"`
	Increment Statement    `@@ ")"?`
	Body      []*Statement `"{" @@* "}"`
}

type ForInc struct {
	Pos lexer.Position

	ExpressionA Expression   `"from" @@`
	ExpressionB Expression   `"to" @@`
	Variable    Variable     `"as" @@`
	Body        []*Statement `"{" @@* "}"`
}

type Argument struct {
	Pos lexer.Position

	VarType  VarType  `@@`
	Variable Variable `@@`
}

type Value struct {
	Pos lexer.Position

	Integer *Integer `@@`
	Float   *Float   `| @@`
	Boolean *Boolean `| @@`
	String  *String  `| @@`
}

type String struct {
	Pos lexer.Position

	Value string `@String`
}

type Integer struct {
	Pos lexer.Position

	Value int `@Int`
}

type Float struct {
	Pos lexer.Position

	Value float64 `@Float`
}

type Boolean struct {
	Pos lexer.Position

	Value string `@("true"|"false")`
}

type Variable struct {
	Pos lexer.Position

	Value string `@Ident`
}

type ArrayCall struct {
	Pos lexer.Position

	Name  string      `@Ident "["`
	Index *Expression `@@ "]"`
}

type ArrayLiteral struct {
	Pos lexer.Position

	Elements []*Expression `"[" [@@ ("," @@)*] "]"`
}

type Factor struct {
	Pos lexer.Position

	ArrayCall     *ArrayCall    `(@@`
	FunctionCall  *FunctionCall `| @@`
	Value         *Value        `| @@`
	Subexpression *Expression   `| "(" @@ ")"`
	Variable      *Variable     `| @@`
	ArrayLiteral  *ArrayLiteral `| @@)`
}

type OpFactor struct {
	Pos lexer.Position

	Operator string  `@("*" | "/")`
	Factor   *Factor `@@`
}

type Term struct {
	Pos lexer.Position

	Left  *Factor     `@@`
	Right []*OpFactor `@@*`
}

type OpTerm struct {
	Pos lexer.Position

	Operator string `@("+" | "-")`
	Term     *Term  `@@`
}

type ComTerm struct {
	Pos lexer.Position

	Left  *Term     `@@`
	Right []*OpTerm `@@*`
}

type OpComTerm struct {
	Pos lexer.Position

	Operator string   `@("=""=" | "!""=" | ">" | "<" | ">""=" | "<""=")`
	Term     *ComTerm `@@`
}

type Expression struct {
	Pos lexer.Position

	Left  *ComTerm     `@@`
	Right []*OpComTerm `@@*`
}

var (
	karboScriptLexer = lexer.NewTextScannerLexer(func(s *scanner.Scanner) {
		s.Mode &^= scanner.ScanChars
	})

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
