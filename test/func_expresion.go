package test

import (
	karboscript "karboScript/src"
)

func ExampleExpresionTest() {
	ast, _ := karboscript.ParseString(`
		function main() {
			  out(12 + 14);
		}
	`)
	
	opcodes, _ := karboscript.GetOpcodes(ast)
	_ = karboscript.Execute(&opcodes)

	// Output:
	// 26
}
