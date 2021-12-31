package compiler_util

import jsonconf "mik/mic_/json_conf"

var SUCCESS bool = true
var FAILURE bool = false

func Generate(ast *RootNode, target, out_pth, inptPtr string, conf jsonconf.Config) bool {
	switch target {
	case "c":
		return GenerateC(ast, out_pth, inptPtr, conf)
	case "wasm":
		// generate web assembly
		return FAILURE
	case "asm":
		return GenerateAsm(ast, out_pth)
	default:
		return FAILURE
	}
}
