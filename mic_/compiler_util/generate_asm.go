package compiler_util

import (
	"os"
	"reflect"
)

var ast_gasm *RootNode
var current_node_gasm Node = nil
var gasm_pos int = -1
var is_eoast_asm bool = false
var ditch_nil_ptr_asm string = "_start:\n\t"
var out_code *string = &ditch_nil_ptr

// advances to the next nodes
func gasm_advance() {
	gasm_pos++
	if gasm_pos < len(ast_gasm.Nodes) {
		current_node_gasm = ast_gasm.Nodes[gasm_pos]
	} else {
		current_node_gasm = nil
		is_eoast_asm = true
	}
}

func GenerateAsm(ast_ *RootNode, out_pth string) bool {
	ast_gasm = ast_
	gasm_advance()

	for !is_eoast_asm {
		gen_asm(current_node_gasm)
		gasm_advance()
	}

	err := os.WriteFile(out_pth, []byte(*out_code), 0755)
	if err != nil {
		panic(err)
	}

	return SUCCESS
}

// generate the fucking nodes
func gen_asm(node Node) {
	switch reflect.TypeOf(node).Name() {
	case "AssignemntNode":

		break
	case "":
		break
	default:
		break
	}
}
