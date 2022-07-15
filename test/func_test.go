package test

import (
	karboscript "karboScript/src"
)

func ExampleFuncTest() {
	//ast, _ := karboscript.ParseString(`
	//	function main() {
	//		  out(a1(), a2());
	//	}
	//
	//	function a1() {
	//		  return 12;
	//	}
	//	function a2() {
	//		  return 55;
	//	}
	//`)
	//

	ast, _ := karboscript.Parse("testFunc.ks")
	opcodes, _ := karboscript.GetOpcodes(ast)
	_ = karboscript.Execute(&opcodes)

	// Output:
	// 12 55
}
