package test

import (
	karboscript "karboScript/src"
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
	ast, err := karboscript.ParseString("function main() { out(1000 + test(123) * 2 + 22, test(123)); } function test(int $test) { return $test + 200;}")
	
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
	ast, err := karboscript.ParseString("function test() { $a = 100; $aaa = 12 + $a; return $aaa;}function main() { out(test());}	")
	
	if err != nil {
	}

	opcodes, _ := karboscript.GetOpcodes(ast)
	_ = karboscript.Execute(&opcodes)

	// Output:
	// 112
}
