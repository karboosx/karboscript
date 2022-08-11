package test

import (
	karboscript "karboScript/src"

	"fmt"
)

func ExampleFuncTest() {
	ast, _ := karboscript.ParseString("function main() { out(a1(), a2()); }function a1() { return 12; }function a2() { return 55; }")
	opcodes, _ := karboscript.GetOpcodes(ast)
	_ = karboscript.Execute(&opcodes)

	// Output:
	// 12 55
}


func ExampleExpresionTest() {
	ast, _ := karboscript.ParseString("function main() { out(12 + 14); }")
	
	opcodes, _ := karboscript.GetOpcodes(ast)
	_ = karboscript.Execute(&opcodes)

	// Output:
	// 26
}

func ExampleExpresionWithFuncTest() {
	ast, err := karboscript.ParseString("function main() { out(1000 + test() * 2 + 22); } function test() { return 100;}")
	
	if err != nil {
	}

	opcodes, _ := karboscript.GetOpcodes(ast)
	_ = karboscript.Execute(&opcodes)

	// Output:
	// 1222
}
func ExampleFuncWithArgTest() {
	ast, err := karboscript.ParseString("function main() { out(1000 + test(123) * 2 + 22, test(123)); } function test(int test) { return test + 200;}")
	
	if err != nil {
	}

	opcodes, _ := karboscript.GetOpcodes(ast)
	_ = karboscript.Execute(&opcodes)

	// Output:
	// 1668 323
}
func ExampleCompareTest() {
	ast, err := karboscript.ParseString("function main() { out(12 > 10, 10 == 10, 30 == 10, 10 != 10); }")
	
	if err != nil {
	}

	opcodes, _ := karboscript.GetOpcodes(ast)
	_ = karboscript.Execute(&opcodes)

	// Output:
	// true true false false
}
func ExampleLocalVarTest() {
	ast, err := karboscript.ParseString("function test() { int a = 100; int aaa = 12 + a; return aaa;} function main() { out(test());}	")
	
	if err != nil {
	}

	opcodes, _ := karboscript.GetOpcodes(ast)
	_ = karboscript.Execute(&opcodes)

	// Output:
	// 112
}
func ExampleLocalVarKeepsLocalScopeTest() {
	ast, err := karboscript.ParseString("function test() { out(a); }function main() { int a = 10; test();}	")
	
	if err != nil {
	}

	opcodes, _ := karboscript.GetOpcodes(ast)
	err = karboscript.Execute(&opcodes)

	fmt.Println(err)
		
	// Output:
	// Undeclared variable: a
}

func ExampleIfTest() {
	ast, err := karboscript.ParseString("function main() {    if (10 == 10) {        out(\"10 == 10\");    }    if (500 < 200) {        out(\"500 < 200\");    }    if (12 > 10) {        out(\"12 > 10\");    }}	")
	
	if err != nil {
	}

	opcodes, _ := karboscript.GetOpcodes(ast)
	_ = karboscript.Execute(&opcodes)

	// Output:
	// 10 == 10
	// 12 > 10
}
func ExampleWhileTest() {
	ast, err := karboscript.ParseString("function main() {    int a = 0;    while (a < 10) {        out (a);        a = a + 1;    }}")
	
	if err != nil {
	}

	opcodes, _ := karboscript.GetOpcodes(ast)
	_ = karboscript.Execute(&opcodes)

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
}

func ExampleDoubleWhileTest() {
	ast, err := karboscript.ParseString("function main() {    int a = 1;    int b = 1;    while (a < 3) {        b = 1;        while (b < 3) {            out (a, b);            b=b+1;        }        a=a+1;    }}")
	
	if err != nil {
	}

	opcodes, _ := karboscript.GetOpcodes(ast)
	_ = karboscript.Execute(&opcodes)

	// Output:
	// 1 1
	// 1 2
	// 2 1
	// 2 2
}
func ExampleForTest() {
	ast, err := karboscript.ParseString("function main() {    for int i=0; i<10; i=i+1; {        out(i);    }}")
	
	if err != nil {
	}

	opcodes, _ := karboscript.GetOpcodes(ast)
	_ = karboscript.Execute(&opcodes)

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
}

func ExampleForIncTest() {
	ast, err := karboscript.ParseString("function main() {    from 0 to 10 as i {        out(i);    }}")
	
	if err != nil {
	}

	opcodes, _ := karboscript.GetOpcodes(ast)
	_ = karboscript.Execute(&opcodes)

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
}
