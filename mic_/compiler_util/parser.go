package compiler_util

import (
	"fmt"
	"reflect"
	"strings"
)

var tokens []Token
var func_names []string
var var_names []string

var pos int = -1
var current_token Token
var is_eot = false
var root_node RootNode = RootNode{}
var func_on bool = false
var VARS map[string]Node = make(map[string]Node)        // holds assigned vars
var FUNCTIONS map[string]Node = make(map[string]Node)   // holds all functions that have a body
var STRUCTS map[string][]Node = make(map[string][]Node) // holds all structs and their vars
var GLOBALS map[string]Node = make(map[string]Node)
var BOOL bool = false
var PREV_IF bool = false
var LOOP bool = false
var global bool = false
var CURRENT_LN int

var expected_ptrs int = 0
var expect_t bool = false
var expect_ret = false
var func_ret_expect string
var current_expected_t string

var current_ln_str string = ""
var expect_doc bool = false

// return error loc string
func CreateErrorLocation() string {
	return fmt.Sprintf("%s at ln %d -> %s", current_token.section, current_token.ln_count-1, current_ln_str)
}

// parse that shit wohooooooo
func Parse(Tokens *[]Token, Func_names, Var_names, ign_sections []string) RootNode {
	tokens = *Tokens
	func_names = Func_names
	var_names = Var_names
	p_advance()

	fmt.Println(current_token)
	for !is_eot {
		if StringInSlice(current_token.section, ign_sections) {
			if current_token.type_ == TT_ID || current_token.type_ == TT_MUL || current_token.type_ == TT_DOC {
				_ = mkID()
				//fmt.Println(res)
				//fmt.Println(VARS)
			} else if current_token.type_ == TT_DEBUG {
				root_node = root_node.AddNodeToRoot(DebugNode{})
			} else {
				NewError("ParsingError", "A function decleration, struct decleration, variable assignement or refference was expected but not found", CreateErrorLocation(), true)
			}
			continue
		}
		if current_token.type_ == TT_ID || current_token.type_ == TT_MUL || current_token.type_ == TT_DOC {
			res := mkID()
			//fmt.Println(res)
			//fmt.Println(VARS)
			root_node = root_node.AddNodeToRoot(res)
		} else if current_token.type_ == TT_DEBUG {
			root_node = root_node.AddNodeToRoot(DebugNode{})
		} else {
			NewError("ParsingError", "A function decleration, struct decleration, variable assignement or refference was expected but not found", CreateErrorLocation(), true)
		}
	}
	return root_node
}

/*-------------------Get and increment methods-------------------*/

// compares types to find errors
func type_compare(comp string) {
	if strings.ToLower(current_expected_t) != strings.ToLower(comp) && expect_t {
		NewError("TypeMissmatchError", fmt.Sprintf("%s was expected but got %s", strings.ToLower(current_expected_t), strings.ToLower(comp)), CreateErrorLocation(), true)
	}
}

// increments the pos and sets the token if the position still is in the index of tokens
func p_advance() {
	pos++
	if pos < len(tokens) {
		current_token = tokens[pos]
		if current_token.ln_count != CURRENT_LN {
			// generate current ln string
			current_ln_str = current_token.value
			i := 1
			for {
				current_token_in_string := getToken(i)
				if current_token_in_string.ln_count != current_token.ln_count || is_eot {
					break
				} else {
					current_ln_str += current_token_in_string.value
				}
				i++
			}
		}
		CURRENT_LN = current_token.ln_count
	} else {
		tok := current_token
		current_token = Token{section: tok.section, ln_count: tok.ln_count + 1, type_: TT_EOF, value: ""}
		is_eot = true
	}
}

// gets a token in the future by index for example -> getToken(1) would give you the next token
func getToken(num int) Token {
	pos_now := pos
	if pos_now+num < len(tokens) {
		return tokens[pos_now+num]
	} else {
		//is_eot = true
		return Token{}
	}
}

// tries to make a name for example -> hello.id is legal
func getName() Token {
	name_node := current_token
	if (getToken(1).type_ == TT_DOT || getToken(1).type_ == TT_ARROW) && getToken(2).type_ == TT_ID && (!StringInSlice(getToken(2).value, TYPES) && !StringInSlice(getToken(2).value, INSTRUCTIONS) && !StringInSlice(getToken(2).value, CUSTOM_TYPES)) {
		// iterates through the tokens for as long as it's a continous name with dots like -> hello.type
		for (getToken(1).type_ == TT_DOT || getToken(1).type_ == TT_ARROW) && getToken(2).type_ == TT_ID && (!StringInSlice(getToken(2).value, TYPES) && !StringInSlice(getToken(2).value, INSTRUCTIONS) && !StringInSlice(getToken(2).value, CUSTOM_TYPES)) {
			p_advance()
			name_node.value += current_token.value
			p_advance()
			name_node.value += current_token.value
		}
	}
	return name_node
}

// gets as many pointers as there are behind a type
func getPointers() int {
	ptrs := 0
	for current_token.type_ == TT_MUL {
		ptrs++
		p_advance()
	}
	return ptrs
}

// gets stuff you would parse in a func definiton
func p_func_parse(end_type string) []Node {
	var vars []Node
	for current_token.type_ != end_type {
		var node Node

		tok := current_token
		next_tok := getToken(1)
		next_next_tok := getToken(2)

		// check if it's a comma to much
		if next_tok.type_ == TT_COMMA && (StringInSlice(next_next_tok.value, TYPES) || StringInSlice(next_next_tok.value, CUSTOM_TYPES)) {
			NewError("ArgumentExpectedError", "After a comma in struct a value is expected but there was none!", CreateErrorLocation(), true)
		}
		if StringInSlice(tok.value, TYPES) || StringInSlice(tok.value, CUSTOM_TYPES) {
			node = p_assign(tok.value, false)
		} else if tok.type_ == TT_COMMA {
			p_advance()
			continue
		} else {
			NewError("IvalidTypeError", "A type that doesn't exist was specified"+tok.type_, CreateErrorLocation(), true)
		}
		vars = append(vars, node)
	}
	p_advance()
	return vars
}

// gets the inputed tokens at length
func p_func_call_parse(end_type string, expected_len int, len_expected bool) []LiteralNode {
	var vars []LiteralNode
	var iters int = 0

	for current_token.type_ != end_type {
		if len_expected {
			if iters > expected_len {
				break
			}
		}
		var node LiteralNode

		tok := current_token
		next_tok := getToken(1)
		next_next_tok := getToken(2)

		// check if it's a comma to much
		if next_tok.type_ == TT_COMMA && (StringInSlice(next_next_tok.value, TYPES) || StringInSlice(next_next_tok.value, CUSTOM_TYPES)) {
			NewError("ArgumentExpectedError", "After a comma in struct a value is expected but there was none!", CreateErrorLocation(), true)
		}
		if StringInMap(tok.value, VARS) || StringInMap(tok.value, FUNCTIONS) {
			node = p_expr()
			vars = append(vars, node)
		} else if tok.type_ == TT_COMMA {
			p_advance()
			continue
		} else {
			node = p_expr()
			vars = append(vars, node)
		}
		iters++
	}
	return vars
}

// makes struct args of of a list in STRUCTS
func p_make_struct_args(name, type_ string, is_ptr, is_n_ptr bool) {
	str_get := STRUCTS[type_]
	for i := 0; len(str_get) > i; i++ {
		current_struct_attr := str_get[i]
		if current_struct_attr.What_type() == "AssignementNode" {
			current_struct_attr2 := current_struct_attr.(AssignemntNode)
			if StringInSlice(current_struct_attr2.Asgn_type, CUSTOM_TYPES) {
				if is_ptr {
					var is_ptr_ bool = false
					var is_n_ptr_ bool = false
					if current_struct_attr2.Ptrs > 0 {
						is_ptr_ = true
						is_n_ptr_ = false
					} else {
						is_ptr_ = false
						is_n_ptr_ = true
					}
					VARS[name+"."+current_struct_attr2.Var_name] = AssignemntNode{current_struct_attr2.Asgn_type, current_struct_attr2.Ptrs, false, name + "." + current_struct_attr2.Var_name, current_struct_attr2.Content}
					p_make_struct_args(name+"->"+current_struct_attr2.Var_name, current_struct_attr2.Asgn_type, is_ptr_, is_n_ptr_)
					//p_make_struct_args(name+"->"+current_struct_attr2.Var_name, current_struct_attr2.Asgn_type, true, false)
				} else {
					var is_ptr_ bool = false
					var is_n_ptr_ bool = false
					if current_struct_attr2.Ptrs > 0 {
						is_ptr_ = true
						is_n_ptr_ = false
					} else {
						is_ptr_ = false
						is_n_ptr_ = true
					}
					VARS[name+"->"+current_struct_attr2.Var_name] = AssignemntNode{current_struct_attr2.Asgn_type, current_struct_attr2.Ptrs, false, name + "->" + current_struct_attr2.Var_name, current_struct_attr2.Content}
					p_make_struct_args(name+"."+current_struct_attr2.Var_name, current_struct_attr2.Asgn_type, is_ptr_, is_n_ptr_)
				}
			} else {
				if is_n_ptr {
					VARS[name+"."+current_struct_attr2.Var_name] = AssignemntNode{current_struct_attr2.Asgn_type, current_struct_attr2.Ptrs, false, name + "." + current_struct_attr2.Var_name, current_struct_attr2.Content}
				}
				if is_ptr {
					VARS[name+"->"+current_struct_attr2.Var_name] = AssignemntNode{current_struct_attr2.Asgn_type, current_struct_attr2.Ptrs, false, name + "->" + current_struct_attr2.Var_name, current_struct_attr2.Content}
				}
			}
		} else if current_struct_attr.What_type() == "ArrAssignementNode" {
			current_struct_attr2 := current_struct_attr.(ArrAssignementNode)
			if StringInSlice(current_struct_attr2.Asgn_type, CUSTOM_TYPES) {
				if is_n_ptr {
					p_make_struct_args(name+"."+current_struct_attr2.Array_name, current_struct_attr2.Asgn_type, false, true)
					VARS[name+"."+current_struct_attr2.Array_name] = ArrAssignementNode{current_struct_attr2.Asgn_type, current_struct_attr2.Ptrs, false, name + "." + current_struct_attr2.Array_name, current_struct_attr2.Arr_len}
				}
				if is_ptr {
					p_make_struct_args(name+"->"+current_struct_attr2.Array_name, current_struct_attr2.Asgn_type, true, false)
					VARS[name+"->"+current_struct_attr2.Array_name] = AssignemntNode{current_struct_attr2.Asgn_type, current_struct_attr2.Ptrs, false, name + "->" + current_struct_attr2.Array_name, current_struct_attr2.Arr_len}
				}
			} else {
				if is_n_ptr {
					VARS[name+"."+current_struct_attr2.Array_name] = ArrAssignementNode{current_struct_attr2.Asgn_type, current_struct_attr2.Ptrs, false, name + "." + current_struct_attr2.Array_name, current_struct_attr2.Arr_len}
				}
				if is_ptr {
					VARS[name+"->"+current_struct_attr2.Array_name] = ArrAssignementNode{current_struct_attr2.Asgn_type, current_struct_attr2.Ptrs, false, name + "->" + current_struct_attr2.Array_name, current_struct_attr2.Arr_len}
				}
			}
		} else {
			NewError("WTF", "", CreateErrorLocation(), true)
		}
	}
}

/*-------------------expr, factor and binary op-------------------*/

/* returns a lot of things:
-	function_call: example -> func_name(10, 20)
-	variable_name: example -> (*...|&)var
-	list_slice: example (*..|&)list[5]
-	int: example -> 10
-	flt: example -> 10.5
-	str: example -> "Hello"
-	char: example -> 'a'
*/
func p_factor() LiteralNode {
	tok := current_token
	if current_token.ln_count != CURRENT_LN {
		return UniversalNone{}
	}
	var ptrs int = 0
	var deref bool = false
	var not bool = false
	var minus bool = false
	var bit_not = false
	// returns either a func_call, var_name or array_slcie
	if tok.type_ == TT_MINUS {
		minus = true
		p_advance()
	} else if tok.type_ == TT_BITNOT {
		bit_not = true
		p_advance()
	}

	if tok.type_ == TT_MUL {
		ptrs = getPointers()
	} else if tok.type_ == TT_KAND {
		deref = true
		p_advance()
	} else if tok.type_ == TT_NOT && BOOL {
		not = true
		p_advance()
	}
	tok = current_token
	if tok.type_ == TT_ID {
		tok = getName()
		p_advance()
		if StringInMap(tok.value, VARS) || StringInSlice(tok.value, var_names) {
			if reflect.TypeOf(VARS[tok.value]).Name() == "ArrAssignementNode" {
				comp_t := VARS[tok.value].(ArrAssignementNode)
				type_compare(comp_t.Asgn_type)

				if expected_ptrs > 0 && current_token.type_ != TT_LBRK && getToken(1).ln_count != tok.ln_count {
					return VarNameNode{comp_t.Array_name, ptrs, deref, not, minus, bit_not}
				}
			}
			// return list slice if [] is found
			if current_token.type_ == TT_LBRK {
				current_expected_t = TT_INT
				expected_ptrs = 0
				expect_t = true
				p_advance()
				arr_pos := p_expr()

				if current_token.type_ == TT_RBRK {
					p_advance()
					return ListSliceNode{Name: tok.value, Pos: arr_pos, Ptrs: ptrs, Deref: deref, Not: not, Minus: minus, BitNot: bit_not}
				} else {
					NewError("ArrayNotClosed", "A array was assagnid with '[' but no ']' was found. ", CreateErrorLocation(), true)
				}
			}
			comp_t := VARS[tok.value].(AssignemntNode)
			type_compare(comp_t.Asgn_type)
			if ptrs > comp_t.Ptrs {
				// is illegal 'cause it's trying to use more ptrs than there is
				NewError("ToManyPointersException", "", CreateErrorLocation(), true)
			}
			if expected_ptrs == 0 && ptrs != comp_t.Ptrs { // checks if is prefixed by right amount of ptrs
				NewError("PointerMissmatchException", "", CreateErrorLocation(), true)
			}
			if ptrs <= comp_t.Ptrs || (expected_ptrs > 0 && deref) { // checks if there is the right amount of ptrs or if it's derefed
				return VarNameNode{tok.value, ptrs, deref, not, minus, bit_not}
			} else {
				NewError("PointerMissmatchException", tok.value, CreateErrorLocation(), true)
			}
		} else if StringInMap(tok.value, FUNCTIONS) || StringInSlice(tok.value, func_names) {
			call_name := tok.value
			if StringInMap(tok.value, FUNCTIONS) {
				switch reflect.TypeOf(FUNCTIONS[call_name]).Name() {
				case "FunctionNode":
					type_compare(FUNCTIONS[call_name].(FunctionNode).Ret_type)
				case "AsmFunctionNode":
					type_compare(FUNCTIONS[call_name].(AsmFunctionNode).Ret_type)
				default:
					break
				}
			}
			expect_t = false
			if current_token.type_ == TT_LPAREN {
				p_advance()
				var call_args []LiteralNode
				if StringInMap(tok.value, FUNCTIONS) {
					call_args = p_func_call_parse(TT_RPAREN, Arg_len(FUNCTIONS[call_name]), true)
				} else {
					call_args = p_func_call_parse(TT_RPAREN, 0, false)
				}
				p_advance()
				return FuncCallNode{Call_name: call_name, Func_parse: call_args, Minus: minus, BitNot: bit_not}
			}
		} else {
			fmt.Println(tok)
			NewError("RefferenceError", "The refferenced var, func or struct is not defined", CreateErrorLocation(), true)
		}
	}
	p_advance()
	type_compare(tok.type_)
	if expected_ptrs == ptrs {
		return VarNameNode{tok.value, ptrs, deref, not, minus, bit_not}
	} else {
		NewError("PointerMissmatchException", current_token.value, CreateErrorLocation(), true)
	}
	return DirectNode{Type_: tok.type_, Value: tok.value, Minus: minus, BitNot: bit_not}
}
func p_term() LiteralNode {
	ops := []string{TT_MUL, TT_DIV}
	return binOp(true, false, ops)
}
func p_expr() LiteralNode {
	ops := []string{TT_PLUS, TT_MINUS, TT_SHIFTL, TT_SHIFTR, TT_BITAND, TT_BITOR, TT_BITXOR}
	return binOp(false, true, ops)
}

func p_bool_expr() LiteralNode {
	var bool_ops []string = []string{
		TT_LTHEN,
		TT_LEQ,
		TT_GTHEN,
		TT_GEQ,
		TT_KAND,
		TT_AND,
		TT_OR,
		TT_NOT,
		TT_NEQ,
	}
	var current_ln int = current_token.ln_count
	var left LiteralNode
	BOOL = true
	left = p_expr()
	for StringInSlice(current_token.type_, bool_ops) && !is_eot && CURRENT_LN == current_ln {
		op_tok := current_token
		p_advance()

		var right LiteralNode
		expected_ptrs = 0
		right = p_expr()
		left = BoolOpNode{left, op_tok.value, right}
	}
	BOOL = false
	return left
}

// makes a binOp node and returns it
func binOp(factor, term bool, ops []string) LiteralNode {
	var left LiteralNode
	var current_ln int = current_token.ln_count
	if term {
		left = p_term()
	} else if factor {
		left = p_factor()
	}

	for StringInSlice(current_token.type_, ops) && !is_eot && CURRENT_LN == current_ln {
		op_tok := current_token
		p_advance()

		var right LiteralNode
		if term {
			right = p_term()
		} else if factor {
			right = p_factor()
		}
		left = BinOpNode{left, op_tok.value, right}
	}
	return left
}

/*-----------Assign for vars functions and also reassgn-----------*/

// makes an assignement with -> type name([expr()]) = expr()
func p_assign(type_ string, glob bool) Node {
	if type_ == "cock" {
		current_expected_t = "int"
	} else {
		current_expected_t = type_
	}

	expect_t = true
	p_advance()
	CURRENT_LN = current_token.ln_count
	// gets all pointers behind the type
	ptrs := getPointers()

	assignement_name := current_token

	if StringInSlice(type_, CUSTOM_TYPES) {
		if ptrs != 0 {
			p_make_struct_args(assignement_name.value, type_, true, false)
		} else {
			p_make_struct_args(assignement_name.value, type_, false, true)
		}
	}

	// checks if the variable allready exists
	if StringInMap(assignement_name.value, VARS) || StringInSlice(assignement_name.value, var_names) {
		NewError("VariableAllreadyAssigned", fmt.Sprintf("The Variable \"%s\" is allready assigned.", assignement_name.value), CreateErrorLocation(), true)
	}

	p_advance()
	// returns a array or a assignement based on if there are brackets or not
	if current_token.type_ == TT_LBRK {
		p_advance()
		var arr_len LiteralNode
		if current_token.type_ != TT_RBRK {
			current_expected_t = TT_INT
			expected_ptrs = 0
			arr_len = p_expr()
		} else {
			// array needs to be initialized with zeros in the c file -> type name[] = 0
			arr_len = UniversalNone{}
		}

		// checks if the arr expr is closed
		if current_token.type_ == TT_RBRK {
			p_advance()
			var_node := ArrAssignementNode{Asgn_type: type_, Ptrs: ptrs, Array_name: assignement_name.value, Arr_len: arr_len, Global: global}
			VARS[assignement_name.value] = var_node
			if glob {
				GLOBALS[assignement_name.value] = var_node
				global = false
			}
			expect_t = false
			return var_node
		} else {
			NewError("ArrayNotClosed", "A array was assagnid with '[' but no ']' was found. ", CreateErrorLocation(), true)
		}
	} else if current_token.type_ == TT_ASSGN {
		expected_ptrs = ptrs
		p_advance()
		// make a typecast node if one is found
		if current_token.type_ == TT_ID && current_token.value == "tcst" {
			expect_t = false
			p_advance()
			if current_token.type_ == TT_LPAREN {
				p_advance()
				var_name := p_expr()
				if current_token.type_ == TT_RPAREN {
					p_advance()
					ret_var := AssignemntNode{Asgn_type: type_, Ptrs: ptrs, Var_name: assignement_name.value, Content: TypeCastNode{Tcast: var_name, Dtype: DataTypeNode{Dtype: type_, Ptrs: ptrs}}}
					VARS[assignement_name.value] = ret_var
					if glob {
						GLOBALS[assignement_name.value] = ret_var
						global = false
					}
					return ret_var
				} else {
					NewError("ParantheseExpectes", "A parenthese in a typecast was expected but nor found", CreateErrorLocation(), true)
				}
			} else {
				NewError("ParantheseExpectes", "A parenthese in a typecast was expected but nor found", CreateErrorLocation(), true)
			}
		} else {
			content := p_expr()
			var_node := AssignemntNode{Asgn_type: type_, Ptrs: ptrs, Var_name: assignement_name.value, Content: content, Global: global}
			VARS[assignement_name.value] = var_node
			if glob {
				GLOBALS[assignement_name.value] = var_node
				global = false
			}
			expect_t = false
			return var_node
		}
	} else {
		var_node := AssignemntNode{Asgn_type: type_, Ptrs: ptrs, Var_name: assignement_name.value, Content: UniversalNone{}, Global: global}
		VARS[assignement_name.value] = var_node
		if glob {
			GLOBALS[assignement_name.value] = var_node
			global = false
		}
		expect_t = false
		return var_node
	}
	return AssignemntNode{}
}

// reassigns a variable or calls a function with no return type
func p_reassign(ptr bool) Node {
	// get the name as well as what the existing var will be assigned to. Then get the expected type and typecast it. After that update the variable to it's current status
	// it's only problematic for arrays. Maybe just don't update but return a reassignement node
	var ptrs int = 0
	//var glob bool
	var content LiteralNode
	if ptr {
		ptrs = getPointers()
	}
	CURRENT_LN = current_token.ln_count
	name := getName()

	if StringInMap(name.value, VARS) && reflect.TypeOf(VARS[name.value]).Name() == "AssignemntNode" {
		ptr_comp := VARS[name.value].(AssignemntNode)
		if ptrs == ptr_comp.Ptrs { // checks if the variable is prefixed with all the ptrs it got
			expected_ptrs = 0 // is possible 'cause basically it's like a normal variable
		} else if ptrs > VARS[name.value].(AssignemntNode).Ptrs {
			// is illegal 'cause it's trying to use more ptrs than there is
			NewError("ToManyPointersException", "", CreateErrorLocation(), true)
		} else {
			expected_ptrs = ptr_comp.Ptrs - ptrs
		}
	} else if StringInMap(name.value, VARS) && reflect.TypeOf(VARS[name.value]).Name() == "ArrAssignementNode" {
		ptr_comp := VARS[name.value].(ArrAssignementNode)
		if ptrs == ptr_comp.Ptrs { // checks if the variable is prefixed with all the ptrs it got
			expected_ptrs = 0 // is possible 'cause basically it's like a normal variable
		} else if ptrs > VARS[name.value].(AssignemntNode).Ptrs {
			// is illegal 'cause it's trying to use more ptrs than there is
			NewError("ToManyPointersException", "", CreateErrorLocation(), true)
		} else {
			expected_ptrs = ptr_comp.Ptrs - ptrs
		}
	}

	p_advance()
	/*if StringInMap(name.value, GLOBALS) {
		glob = true
	}*/
	if current_token.type_ == TT_LBRK {
		p_advance()
		current_expected_t = TT_INT
		prev_expected := expected_ptrs
		expected_ptrs = 0
		expect_t = true
		arr_idx := p_expr()
		expect_t = false
		if current_token.type_ == TT_RBRK {
			p_advance()
			if current_token.type_ == TT_REASSGN || current_token.type_ == TT_PLUSEQ || current_token.type_ == TT_MINUSEQ || current_token.type_ == TT_MULEQ || current_token.type_ == TT_DIVEQ {
				reassgn_t := current_token.value
				p_advance()
				if StringInMap(name.value, VARS) {
					current_expected_t = VARS[name.value].(ArrAssignementNode).Asgn_type
					expect_t = true
				} else {
					expect_t = false
				}
				expected_ptrs = prev_expected
				content = p_expr()
				expect_t = false
				return ArrReAssignementNode{Reassgn_t: reassgn_t, Re_type: name.value, Ptrs: ptrs, Arr_idx: arr_idx, Content: content}
			} else {
				// make error
				NewError("", "", "", true)
			}
		} else {
			// make error
			NewError("", "", "", true)
		}
	} else if current_token.type_ == TT_REASSGN || current_token.type_ == TT_PLUSEQ || current_token.type_ == TT_MINUSEQ || current_token.type_ == TT_MULEQ || current_token.type_ == TT_DIVEQ {
		reassgn_t := current_token.value
		p_advance()
		current_expected_t = VARS[name.value].(AssignemntNode).Asgn_type
		expect_t = true
		content := p_expr()
		expect_t = false
		ret_var := ReAssignmentNode{Reassgn_t: reassgn_t, Re_type: name.value, Ptrs: ptrs, Content: content}
		return ret_var
	} else if current_token.type_ == TT_LPAREN {
		var call_args []LiteralNode = []LiteralNode{UniversalNone{}}
		p_advance()
		if StringInMap(name.value, FUNCTIONS) {
			call_args = p_func_call_parse(TT_RPAREN, Arg_len(FUNCTIONS[name.value]), true)
		} else if current_token.type_ == TT_RPAREN {
			p_advance()
			return FuncCallNode{Call_name: name.value, Func_parse: call_args}
		} else {
			call_args = p_func_call_parse(TT_RPAREN, 0, false)
		}
		if current_token.type_ == TT_RPAREN {
			p_advance()
			return FuncCallNode{Call_name: name.value, Func_parse: call_args}
		} else {
			// make error
			NewError("ParentheseExpectedError", "A ')' was expected but not found", CreateErrorLocation(), true)
		}
	} else {
		// make error
		NewError("UnknownOperatorError", "Expected a reassignement, but no ?= or any other operator was found", CreateErrorLocation(), true)
	}
	return ReAssignmentNode{}
}

// makes a struct
func p_struct(estruct bool) Node {
	var typedef bool
	p_advance()
	struct_name := current_token
	p_advance()

	if current_token.value == "typedef" {
		typedef = true
		p_advance()
	} else if (current_token.type_ == TT_ID || current_token.type_ == TT_MUL) && !StringInMap(current_token.value, VARS) && getToken(1).type_ != TT_LCURL && StringInMapArray(struct_name.value, STRUCTS) && !StringInSlice(struct_name.value, CUSTOM_TYPES) {
		// initialize struct with variable
		var ptrs int = 0
		if current_token.type_ == TT_MUL {
			ptrs = getPointers()
		}
		var_name := current_token.value
		p_advance()
		if ptrs != 0 {
			p_make_struct_args(var_name, struct_name.value, true, false)
		} else {
			p_make_struct_args(var_name, struct_name.value, false, true)
		}
		return RefStructNode{struct_name: struct_name.value, ptrs: ptrs, var_name: var_name}
	}
	if current_token.type_ == TT_LCURL {
		p_advance()
		prev_vars := VARS
		VARS = make(map[string]Node)
		vars := p_func_parse(TT_RCURL)
		VARS = prev_vars
		STRUCTS[struct_name.value] = vars
		if typedef {
			CUSTOM_TYPES = append(CUSTOM_TYPES, struct_name.value)
		}
		return StructNode{Name: struct_name.value, Typedef: typedef, Estruct: estruct, Vars: vars}
	}
	return StructNode{}
}

// makes a function and returns it
func p_mikf() Node {
	p_advance()
	// get function name
	f_name := current_token.value

	// clear var scope and add globals
	old_vars := VARS
	VARS = make(map[string]Node)
	AddGlobToVars(VARS, GLOBALS)

	p_advance()
	if current_token.type_ == TT_LPAREN {
		p_advance()
		call_args := p_func_parse(TT_RPAREN)

		if current_token.type_ == TT_ARROW {
			p_advance()
			if StringInSlice(current_token.value, TYPES) || StringInSlice(current_token.value, CUSTOM_TYPES) {
				ret_type := current_token.value
				if ret_type != "void" {
					expect_ret = true
					func_ret_expect = ret_type
				}
				p_advance()

				if current_token.type_ == TT_LCURL {
					func_on = true
					p_advance()
					var code []Node
					for current_token.type_ != TT_EOF && current_token.type_ != TT_RCURL {
						if current_token.type_ == TT_ID || current_token.type_ == TT_MUL {
							res := mkID()
							//fmt.Println(res)
							//fmt.Println(VARS)
							code = append(code, res)
						} else if current_token.type_ == TT_DEBUG {
							code = append(code, DebugNode{})
						} else {
							NewError("ParsingError", "A function decleration, struct decleration, variable assignement or refference was expected but not found", CreateErrorLocation(), true)
						}
					}
					p_advance()
					ret_func := FunctionNode{Func_name: f_name, Arg_parse: call_args, Ret_type: ret_type, Code_block: code}
					VARS = old_vars
					AddGlobToVars(VARS, GLOBALS)
					FUNCTIONS[f_name] = ret_func
					if expect_ret {
						NewError("NoReturnError", "A return was expected but not found", CreateErrorLocation(), true)
					}
					return ret_func
				} else {
					ret_func := FunctionNode{Decl: true, Func_name: f_name, Arg_parse: call_args, Ret_type: ret_type, Code_block: []Node{}}
					VARS = old_vars
					AddGlobToVars(VARS, GLOBALS)
					FUNCTIONS[f_name] = ret_func
					return ret_func
				}
			} else {
				NewError("IllegalTypeError", "A non existing type: "+current_token.value+" was found.", CreateErrorLocation(), true)
			}
		} else {
			NewError("NoReturnTypeError", "A return type with '-> type' was expected but not found.", CreateErrorLocation(), true)
		}
	} else {
		NewError("ParentheseNotFoundError", "Expected a () after function name but found none.", CreateErrorLocation(), true)
	}

	return FunctionNode{}
}

// makes a assembly function and returns it
func p_mikas() AsmFunctionNode {
	p_advance()
	// get function name
	f_name := current_token.value

	// clear var scope and add globals
	old_vars := VARS
	VARS = make(map[string]Node)
	AddGlobToVars(VARS, GLOBALS)

	p_advance()
	if current_token.type_ == TT_LPAREN {
		p_advance()
		call_args := p_func_parse(TT_RPAREN)

		if current_token.type_ == TT_ARROW {
			p_advance()
			if StringInSlice(current_token.value, TYPES) || StringInSlice(current_token.value, CUSTOM_TYPES) {
				ret_type := current_token.value
				p_advance()

				if current_token.type_ == TT_LCURL {
					p_advance()
					if current_token.type_ == TT_STRING {
						asm_code := current_token.value
						p_advance()

						if current_token.type_ == TT_RCURL {
							p_advance()
							VARS = old_vars
							ret_func := AsmFunctionNode{Func_name: f_name, Arg_parse: call_args, Ret_type: ret_type, Asm_block: asm_code}
							FUNCTIONS[f_name] = ret_func
							return ret_func
						} else {
							NewError("CurlyBracketExpectError", "A Curly Bracket was expected after assembly code.", CreateErrorLocation(), true)
						}
					} else {
						NewError("IllegalAssemblyCodeError", "", CreateErrorLocation(), true)
					}
				} else {
					NewError("CurlyBracketExpectError", "A Curly Bracket was expected after assembly code.", CreateErrorLocation(), true)
				}
			} else {
				NewError("IllegalTypeError", "A non existing type: "+current_token.value+" was found.", CreateErrorLocation(), true)
			}
		} else {
			NewError("NoReturnTypeError", "A return type with '-> type' was expected but not found.", CreateErrorLocation(), true)
		}
	} else {
		NewError("ParentheseNotFoundError", "Expected a () after function name but found none.", CreateErrorLocation(), true)
	}
	return AsmFunctionNode{}
}

// makes a if statement and returns it
func p_if(elif bool) IfNode {
	var bool_block LiteralNode
	p_advance()
	CURRENT_LN = current_token.ln_count

	if current_token.type_ == TT_LPAREN {
		p_advance()
		bool_block = p_bool_expr()
		if current_token.type_ == TT_RPAREN {
			p_advance()
			if current_token.type_ == TT_LCURL {
				p_advance()
				var code []Node
				for current_token.type_ != TT_EOF && current_token.type_ != TT_RCURL {
					if current_token.type_ == TT_ID || current_token.type_ == TT_MUL {
						res := mkID()
						//fmt.Println(res)
						//fmt.Println(VARS)
						code = append(code, res)
					} else if current_token.type_ == TT_DEBUG {
						code = append(code, DebugNode{})
					} else {
						NewError("ParsingError", "A function decleration, struct decleration, variable assignement or refference was expected but not found", CreateErrorLocation(), true)
					}
				}
				p_advance()
				ret_if := IfNode{elif: elif, bool_: bool_block, codeblock: code}
				PREV_IF = true
				return ret_if
			} else {
				NewError("CodeBlockExpectedError", "A { was expexted after if (bool) ...", CreateErrorLocation(), true)
			}
		} else {
			NewError("ParantheseExpectedError", "After bool condition a ')' was expected but not found.", CreateErrorLocation(), true)
		}
	} else {
		NewError("ParantheseExpectedError", "After 'if' condition a '(' was expected but not found.", CreateErrorLocation(), true)
	}
	return IfNode{}
}

// makes a else statemtn and returns it
func p_else() ElseNode {
	p_advance()

	if current_token.type_ == TT_LCURL {
		p_advance()
		var code []Node
		for current_token.type_ != TT_EOF && current_token.type_ != TT_RCURL {
			if current_token.type_ == TT_ID || current_token.type_ == TT_MUL {
				res := mkID()
				//fmt.Println(res)
				//fmt.Println(VARS)
				code = append(code, res)
			} else if current_token.type_ == TT_DEBUG {
				code = append(code, DebugNode{})
			} else {
				NewError("ParsingError", "A function decleration, struct decleration, variable assignement or refference was expected but not found", CreateErrorLocation(), true)
			}
		}
		p_advance()
		ret_else := ElseNode{codeblock: code}
		return ret_else
	} else {
		NewError("CodeBlockExpectedError", "A code block was expecred after 'else'.", CreateErrorLocation(), true)
	}
	return ElseNode{}
}

// makes a while statement and returns it
func p_while() WhileNode {
	LOOP = true
	var bool_block LiteralNode
	p_advance()
	CURRENT_LN = current_token.ln_count

	if current_token.type_ == TT_LPAREN {
		p_advance()
		bool_block = p_bool_expr()
		if current_token.type_ == TT_RPAREN {
			p_advance()
			if current_token.type_ == TT_LCURL {
				p_advance()
				var code []Node
				for current_token.type_ != TT_EOF && current_token.type_ != TT_RCURL {
					if current_token.type_ == TT_ID || current_token.type_ == TT_MUL {
						res := mkID()
						//fmt.Println(res)
						//fmt.Println(VARS)
						code = append(code, res)
					} else if current_token.type_ == TT_DEBUG {
						code = append(code, DebugNode{})
					} else {
						NewError("ParsingError", "A function decleration, struct decleration, variable assignement or refference was expected but not found", CreateErrorLocation(), true)
					}
				}
				p_advance()
				ret_while := WhileNode{bool_: bool_block, codeblock: code}
				return ret_while
			} else {
				NewError("CodeBlockExpectedError", "A '{' was expexted after while (bool) ...", CreateErrorLocation(), true)
			}
		} else {
			NewError("ParantheseExpectedError", "After bool condition a ')' was expected but not found.", CreateErrorLocation(), true)
		}
	} else {
		NewError("ParantheseExpectedError", "After 'while' condition a '(' was expected but not found.", CreateErrorLocation(), true)
	}
	return WhileNode{}
}

/*-------------------expr, factor and binary op-------------------*/
func mkID() Node {
	var node Node
	tok := current_token
	// check if it is a assignement by looking for a dtype
	//fmt.Println(tok)
	if expect_doc {
		if !(StringInSlice(tok.value, []string{"mikf", "mikas", "struct", "estruct"}) || StringInSlice(tok.value, TYPES) || StringInSlice(tok.value, CUSTOM_TYPES) || tok.value == "global") {
			NewError("NoDocExpectedError", "A doc is only expexted before functions, variables, structs, and assembly functions", CreateErrorLocation(), true)
		} else {
			expect_doc = false
		}

	}
	if PREV_IF && tok.value != "elif" && tok.value != "else" && tok.value != "if" {
		PREV_IF = false
	}
	if StringInSlice(tok.value, TYPES) || StringInSlice(tok.value, CUSTOM_TYPES) {
		node = p_assign(tok.value, global)
	} else if tok.value == "global" && !global {
		global = true
		p_advance()
		node = mkID()
	} else if tok.value == "return" {
		// check if return is possible
		if func_on && expect_ret {
			p_advance()
			CURRENT_LN = current_token.ln_count
			current_expected_t = func_ret_expect
			expect_t = true
			ret := p_expr()
			expect_ret = false
			expect_t = false
			node = ReturnNode{Return_val: ret}
		} else {
			NewError("ReturnNotExpected", "You tried to return without being in a function.", CreateErrorLocation(), true)
		}
	} else if tok.value == "break" && LOOP {
		node = DirectNode{"", "break", false, false}
		p_advance()
		LOOP = false
	} else if tok.value == "continue" && LOOP {
		node = DirectNode{"", "continue", false, false}
		p_advance()
		LOOP = false
	} else if tok.value == "mikf" {
		node = p_mikf()
		//fmt.Println(node)
	} else if tok.value == "mikas" {
		node = p_mikas()
	} else if tok.value == "struct" {
		node = p_struct(false)
	} else if tok.value == "estruct" {
		node = p_struct(true)
	} else if tok.value == "if" {
		node = p_if(false)
	} else if tok.value == "elif" {
		if !PREV_IF {
			NewError("UnexpectedOperationError", "No elif is expected without a if to even begin with.", CreateErrorLocation(), true)
		} else {
			node = p_if(true)
		}
	} else if tok.value == "else" {
		if !PREV_IF {
			NewError("UnexpectedOperationError", "No else is expected without a if or elif to even begin with.", CreateErrorLocation(), true)
		} else {
			node = p_else()
			PREV_IF = false
		}
	} else if tok.value == "while" {
		node = p_while()
	} else if StringInMap(getName().value, VARS) || StringInSlice(getName().value, var_names) || StringInMap(getName().value, FUNCTIONS) || StringInSlice(getName().value, func_names) || StringInMap(getName().value, GLOBALS) || tok.type_ == TT_MUL {
		if tok.type_ == TT_MUL {
			node = p_reassign(true)
		} else {
			node = p_reassign(false)
		}
	} else if tok.type_ == TT_DOC {
		node = DocNode{current_token.value}
		p_advance()
		expect_doc = true
	} else {
		fmt.Println(current_token)
		NewError("", "", "", true)
	}
	return node
}
