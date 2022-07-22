package test

import (
	karboscript "karboScript/src"
)

func ExampleFuncTest() {
	ast, _ := karboscript.Parse("testFunc.ks")
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
