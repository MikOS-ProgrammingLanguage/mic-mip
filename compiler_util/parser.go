package compiler_util

import (
	"fmt"
)

var tokens []Token
var illegal_names []string
var pos int = -1
var current_token Token
var is_eot = false
var root_node RootNode = RootNode{}
var func_on bool = false
var VARS map[string]Node = make(map[string]Node)        // holds assigned vars
var FUNCTIONS map[string]Node = make(map[string]Node)   // holds all functions that have a body
var STRUCTS map[string][]Node = make(map[string][]Node) // holds all structs and their vars

// parse that shit wohooooooo
func Parse(Tokens *[]Token, Illegal_names []string) RootNode {
	tokens = *Tokens
	illegal_names = Illegal_names
	p_advance()

	for current_token.type_ != TT_EOF {
		if current_token.type_ == TT_ID {
			res := mkID()
			//fmt.Println(res)
			//fmt.Println(VARS)
			root_node = root_node.AddNodeToRoot(res)
		} else if current_token.type_ == TT_DEBUG {
			root_node = root_node.AddNodeToRoot(DebugNode{})
		} else {
			NewError("ParsingError", "A function decleration, struct decleration, variable assignement or refference was expected but not found", fmt.Sprintf("%s at ln %d", current_token.section, current_token.ln_count), true)
		}
	}
	return root_node
}

/*-------------------Get and increment methods-------------------*/

// increments the pos and sets the token if the position still is in the index of tokens
func p_advance() {
	pos++
	if pos < len(tokens) {
		current_token = tokens[pos]
	} else {
		current_token = Token{}
		is_eot = true
	}
}

// gets a token in the future by index for example -> getToken(1) would give you the next token
func getToken(num int) Token {
	pos_now := pos
	if pos_now+num < len(tokens) {
		return tokens[pos_now+num]
	} else {
		is_eot = true
		return Token{}
	}
}

// tries to make a name for example -> hello.id is legal
func getName() Token {
	name_node := current_token
	if getToken(1).type_ == TT_DOT && getToken(2).type_ == TT_ID && (!StringInSlice(getToken(2).value, TYPES) && !StringInSlice(getToken(2).value, INSTRUCTIONS) && !StringInSlice(getToken(2).value, CUSTOM_TYPES)) {
		// iterates through the tokens for as long as it's a continous name with dots like -> hello.type
		for getToken(1).type_ == TT_DOT && getToken(2).type_ == TT_ID && (!StringInSlice(getToken(2).value, TYPES) && !StringInSlice(getToken(2).value, INSTRUCTIONS) && !StringInSlice(getToken(2).value, CUSTOM_TYPES)) {
			p_advance()
			name_node.value += "."
			p_advance()
			name_node.value += current_token.value
		}
	} else if getToken(1).type_ == TT_ARROW && getToken(2).type_ == TT_ID && (!StringInSlice(getToken(2).value, TYPES) && !StringInSlice(getToken(2).value, INSTRUCTIONS) && !StringInSlice(getToken(2).value, CUSTOM_TYPES)) {
		// iterates through the tokens for as long as it's a continous name with arrows like -> hello->type
		for getToken(1).type_ == TT_ARROW && getToken(2).type_ == TT_ID && (!StringInSlice(getToken(2).value, TYPES) && !StringInSlice(getToken(2).value, INSTRUCTIONS) && !StringInSlice(getToken(2).value, CUSTOM_TYPES)) {
			p_advance()
			name_node.value += "->"
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
			NewError("ArgumentExpectedError", "After a comma in struct a value is expected but there was none!", fmt.Sprintf("%s at ln %d", next_next_tok.section, next_next_tok.ln_count), true)
		}
		if StringInSlice(tok.value, TYPES) || StringInSlice(tok.value, CUSTOM_TYPES) {
			node = p_assign(tok.value)
		} else if tok.type_ == TT_COMMA {
			p_advance()
			continue
		} else {
			NewError("IvalidTypeError", "A type that doesn't exist was specified", fmt.Sprintf("%s at ln %d", tok.section, tok.ln_count), true)
		}
		vars = append(vars, node)
	}
	p_advance()
	return vars
}

// makes struct args of of a list in STRUCTS
func p_make_struct_args(name, type_ string, is_ptr, is_n_ptr bool) {
	str_get := STRUCTS[type_]
	for i := 0; len(str_get) > i; i++ {
		current_struct_attr := str_get[i]
		if current_struct_attr.what_type() == "AssignementNode" {
			current_struct_attr2 := current_struct_attr.(AssignemntNode)
			if StringInSlice(current_struct_attr2.asgn_type, CUSTOM_TYPES) {
				if is_ptr {
					VARS[name+"."+current_struct_attr2.var_name] = AssignemntNode{current_struct_attr2.asgn_type, current_struct_attr2.ptrs, name + "." + current_struct_attr2.var_name, current_struct_attr2.content}
					p_make_struct_args(name+"->"+current_struct_attr2.var_name, current_struct_attr2.asgn_type, true, false)
				} else {
					VARS[name+"->"+current_struct_attr2.var_name] = AssignemntNode{current_struct_attr2.asgn_type, current_struct_attr2.ptrs, name + "->" + current_struct_attr2.var_name, current_struct_attr2.content}
					p_make_struct_args(name+"."+current_struct_attr2.var_name, current_struct_attr2.asgn_type, false, true)
				}
			} else {
				if is_n_ptr {
					VARS[name+"."+current_struct_attr2.var_name] = AssignemntNode{current_struct_attr2.asgn_type, current_struct_attr2.ptrs, name + "." + current_struct_attr2.var_name, current_struct_attr2.content}
				}
				if is_ptr {
					VARS[name+"->"+current_struct_attr2.var_name] = AssignemntNode{current_struct_attr2.asgn_type, current_struct_attr2.ptrs, name + "->" + current_struct_attr2.var_name, current_struct_attr2.content}
				}
			}
		} else if current_struct_attr.what_type() == "ArrAssignementNode" {
			current_struct_attr2 := current_struct_attr.(ArrAssignementNode)
			if StringInSlice(current_struct_attr2.asgn_type, CUSTOM_TYPES) {
				if is_n_ptr {
					p_make_struct_args(name+"."+current_struct_attr2.array_name, current_struct_attr2.asgn_type, false, true)
					VARS[name+"."+current_struct_attr2.array_name] = ArrAssignementNode{current_struct_attr2.asgn_type, current_struct_attr2.ptrs, name + "." + current_struct_attr2.array_name, current_struct_attr2.arr_len}
				}
				if is_ptr {
					p_make_struct_args(name+"->"+current_struct_attr2.array_name, current_struct_attr2.asgn_type, true, false)
					VARS[name+"->"+current_struct_attr2.array_name] = AssignemntNode{current_struct_attr2.asgn_type, current_struct_attr2.ptrs, name + "->" + current_struct_attr2.array_name, current_struct_attr2.arr_len}
				}
			} else {
				if is_n_ptr {
					VARS[name+"."+current_struct_attr2.array_name] = ArrAssignementNode{current_struct_attr2.asgn_type, current_struct_attr2.ptrs, name + "." + current_struct_attr2.array_name, current_struct_attr2.arr_len}
				}
				if is_ptr {
					VARS[name+"->"+current_struct_attr2.array_name] = ArrAssignementNode{current_struct_attr2.asgn_type, current_struct_attr2.ptrs, name + "->" + current_struct_attr2.array_name, current_struct_attr2.arr_len}
				}
			}
		} else {
			NewError("WTF", "", "", true)
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
	var ptrs int = 0
	var deref bool = false
	// returns either a func_call, var_name or array_slcie
	if tok.type_ == TT_MUL {
		ptrs = getPointers()
	} else if tok.type_ == TT_KAND {
		deref = true
		p_advance()
	}
	tok = current_token
	if tok.type_ == TT_ID {
		tok = getName()
		p_advance()
		if StringInMap(tok.value, VARS) {
			// return list slice if [] is found
			if current_token.type_ == TT_LBRK {
				p_advance()
				arr_pos := p_expr()
				if current_token.type_ == TT_RBRK {
					p_advance()
					return ListSliceNode{name: tok.value, pos: arr_pos, ptrs: ptrs, deref: deref}
				} else {
					NewError("ArrayNotClosed", "A array was assagnid with '[' but no ']' was found. ", fmt.Sprintf("%s at ln %d", current_token.section, current_token.ln_count), true)
				}
			}
			return VarNameNode{tok.value, ptrs, deref}
		} else if StringInMap(tok.value, FUNCTIONS) {
			// make and return a functioncall
		} else {
			fmt.Println(tok.value)
			NewError("RefferenceError", "The refferenced var, func or struct is not defined", fmt.Sprintf("%s at ln %d", tok.section, tok.ln_count), true)
		}
	} else if tok.type_ == TT_INT {
		// return a node for either int, flt, str or char
	}
	p_advance()
	return DataTypeNode{tok.value, ptrs}
}

func p_term() LiteralNode {
	ops := []string{TT_MUL, TT_DIV}
	return binOp(true, false, ops)
}
func p_expr() LiteralNode {
	ops := []string{TT_PLUS, TT_MINUS}
	return binOp(false, true, ops)
}

// makes a binOp node and returns it
func binOp(factor, term bool, ops []string) LiteralNode {
	var left LiteralNode
	if term {
		left = p_term()
	} else if factor {
		left = p_factor()
	}

	for StringInSlice(current_token.type_, ops) && !is_eot {
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
func p_assign(type_ string) Node {
	p_advance()

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
	if StringInMap(assignement_name.value, VARS) {
		NewError("VariableAllreadyAssigned", fmt.Sprintf("The Variable \"%s\" is allready assigned.", assignement_name.value), fmt.Sprintf("%s at ln %d", assignement_name.section, assignement_name.ln_count), true)
	}

	p_advance()
	// returns a array or a assignement based on if there are brackets or not
	if current_token.type_ == TT_LBRK {
		p_advance()
		var arr_len LiteralNode
		if current_token.type_ != TT_RBRK {
			arr_len = p_expr()
		} else {
			// array needs to be initialized with zeros in the c file -> type name[] = 0
			arr_len = UniversalNone{}
		}

		// checks if the arr expr is closed
		if current_token.type_ == TT_RBRK {
			p_advance()
			var_node := ArrAssignementNode{asgn_type: type_, ptrs: ptrs, array_name: assignement_name.value, arr_len: arr_len}
			VARS[assignement_name.value] = var_node
			return var_node
		} else {
			NewError("ArrayNotClosed", "A array was assagnid with '[' but no ']' was found. ", fmt.Sprintf("%s at ln %d", current_token.section, current_token.ln_count), true)
		}
	} else if current_token.type_ == TT_ASSGN {
		p_advance()
		// make a typecast node if one is found
		if current_token.type_ == TT_ID && current_token.value == "tcst" {
			p_advance()
			if current_token.type_ == TT_LPAREN {
				p_advance()
				var_name := p_expr()
				if current_token.type_ == TT_COMMA {
					p_advance()
					if StringInSlice(current_token.value, TYPES) || StringInSlice(current_token.value, CUSTOM_TYPES) {
						type2 := getName()
						p_advance()
						if current_token.type_ == TT_RPAREN {
							p_advance()
							ret_var := AssignemntNode{asgn_type: type_, ptrs: ptrs, var_name: assignement_name.value, content: TypeCastNode{tcast: var_name, dtype: type2.value}}
							return ret_var
						}
					} else {
						NewError("TypeExpectedButNotFound", "A type in Typecast was expected but not found", fmt.Sprintf("%s at ln %d", current_token.section, current_token.ln_count), true)
					}
				} else {
					NewError("TypeExpectedButNotFound", "A type in Typecast was expected but not found", fmt.Sprintf("%s at ln %d", current_token.section, current_token.ln_count), true)
				}
			} else {
				NewError("ParantheseExpectes", "A parenthese in a typecast was expected but nor found", fmt.Sprintf("%s at ln %d", current_token.section, current_token.ln_count), true)
			}
		} else {
			content := p_expr()
			var_node := AssignemntNode{asgn_type: type_, ptrs: ptrs, var_name: assignement_name.value, content: content}
			VARS[assignement_name.value] = var_node
			return var_node
		}
	} else {
		var_node := AssignemntNode{asgn_type: type_, ptrs: ptrs, var_name: assignement_name.value, content: UniversalNone{}}
		VARS[assignement_name.value] = var_node
		return var_node
	}
	return AssignemntNode{}
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
		return StructNode{name: struct_name.value, typedef: typedef, estruct: estruct, vars: vars}
	}
	return StructNode{}
}

// makes a function and returns it
func p_mikf() FunctionNode {

	return FunctionNode{}
}

// makes a assembly function and returns it
func p_mikas() AsmFunctionNode {
	return AsmFunctionNode{}
}

/*-------------------expr, factor and binary op-------------------*/
func mkID() Node {
	var node Node
	tok := current_token
	// check if it is a assignement by looking for a dtype
	if StringInSlice(tok.value, TYPES) || StringInSlice(tok.value, CUSTOM_TYPES) {
		node = p_assign(tok.value)
	} else if tok.value == "return" {
		// check if return is possible
		if func_on {
			// make return statement
		} else {
			NewError("ReturnNotExpected", "You tried to return without being in a function.", fmt.Sprintf("%s at ln %d", tok.section, tok.ln_count), true)
		}
	} else if tok.value == "mikf" {
		node = p_mikf()
	} else if tok.value == "mikas" {
		node = p_mikas()
	} else if tok.value == "struct" {
		node = p_struct(false)
	} else if tok.value == "estruct" {
		node = p_struct(true)
	} else if tok.value == "if" {
		// make a if statement
	} else if tok.value == "elif" {
		// make a elif
	} else if tok.value == "else" {
		// make else
	} else if tok.value == "while" {
		// make while
	} else if tok.value == "for" {
		// make for
	} else {
		// look for function call or reassignement
		fmt.Println(current_token)
		NewError("", "", "", true)
	}
	return node
}
