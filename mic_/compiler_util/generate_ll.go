package compiler_util

import (
	"fmt"
	"os"
	//"os/exec"
	"reflect"
	"strconv"
	"strings"
)

var current_function_variable_count int = 1 // counts the temporary variables specified within a function
var main_function_variable_count int = 1    // counts the temporary variables specified within the main function
var temporary_string_count int = 1          // counts the temporary string variables

var var_is_referenced bool = false

var VARS_LLVM []string

func get_align(ptrs int, i string) string {
	if i == "char" || i == "i8" {
		return ", align 1"
	}
	if ptrs > 0 {
		return ", align 8"
	} else {
		return ", align 4"
	}
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
	case "-":
		return "sub"
	case "*":
		return "mul"
	case "/":
		return "div"
	}
	return ""
}
func operator_prefix_cast_llvm(i, ii string) string {
	switch i {
	case "i32":
		if ii == "div" {
			return "s"
		}
		return ""
	case "i8*":
		if ii == "div" {
			return "s"
		}
		return ""
	case "i8":
		if ii == "div" {
			return "s"
		}
		return ""
	case "float":
		return "f"
	case "i64":
		if ii == "div" {
			return "s"
		}
		return ""
	}
	return ""
}

func type_convert_llvm(i string) string {
	switch i {
	case "int":
		return "i32"
	case "str":
		return "i8*"
	case "char":
		return "i8"
	case "flt":
		return "float"
	case "cock":
		return "i64"
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
	*out_code += "\ndefine i32 @main() {\n"
	for !is_eoast_asm {
		generate_asm(current_node_gasm)
		gasm_advance()
	}

	*out_code += "\tret i32 1\n}"
	*out_code = *global_vars + *out_code
	err := os.WriteFile(out_pth+".ll", []byte(*out_code), 0755)
	if err != nil {
		panic(err)
	}

	//_, _ = exec.Command("clang", out_pth+".ll", "-o", out_pth+".o").CombinedOutput()
	//_, _ = exec.Command("llc", "-filetype=obj", out_pth+".ll", "-o", out_pth).CombinedOutput()
	//_, _ = exec.Command("rm", "-rf", out_pth+".ll").CombinedOutput()

	return SUCCESS
}

func generate_second_and_third_class_node(node LiteralNode, type_ string) (string, bool) {
	switch reflect.TypeOf(node).Name() {
	case "BinOpNode":

		new_op := node.(BinOpNode)

		left, _ := generate_second_and_third_class_node(new_op.Left_node, type_)
		right, _ := generate_second_and_third_class_node(new_op.Right_node, type_)

		left_code := type_ + " " + left
		operator := operator_prefix_cast_llvm(type_, operator_cast_llvm(new_op.Op_tok)) + operator_cast_llvm(new_op.Op_tok)
		right_code := right

		code := fmt.Sprintf("\t%%%d = ", current_function_variable_count)
		code += operator + " "
		code += left_code + ", "
		code += right_code + "\n"

		*out_code += code
		current_function_variable_count++
		return fmt.Sprintf("\t%%%d", current_function_variable_count-1), false

	case "DirectNode":
		new_op := node.(DirectNode)
		return new_op.Value, true
	case "VarNameNode":
		var_is_referenced = false
		new_op := node.(VarNameNode)

		// return "type %var_name" if the current var node is a var elese return the literal
		if StringInSlice(new_op.Name, VARS_LLVM) {
			var_is_referenced = true
			code := fmt.Sprintf("\t%%%d = load %s, %s%s %%%s %s\n", current_function_variable_count, type_, type_, get_ptrs(new_op.Ptrs+1), new_op.Name, get_align(new_op.Ptrs+current_ptr_offset, type_))

			*out_code += code
			current_function_variable_count++

			return fmt.Sprintf("\t%%%d", (current_function_variable_count - 1)), true
		} else {
			if type_ == "i8*" {
				return strings.ReplaceAll(new_op.Name, "\"", ""), true
			} else if type_ == "i8" {
				value := int(rune(strings.ReplaceAll(new_op.Name, "'", "")[0]))
				return strconv.Itoa(value), true
			} else {
				return new_op.Name, true
			}
		}
	}
	return "", true
}

// generate the fucking nodes
func generate_asm(node Node) string {
	switch reflect.TypeOf(node).Name() {
	case "AssignmentNode":
		node_cast := node.(AssignmentNode)
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
			gen_cnt, _ := generate_second_and_third_class_node(node_cast.Content, type_convert_llvm(node_cast.Asgn_type))
			assignment_code += fmt.Sprintf("@.str.%s =%s [%d x i8] c\"%s\\00\", align 1\n", node_cast.Var_name, string_constant, len(gen_cnt)+1, gen_cnt)
			string_const_asgn_length_and_val = fmt.Sprintf("[%d x i8]", len(gen_cnt)+1)
			string_const_asgn = assignment_code
			//node_cast.Ptrs++
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

		if string_flg {
			assignment_code += converted_type + get_ptrs(node_cast.Ptrs) + ", align 8"
		} else {
			assignment_code += converted_type + get_ptrs(node_cast.Ptrs) + get_align(node_cast.Ptrs, node_cast.Asgn_type)
		}
		*out_code += assignment_code + "\n"
		// -> %var_name = alloca type, align 4 (pointer on the stack to a type)

		if reflect.TypeOf(node_cast.Content).Name() != "UniversalNone" {
			func_count_old := current_function_variable_count
			// store the actual value into %var_name
			code, is_req := generate_second_and_third_class_node(node_cast.Content, converted_type)

			// assign var_name to a var. This works because the last temporary variable is always the final value to store
			if !string_flg {
				if !is_req {
					*out_code += fmt.Sprintf("\tstore %s %%%d, %s%s %%%s %s\n", converted_type, current_function_variable_count-1, converted_type, get_ptrs(node_cast.Ptrs+1), node_cast.Var_name, get_align(node_cast.Ptrs, node_cast.Asgn_type))
				} else {
					*out_code += fmt.Sprintf("\tstore %s %s, %s%s %%%s %s\n", converted_type, code, converted_type, get_ptrs(node_cast.Ptrs+1), node_cast.Var_name, get_align(node_cast.Ptrs, node_cast.Asgn_type))
				}
			} else {
				current_ptr_offset = 0
				if var_is_referenced {
					*out_code += fmt.Sprintf("\tstore %s %%%d, %s%s %%%s %s\n", converted_type, current_function_variable_count-1, converted_type, get_ptrs(node_cast.Ptrs+1), node_cast.Var_name, get_align(node_cast.Ptrs+1, node_cast.Asgn_type))
				} else {
					*out_code += fmt.Sprintf("\tstore %s getelementptr inbounds (%s, %s%s %s, i64 0, i64 0),%s%s %%%s %s\n", converted_type, string_const_asgn_length_and_val, string_const_asgn_length_and_val, get_ptrs(node_cast.Ptrs+1), "@.str."+node_cast.Var_name, converted_type, get_ptrs(node_cast.Ptrs+1), node_cast.Var_name, ", align 8")
					if func_count_old == current_function_variable_count {
						*global_vars += string_const_asgn
					}
				}
			}
		}

		VARS_LLVM = append(VARS_LLVM, node_cast.Var_name)
		return ""
	case "ReAssignmentNode":
		node_cast := node.(ReAssignmentNode)
		
		switch node_cast.Reassgn_t {
		case "+=":
			// load self (%0 = load i32* self)
			// %1 = add i32 %0, expr_erg
			// store %1, i32* self
			converted_type := type_convert_llvm(VARS[node_cast.Re_type].(AssignmentNode).Asgn_type)

			*out_code += fmt.Sprintf("\t%%%d = load %s, %s%s %%%s %s\n", current_function_variable_count, converted_type, converted_type, get_ptrs(VARS[node_cast.Re_type].(AssignmentNode).Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(VARS[node_cast.Re_type].(AssignmentNode).Ptrs+current_ptr_offset, converted_type))
			var_pos := current_function_variable_count
			current_function_variable_count++

			code, is_req := generate_second_and_third_class_node(node_cast.Content, node_cast.Re_type)
			if is_req {
				*out_code += fmt.Sprintf("\t%%%d = %sadd %s %s, %%%d\n", current_function_variable_count, operator_prefix_cast_llvm(converted_type, "add"), converted_type, code, var_pos)
				current_function_variable_count++
				//*out_code += fmt.Sprintf("\tstore %s %%%d, %s%s %%%s %s\n", converted_type, current_function_variable_count-1, converted_type, get_ptrs(node_cast.Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(node_cast.Ptrs, VARS[node_cast.Re_type].(AssignmentNode).Asgn_type))
			} else {
				*out_code += fmt.Sprintf("\t%%%d = %sadd %s %%%d, %%%d\n", current_function_variable_count, operator_prefix_cast_llvm(converted_type, "add"), converted_type, current_function_variable_count-1, var_pos)
				current_function_variable_count++
				//*out_code += fmt.Sprintf("\tstore %s %s, %s%s %%%s %s\n", converted_type, code, converted_type, get_ptrs(node_cast.Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(node_cast.Ptrs, VARS[node_cast.Re_type].(AssignmentNode).Asgn_type))
			}
			*out_code += fmt.Sprintf("\tstore %s %%%d, %s%s %%%s %s\n", converted_type, current_function_variable_count-1, converted_type, get_ptrs(node_cast.Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(node_cast.Ptrs, VARS[node_cast.Re_type].(AssignmentNode).Asgn_type))
		case "-=":
			// load self (%0 = load i32* self)
			// %1 = sub i32 %0, expr_erg
			// store %1, i32* self
			converted_type := type_convert_llvm(VARS[node_cast.Re_type].(AssignmentNode).Asgn_type)

			*out_code += fmt.Sprintf("\t%%%d = load %s, %s%s %%%s %s\n", current_function_variable_count, converted_type, converted_type, get_ptrs(VARS[node_cast.Re_type].(AssignmentNode).Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(VARS[node_cast.Re_type].(AssignmentNode).Ptrs+current_ptr_offset, converted_type))
			var_pos := current_function_variable_count
			current_function_variable_count++

			code, is_req := generate_second_and_third_class_node(node_cast.Content, node_cast.Re_type)
			if is_req {
				*out_code += fmt.Sprintf("\t%%%d = %ssub %s %s, %%%d\n", current_function_variable_count, operator_prefix_cast_llvm(converted_type, "sub"), converted_type, code, var_pos)
				current_function_variable_count++
				//*out_code += fmt.Sprintf("\tstore %s %%%d, %s%s %%%s %s\n", converted_type, current_function_variable_count-1, converted_type, get_ptrs(node_cast.Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(node_cast.Ptrs, VARS[node_cast.Re_type].(AssignmentNode).Asgn_type))
			} else {
				*out_code += fmt.Sprintf("\t%%%d = %ssub %s %%%d, %%%d\n", current_function_variable_count, operator_prefix_cast_llvm(converted_type, "sub"), converted_type, current_function_variable_count-1, var_pos)
				current_function_variable_count++
				//*out_code += fmt.Sprintf("\tstore %s %s, %s%s %%%s %s\n", converted_type, code, converted_type, get_ptrs(node_cast.Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(node_cast.Ptrs, VARS[node_cast.Re_type].(AssignmentNode).Asgn_type))
			}
			*out_code += fmt.Sprintf("\tstore %s %%%d, %s%s %%%s %s\n", converted_type, current_function_variable_count-1, converted_type, get_ptrs(node_cast.Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(node_cast.Ptrs, VARS[node_cast.Re_type].(AssignmentNode).Asgn_type))
		case "*=":
			converted_type := type_convert_llvm(VARS[node_cast.Re_type].(AssignmentNode).Asgn_type)

			*out_code += fmt.Sprintf("\t%%%d = load %s, %s%s %%%s %s\n", current_function_variable_count, converted_type, converted_type, get_ptrs(VARS[node_cast.Re_type].(AssignmentNode).Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(VARS[node_cast.Re_type].(AssignmentNode).Ptrs+current_ptr_offset, converted_type))
			var_pos := current_function_variable_count
			current_function_variable_count++

			code, is_req := generate_second_and_third_class_node(node_cast.Content, node_cast.Re_type)
			if is_req {
				*out_code += fmt.Sprintf("\t%%%d = %smul %s %s, %%%d\n", current_function_variable_count, operator_prefix_cast_llvm(converted_type, "mul"), converted_type, code, var_pos)
				current_function_variable_count++
				//*out_code += fmt.Sprintf("\tstore %s %%%d, %s%s %%%s %s\n", converted_type, current_function_variable_count-1, converted_type, get_ptrs(node_cast.Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(node_cast.Ptrs, VARS[node_cast.Re_type].(AssignmentNode).Asgn_type))
			} else {
				*out_code += fmt.Sprintf("\t%%%d = %smul %s %%%d, %%%d\n", current_function_variable_count, operator_prefix_cast_llvm(converted_type, "mul"), converted_type, current_function_variable_count-1, var_pos)
				current_function_variable_count++
				//*out_code += fmt.Sprintf("\tstore %s %s, %s%s %%%s %s\n", converted_type, code, converted_type, get_ptrs(node_cast.Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(node_cast.Ptrs, VARS[node_cast.Re_type].(AssignmentNode).Asgn_type))
			}
			*out_code += fmt.Sprintf("\tstore %s %%%d, %s%s %%%s %s\n", converted_type, current_function_variable_count-1, converted_type, get_ptrs(node_cast.Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(node_cast.Ptrs, VARS[node_cast.Re_type].(AssignmentNode).Asgn_type))
		case "/*":
			converted_type := type_convert_llvm(VARS[node_cast.Re_type].(AssignmentNode).Asgn_type)

			*out_code += fmt.Sprintf("\t%%%d = load %s, %s%s %%%s %s\n", current_function_variable_count, converted_type, converted_type, get_ptrs(VARS[node_cast.Re_type].(AssignmentNode).Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(VARS[node_cast.Re_type].(AssignmentNode).Ptrs+current_ptr_offset, converted_type))
			var_pos := current_function_variable_count
			current_function_variable_count++

			code, is_req := generate_second_and_third_class_node(node_cast.Content, node_cast.Re_type)
			if is_req {
				*out_code += fmt.Sprintf("\t%%%d = %sdiv %s %s, %%%d\n", current_function_variable_count, operator_prefix_cast_llvm(converted_type, "div"), converted_type, code, var_pos)
				current_function_variable_count++
				//*out_code += fmt.Sprintf("\tstore %s %%%d, %s%s %%%s %s\n", converted_type, current_function_variable_count-1, converted_type, get_ptrs(node_cast.Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(node_cast.Ptrs, VARS[node_cast.Re_type].(AssignmentNode).Asgn_type))
			} else {
				*out_code += fmt.Sprintf("\t%%%d = %sdiv %s %%%d, %%%d\n", current_function_variable_count, operator_prefix_cast_llvm(converted_type, "div"), converted_type, current_function_variable_count-1, var_pos)
				current_function_variable_count++
				//*out_code += fmt.Sprintf("\tstore %s %s, %s%s %%%s %s\n", converted_type, code, converted_type, get_ptrs(node_cast.Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(node_cast.Ptrs, VARS[node_cast.Re_type].(AssignmentNode).Asgn_type))
			}
			*out_code += fmt.Sprintf("\tstore %s %%%d, %s%s %%%s %s\n", converted_type, current_function_variable_count-1, converted_type, get_ptrs(node_cast.Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(node_cast.Ptrs, VARS[node_cast.Re_type].(AssignmentNode).Asgn_type))
		case "?=":
			converted_type := type_convert_llvm(VARS[node_cast.Re_type].(AssignmentNode).Asgn_type)

			if VARS[node_cast.Re_type].(AssignmentNode).Asgn_type == "str" {
				func_count_old := current_function_variable_count

				out_cpy := *out_code
				cnt_cpy := current_function_variable_count
				//node_cast.Ptrs++
				assignment_code := ""
				gen_cnt, _ := generate_second_and_third_class_node(node_cast.Content, type_convert_llvm(VARS[node_cast.Re_type].(AssignmentNode).Asgn_type))
				assignment_code += fmt.Sprintf("@.str.%s =%s [%d x i8] c\"%s\\00\", align 1\n", VARS[node_cast.Re_type].(AssignmentNode).Var_name, string_constant, len(gen_cnt)+1, gen_cnt)
				string_const_asgn_length_and_val := fmt.Sprintf("[%d x i8]", len(gen_cnt)+1)
				string_const_asgn := assignment_code
				//node_cast.Ptrs++
				current_ptr_offset++
				*out_code = out_cpy
				current_function_variable_count = cnt_cpy
				assignment_code = "\t"

				current_ptr_offset = 0
				if var_is_referenced {
					*out_code += fmt.Sprintf("\tstore %s %%%d, %s%s %%%s %s\n", converted_type, current_function_variable_count-1, converted_type, get_ptrs(node_cast.Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(node_cast.Ptrs+1, VARS[node_cast.Re_type].(AssignmentNode).Asgn_type))
				} else {
					*out_code += fmt.Sprintf("\tstore %s getelementptr inbounds (%s, %s%s %s, i64 0, i64 0),%s%s %%%s %s\n", converted_type, string_const_asgn_length_and_val, string_const_asgn_length_and_val, get_ptrs(node_cast.Ptrs+1), "@.str."+VARS[node_cast.Re_type].(AssignmentNode).Var_name, converted_type, get_ptrs(node_cast.Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, ", align 8")
					if func_count_old == current_function_variable_count {
						*global_vars += string_const_asgn
					}
				}
			} else {
				code, is_req := generate_second_and_third_class_node(node_cast.Content, converted_type)
				if !is_req {
					*out_code += fmt.Sprintf("\tstore %s %%%d, %s%s %%%s %s\n", converted_type, current_function_variable_count-1, converted_type, get_ptrs(node_cast.Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(node_cast.Ptrs, VARS[node_cast.Re_type].(AssignmentNode).Asgn_type))
				} else {
					*out_code += fmt.Sprintf("\tstore %s %s, %s%s %%%s %s\n", converted_type, code, converted_type, get_ptrs(node_cast.Ptrs+1), VARS[node_cast.Re_type].(AssignmentNode).Var_name, get_align(node_cast.Ptrs, VARS[node_cast.Re_type].(AssignmentNode).Asgn_type))
				}
			}
		}
	case "IfNode":
		
	default:
		break
	}
	return ""
}
