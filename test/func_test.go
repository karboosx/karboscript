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
