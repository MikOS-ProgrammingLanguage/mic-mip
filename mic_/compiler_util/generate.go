package compiler_util

var SUCCESS bool = true
var FAILURE bool = false

func Generate(ast *RootNode, target, out_pth string) bool {
	switch target {
	case "c":
		return GenerateC(ast, out_pth)
	case "wasm":
		// generate web assembly
		return FAILURE
	case "asm":
		// generate and assemble branchless asm
		return FAILURE
	default:
		return FAILURE
	}
}
