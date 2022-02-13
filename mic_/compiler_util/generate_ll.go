package compiler_util

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

var current_function_variable_count int = 1 // counts the temporary variables specified within a function
var main_function_variable_count int = 1    // counts the temporary variables specified within the main function

var var_is_refferenced bool = false

var VARS_LLVM []string

func get_align(ptrs int) string {
	fmt.Println(ptrs + 1)
	return fmt.Sprintf(", align %d", 4*(ptrs+1))
}
func get_ptrs(ptrs int) string {
	str := ""
	for i := 0; i < ptrs; i++ {
		str += "*"
	}
	return str
}

var current_ptr_offset int = 0
var ast_gasm *RootNode
var current_node_gasm Node = nil
var gasm_pos int = -1
var is_eoast_asm bool = false
var ditch_nil_ptr_asm string = ""
var ditch_nil_ptr_asm_2 string = ""
var string_constant string = " private unnamed_addr constant "
var out_code *string = &ditch_nil_ptr_asm
var global_vars *string = &ditch_nil_ptr_asm_2

func operator_cast_llvm(i string) string {
	switch i {
	case "+":
		return "add"
	}
	return ""
}

func type_convert_llvm(i string) string {
	switch i {
	case "int":
		return "i32"
	case "str":
		return "i8*"
	}
	return i
}

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

	*global_vars += "; Module: main.ll\n"
	*global_vars += "target datalayout = \"e-m:o-i64:64-i128:128-n32:64-S128\"\n\n"
	*out_code += "\ndefine i32 @main() {\n"
	for !is_eoast_asm {
		generate_asm(current_node_gasm)
		gasm_advance()
	}

	*out_code += "\tret i32 1\n}"
	*out_code = *global_vars + *out_code
	err := os.WriteFile(out_pth, []byte(*out_code), 0755)
	if err != nil {
		panic(err)
	}

	return SUCCESS
}

func generate_second_and_third_class_node(node LiteralNode, type_ string) string {
	switch reflect.TypeOf(node).Name() {
	case "BinOpNode":

		new_op := node.(BinOpNode)

		left_code := type_ + " " + generate_second_and_third_class_node(new_op.Left_node, type_)
		operator := operator_cast_llvm(new_op.Op_tok)
		right_code := generate_second_and_third_class_node(new_op.Right_node, type_)

		code := fmt.Sprintf("\t%%%d = ", current_function_variable_count)
		code += operator + " "
		code += left_code + ", "
		code += right_code + "\n"

		*out_code += code
		current_function_variable_count++
		return fmt.Sprintf("\t%%%d", current_function_variable_count-1)

	case "DirectNode":
		new_op := node.(DirectNode)
		return new_op.Value
	case "VarNameNode":
		var_is_refferenced = false
		new_op := node.(VarNameNode)

		// return "type %var_name" if the current var node is a var elese return the literal
		if StringInSlice(new_op.Name, VARS_LLVM) {
			var_is_refferenced = true
			code := fmt.Sprintf("\t%%%d = load %s, %s%s %%%s %s\n", current_function_variable_count, type_, type_, get_ptrs(new_op.Ptrs+1), new_op.Name, get_align(new_op.Ptrs+current_ptr_offset))

			*out_code += code
			current_function_variable_count++

			return fmt.Sprintf("\t%%%d", (current_function_variable_count - 1))
		} else {
			if type_ == "i8*" {
				return strings.ReplaceAll(new_op.Name, "\"", "")
			} else {
				return new_op.Name
			}
		}
	}
	return ""
}

// generate the fucking nodes
func generate_asm(node Node) string {
	switch reflect.TypeOf(node).Name() {
	case "AssignemntNode":
		node_cast := node.(AssignemntNode)
		assignment_code := "\t"
		var string_flg bool = false
		var string_const_asgn string = ""
		var string_const_asgn_length_and_val = ""

		if node_cast.Asgn_type == "str" {
			out_cpy := *out_code
			cnt_cpy := current_function_variable_count
			string_flg = true
			//node_cast.Ptrs++
			assignment_code = ""
			fmt.Println(node_cast.Content.What_type())
			assignment_code += fmt.Sprintf("@.str.%s =%s [%d x i8] c\"%s\\00\", align 1\n", node_cast.Var_name, string_constant, len(generate_second_and_third_class_node(node_cast.Content, type_convert_llvm(node_cast.Asgn_type)))+1, generate_second_and_third_class_node(node_cast.Content, type_convert_llvm(node_cast.Asgn_type)))
			string_const_asgn_length_and_val = fmt.Sprintf("[%d x i8]", len(generate_second_and_third_class_node(node_cast.Content, type_convert_llvm(node_cast.Asgn_type)))+1)
			string_const_asgn = assignment_code
			node_cast.Ptrs++
			current_ptr_offset++
			*out_code = out_cpy
			current_function_variable_count = cnt_cpy
			assignment_code = "\t"
		}

		// sets the variable to global (@) or local (%)
		if node_cast.Global {
			assignment_code += "@"
		} else {
			assignment_code += "%"
		}

		// write the name of the variable and make it an empty pointer to it's type
		assignment_code += node_cast.Var_name + " = alloca "

		converted_type := type_convert_llvm(node_cast.Asgn_type)

		assignment_code += converted_type + get_align(node_cast.Ptrs)
		*out_code += assignment_code + "\n"
		// -> %var_name = alloca type, align 4 (pointer on the stack to a type)

		func_count_old := current_function_variable_count
		// store the actual value into %var_name
		_ = generate_second_and_third_class_node(node_cast.Content, converted_type)

		// assign var_name to a var. This works because the last temporary variable is always the final value to store
		if !string_flg {
			*out_code += fmt.Sprintf("\tstore %s %%%d, %s%s %%%s %s\n", converted_type, current_function_variable_count-1, converted_type, get_ptrs(node_cast.Ptrs+1), node_cast.Var_name, get_align(node_cast.Ptrs))
		} else {
			current_ptr_offset = 0
			if var_is_refferenced {
				*out_code += fmt.Sprintf("\tstore %s %%%d, %s%s %%%s %s\n", converted_type, current_function_variable_count-1, converted_type, get_ptrs(node_cast.Ptrs), node_cast.Var_name, get_align(node_cast.Ptrs))
			} else {
				*out_code += fmt.Sprintf("\tstore %s getelementptr inbounds (%s, %s%s %s, i64 0, i64 0),%s%s %%%s %s\n", converted_type, string_const_asgn_length_and_val, string_const_asgn_length_and_val, get_ptrs(node_cast.Ptrs), "@.str."+node_cast.Var_name, converted_type, get_ptrs(node_cast.Ptrs), node_cast.Var_name, get_align(node_cast.Ptrs))
				if func_count_old == current_function_variable_count {
					*global_vars += string_const_asgn
				}
			}
		}

		VARS_LLVM = append(VARS_LLVM, node_cast.Var_name)
		return ""
	case "":
		break
	default:
		break
	}
	return ""
}
