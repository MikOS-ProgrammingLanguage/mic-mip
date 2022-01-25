package compiler_util

import (
	"fmt"
	jsonconf "mik/mic_/json_conf"
	"os"
	"reflect"
	"strings"
)

var ditch_nil_ptr string = "// this is just to ditch a nil\n"
var ditch_nil_ptr2 string = "// this is just to ditch a nil\n" // i need that bcs otherwise i would write shit to the same mem address
var ditch_nil_ptr3 string = "// this is just to ditch a nil\n"

var ast_gc *RootNode
var current_node_gc Node = nil
var gc_pos int = -1
var is_eoast_gc bool = false
var in_func_gc bool = false
var arg_parse_gc bool = false
var main_code_gc *string = &ditch_nil_ptr3
var main_code_arr_gc []string
var none_main_code_gc *string = &ditch_nil_ptr
var global_vars_code_gc *string = &ditch_nil_ptr2

func gen_c_advance() {
	gc_pos++
	if gc_pos < len(ast_gc.Nodes) {
		current_node_gc = ast_gc.Nodes[gc_pos]
	} else {
		current_node_gc = UniversalNone{}
		is_eoast_gc = true
	}
}

func tconvert_c(t string) string {
	switch t {
	case "str":
		return "char*"
	case "flt":
		return "float"
	case "int64":
		return "long long"
	case "int32":
		return "int"
	case "int16":
		return "short"
	case "int8":
		return "char"
	case "uint64":
		return "unsigned long long"
	case "uint32":
		return "unsigned int"
	case "uint16":
		return "unsigned short"
	case "uint8":
		return "unsigned char"
	case "cock":
		return "unsigned volatile long long int"
	default:
		return t
	}
}

func gen_c(node Node) string {
	// switch through all first class nodes
	switch reflect.TypeOf(node).Name() {
	case "AssignemntNode":
		// make assignement
		cast_val := node.(AssignemntNode)
		code_ := tconvert_c(cast_val.Asgn_type)
		ptr_code := ""
		var scnd_glob string = ""

		for i := 0; i < cast_val.Ptrs; i++ {
			code_ += "*"
			ptr_code += "*"
		}
		code_ += " " + cast_val.Var_name
		if reflect.TypeOf(cast_val.Content).Name() != "UniversalNone" {
			if !cast_val.Global {
				code_ += "= "
				code_ += gen_c_s_l(cast_val.Content)
			} else {
				code_ += ""
				cntnt := gen_c_s_l(cast_val.Content) + ";\n"
				if in_func_gc {
					scnd_glob = ptr_code + cast_val.Var_name + "=" + cntnt
				} else {
					main_code_arr_gc = append(main_code_arr_gc, ptr_code+cast_val.Var_name+"="+cntnt)
				}
			}
		}

		if cast_val.Global {
			code_ += ";\n"
			*global_vars_code_gc += code_
			return scnd_glob
		} else if in_func_gc {
			code_ += ";\n\t"
			*none_main_code_gc += code_
			return code_
		} else if arg_parse_gc {
			return code_
		} else {
			code_ += ";\n"
			main_code_arr_gc = append(main_code_arr_gc, code_)
			return code_
		}
	case "ReAssignmentNode":
		// make reassignement
		cast_val := node.(ReAssignmentNode)
		code_ := ""
		if VARS[cast_val.Re_type].(AssignemntNode).Asgn_type == "str" {
			code_ += "strcpy("
			if cast_val.Ptrs > 0 {
				for i := 0; i < cast_val.Ptrs; i++ {
					code_ += "*"
				}
			}
			code_ += cast_val.Re_type
			code_ += ", "
			//code_ = strings.Replace(code_, "?=", "=", 1)
			code_ += gen_c_s_l(cast_val.Content)
			code_ += ");\n"
			if in_func_gc {
				return code_
			}
			main_code_arr_gc = append(main_code_arr_gc, code_)
			return code_
		} else {
			if cast_val.Ptrs > 0 {
				for i := 0; i < cast_val.Ptrs; i++ {
					code_ += "*"
				}
			}
			code_ += cast_val.Re_type + cast_val.Reassgn_t
			code_ = strings.Replace(code_, "?=", "=", 1)
			code_ += gen_c_s_l(cast_val.Content)
			code_ += ";\n"
			if in_func_gc {
				return code_
			}
			main_code_arr_gc = append(main_code_arr_gc, code_)
			return code_
		}
	case "ArrAssignementNode":
		// make a array assignement
		cast_val := node.(ArrAssignementNode)
		code_ := tconvert_c(cast_val.Asgn_type) + " "
		for i := 0; i < cast_val.Ptrs; {
			code_ += "*"
		}
		code_ += cast_val.Array_name + "["
		if cast_val.Arr_len != nil {
			code_ += gen_c_s_l(cast_val.Arr_len)
			code_ += "];\n"
		} else {
			code_ += "] = 0;\n"
		}

		if cast_val.Global {
			code_ += *none_main_code_gc
			none_main_code_gc = &code_
		} else if in_func_gc {
			code_ += ";\n\t"
			*none_main_code_gc += code_
			return code_
		} else if arg_parse_gc {
			return code_
		} else {
			main_code_arr_gc = append(main_code_arr_gc, code_)
		}
		return code_
	case "ArrReAssignementNode":
		// make a reassignement for a element in an array
		cast_val := node.(ArrReAssignementNode)
		code_ := ""
		for i := 0; i < cast_val.Ptrs; i++ {
			code_ += "*"
		}
		code_ += cast_val.Re_type + "[" + gen_c_s_l(cast_val.Arr_idx) + "]" + strings.Replace(cast_val.Reassgn_t, "?=", "=", 1) + gen_c_s_l(cast_val.Content) + ";\n"
		if in_func_gc {
			return code_
		}
		main_code_arr_gc = append(main_code_arr_gc, code_)
		return code_
	case "FuncCallNode":
		if arg_parse_gc {
			return gen_c_s_l(node.(FuncCallNode))
		}
		if in_func_gc {
			return gen_c_s_l(node.(FuncCallNode)) + ";\n"
		} else {
			main_code_arr_gc = append(main_code_arr_gc, gen_c_s_l(node.(FuncCallNode))+";\n")
			return gen_c_s_l(node.(FuncCallNode)) + ";\n"
		}
	case "FunctionNode":
		cast_val := current_node_gc.(FunctionNode)
		code_ := tconvert_c(cast_val.Ret_type) + " "

		code_ += cast_val.Func_name + "("
		arg_parse_gc = true
		if len(cast_val.Arg_parse) > 0 {
			for i := 0; i < (len(cast_val.Arg_parse)); i++ {
				code_ += gen_c(cast_val.Arg_parse[i])

				if (i + 1 != len(cast_val.Arg_parse)) {
					code_ += ","
				}
			}
		} else if len(cast_val.Arg_parse) == 0 {
			code_ += ""
		} else {
			temp := cast_val.Arg_parse[0]
			code_ += gen_c(temp)
		}
		arg_parse_gc = false
		code_ += ")"
		if cast_val.Decl {
			*none_main_code_gc += code_ + ";\n"
			return code_ + ";\n"
		}
		code_ += "{\n"
		in_func_gc = true
		none_main_code_gc_prev := none_main_code_gc
		for _, val := range cast_val.Code_block {
			code_ += gen_c(val)
		}
		in_func_gc = false
		code_ += "}\n"

		none_main_code_gc = none_main_code_gc_prev
		*none_main_code_gc += code_

		return code_
	case "ReturnNode":
		cast_val := node.(ReturnNode)
		code_ := "return "
		code_ += gen_c_s_l(cast_val.Return_val)
		code_ += ";\n"
		return code_
	case "AsmFunctionNode":
		// make as assembly function
		cast_val := node.(AsmFunctionNode)
		code_ := "__attribute__((naked)) "
		code_ += cast_val.Ret_type + " " + cast_val.Func_name + "("

		arg_parse_gc = true
		if len(cast_val.Arg_parse) > 0 {
			var len_ int = 0
			for i := 0; i < (len(cast_val.Arg_parse) - 1); i++ {
				len_ = i
				code_ += gen_c(cast_val.Arg_parse[i])
				code_ += ","
			}
			code_ += gen_c(cast_val.Arg_parse[len_+1])
		} else if len(cast_val.Arg_parse) == 0 {
			code_ += ""
		} else {
			temp := cast_val.Arg_parse[0]
			code_ += gen_c(temp)
		}
		arg_parse_gc = false
		code_ += ")" + " { __asm__ __volatile__ (" + cast_val.Asm_block + ");}\n"
		*none_main_code_gc += code_
		return code_
	case "StructNode":
		// make struct or estruct
		cast_val := node.(StructNode)
		code_ := ""

		if cast_val.Typedef {
			code_ += "typedef "
		}
		if cast_val.Estruct {
			code_ += "struct __attribute__((packed))"
			if cast_val.Typedef {
				code_ += ""
			} else {
				code_ += cast_val.Name
			}
			code_ += "{\n"
		} else {
			code_ += "struct "
			if cast_val.Typedef {
				code_ += ""
			} else {
				code_ += cast_val.Name
			}
			code_ += "{\n"
		}
		in_func_gc = true
		current_none_main_gc := *none_main_code_gc
		arg_parse_gc = true
		for _, val := range cast_val.Vars {
			code_ += "\t"
			code_ += gen_c(val) // + ";\n"
		}
		arg_parse_gc = false
		if cast_val.Typedef {
			code_ += "}" + cast_val.Name + ";\n"
		} else {
			code_ += "};\n"
		}
		*none_main_code_gc = current_none_main_gc
		in_func_gc = false
		*none_main_code_gc += code_
		return code_
	case "RefStructNode":
		cast_val := node.(RefStructNode)
		code := ""
		code += "struct "
		code += cast_val.struct_name

		for i := 0; i < cast_val.ptrs; i++ {
			code += "*"
		}
		code += cast_val.var_name + ";\n\t"
		if in_func_gc {
			*none_main_code_gc += code
			return code
		} else {
			main_code_arr_gc = append(main_code_arr_gc, code)
			return code
		}
	case "IfNode":
		cast_val := node.(IfNode)
		code_ := ""
		if cast_val.elif {
			code_ += "else if ("
		} else {
			code_ += "if ("
		}
		code_ += gen_c_s_l(cast_val.bool_) + ") {\n"
		current_main_gc := main_code_arr_gc
		for _, val := range cast_val.codeblock {
			code_ += "\t\t"
			code_ += gen_c(val)
		}
		code_ += "\t}\n"
		main_code_arr_gc = current_main_gc
		if in_func_gc {
			*none_main_code_gc += code_
			return code_
		} else {
			main_code_arr_gc = append(main_code_arr_gc, code_)
			return code_
		}
	case "ElseNode":
		cast_val := node.(ElseNode)
		code_ := "else {\n"
		current_main_gc := main_code_arr_gc
		for _, val := range cast_val.codeblock {
			code_ += "\t\t"
			code_ += gen_c(val)
		}
		code_ += "\t}\n"
		main_code_arr_gc = current_main_gc
		if in_func_gc {
			*none_main_code_gc += code_
			return code_
		} else {
			main_code_arr_gc = append(main_code_arr_gc, code_)
			return code_
		}
	case "WhileNode":
		cast_val := node.(WhileNode)
		code_ := "while ("
		code_ += gen_c_s_l(cast_val.bool_) + ") {\n"
		current_main_gc := main_code_arr_gc
		for _, val := range cast_val.codeblock {
			code_ += "\t\t"
			code_ += gen_c(val)
		}
		code_ += "\t}\n"
		main_code_arr_gc = current_main_gc
		if in_func_gc {
			*none_main_code_gc += code_
			return code_
		} else {
			main_code_arr_gc = append(main_code_arr_gc, code_)
			return code_
		}
	case "DirectNode":
		cast_val := node.(DirectNode)
		if cast_val.Type_ == "" {
			code_ := ""
			code_ += cast_val.Value
			code_ += ";\n"

			if in_func_gc {
				*none_main_code_gc += code_
				return code_
			} else {
				main_code_arr_gc = append(main_code_arr_gc, code_)
				return code_
			}
		} else {
			NewError("", "", "", true)
		}
	default:
		NewError(reflect.TypeOf(current_node_gc).Name(), "", "", true)
	}
	return ""
}

func GenerateC(ast_ *RootNode, out_pth, inptPtr string, conf jsonconf.Config) bool {
	if conf.CLibs != nil {
		lib_code, _, _ := jsonconf.MakeCLib(conf, inptPtr)
		*none_main_code_gc += *lib_code
	}

	ast_gc = ast_
	gen_c_advance()
	for !is_eoast_gc {
		gen_c(current_node_gc)
		gen_c_advance()
	}
	if len(main_code_arr_gc) > 0 {
		pre_code := "int main(void) {\n"
		main_code_gc = &pre_code

		for _, val := range main_code_arr_gc {
			*main_code_gc += "\t"
			*main_code_gc += val
		}
		*main_code_gc += "}"
	}
	data := *global_vars_code_gc + "char* strcpy(char* dest, const char* src) {do {*dest++ = *src++;} while (*src != 0);return 0;} // strcpy\n" + *none_main_code_gc + *main_code_gc
	err := os.WriteFile(out_pth, []byte(data), 0644)
	if err != nil {
		panic(err)
	} else {
		return SUCCESS
	}
}

func op_tok_gc(t string) string {
	switch t {
	case "b&":
		return "&"
	case "b|":
		return "|"
	case "b!":
		return "~"
	default:
		return t
	}
}

func gen_c_s_l(node LiteralNode) string {
	// switch all second and third class objects
	switch reflect.TypeOf(node).Name() {
	case "BoolOpNode":
		new_op := node.(BoolOpNode)
		left_c := gen_c_s_l(new_op.left)
		op := op_tok_gc(new_op.op_tok)
		right_c := gen_c_s_l(new_op.right)
		return left_c + op + right_c
	case "BinOpNode":
		new_op := node.(BinOpNode)
		left_c := gen_c_s_l(new_op.Left_node)
		op := op_tok_gc(new_op.Op_tok)
		right_c := gen_c_s_l(new_op.Right_node)
		return left_c + op + right_c
	case "DataTypeNode":
		new_op := node.(DataTypeNode)
		code_ := tconvert_c(new_op.Dtype)
		if new_op.Ptrs > 0 {
			for i := 0; new_op.Ptrs > i; i++ {
				code_ += "*"
			}
		}
		return code_
	case "FuncCallNode":
		new_op := node.(FuncCallNode)
		code_ := ""
		if new_op.BitNot {
			code_ += "~"
		} else if new_op.Minus {
			code_ += "-"
		}
		code_ += new_op.Call_name + "("
		if len(new_op.Func_parse) > 0 {
			var len_ int = 0
			for i := 0; i < (len(new_op.Func_parse) - 1); i++ {
				len_ = i
				code_ += gen_c_s_l(new_op.Func_parse[i])
				code_ += ","
			}
			code_ += gen_c_s_l(new_op.Func_parse[len_+1])
		} else if len(new_op.Func_parse) == 0 {
			code_ += ""
		} else {
			code_ += gen_c_s_l(new_op.Func_parse[0])
		}
		code_ += ")"
		return code_
	case "DirectNode":
		new_op := node.(DirectNode)
		code_ := ""
		if new_op.BitNot {
			code_ += "~"
		} else if new_op.Minus {
			code_ += "-"
		}
		if new_op.Type_ == TT_STRING {
			return code_ + fmt.Sprintf("\"%s\"", new_op.Value)
		} else if new_op.Type_ == TT_CHAR {
			return code_ + fmt.Sprintf("'%s'", new_op.Value)
		} else {
			return code_ + new_op.Value
		}
	case "VarNameNode":
		new_op := node.(VarNameNode)
		code_ := ""
		if new_op.Not {
			code_ += "!"
		}
		if new_op.BitNot {
			code_ += "~"
		} else if new_op.Minus {
			code_ += "-"
		}
		if new_op.Deref {
			code_ += "&"
		} else if new_op.Ptrs != 0 {
			for i := 0; i < new_op.Ptrs; i++ {
				code_ += "*"
			}
		}
		return code_ + new_op.Name
	case "ListSliceNode":
		new_op := node.(ListSliceNode)
		code_ := ""
		if new_op.Not {
			code_ += "!"
		}
		if new_op.BitNot {
			code_ += "~"
		} else if new_op.Minus {
			code_ += "-"
		}
		if new_op.Ptrs != 0 {
			for i := 0; i < new_op.Ptrs; i++ {
				code_ += "*"
			}
		} else if new_op.Deref {
			code_ += "&"
		}
		code_ += new_op.Name + "["
		code_ += gen_c_s_l(new_op.Pos) + "]"
		return code_
	case "TypeCastNode":
		new_op := node.(TypeCastNode)
		code_ := "("
		code_ += gen_c_s_l(new_op.Dtype) + ")"
		code_ += gen_c_s_l(new_op.Tcast)
		return code_
	default:
		return ""
	}
}
