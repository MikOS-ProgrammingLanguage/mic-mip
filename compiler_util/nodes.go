package compiler_util

// thirs class returns to second class which returns to first class and first class appends to root directly. 3 -> 2 -> 1 -> 0
/* all first class types. First class means that they can be independent for example a bin op is only called with assignements hence it would be second class

- Assignemnt
- ReAssignement
- Function
- Assembly function
- Struct
- Estruct

*/
/* all second class types

- BinOpNode
- FuncParse
- TypeCast

*/
/* all third class types

- datatypeNode
- funcCallNode
*/

// a node is everything that implements is_1_node with ret bool
// is_1_node is needed to distinguish between nodes of different prioritys
type Node interface {
	is_1_node() bool
	what_type() string
}
type SubNode interface {
	is_2_node() bool
}
type LiteralNode interface {
	is_3_node() bool
}

type UniversalNone struct {
}
type RootNode struct {
	nodes []Node
}

// adds a given node to the root node stem
func (rnode RootNode) AddNodeToRoot(node Node) RootNode {
	rnode.nodes = append(rnode.nodes, node)
	return rnode
}

/* First class nodes */
type ReturnNode struct {
	return_val LiteralNode
}
type ReAssignmentNode struct {
	re_type string
	content LiteralNode
}
type AssignemntNode struct {
	asgn_type string
	ptrs      int
	global    bool
	var_name  string
	content   LiteralNode
}
type ArrReAssignementNode struct {
	re_type string
	content LiteralNode
}
type ArrAssignementNode struct {
	asgn_type  string
	ptrs       int
	global     bool
	array_name string
	arr_len    LiteralNode
}
type FunctionNode struct {
	decl       bool
	func_name  string
	arg_parse  []Node
	ret_type   string
	code_block []Node
}
type AsmFunctionNode struct {
	func_name string
	arg_parse []Node
	ret_type  string
	asm_block string
}
type StructNode struct {
	name    string
	typedef bool
	estruct bool
	vars    []Node
}
type DebugNode struct {
}

/* Second class nodes */
type BinOpNode struct {
	left_node  LiteralNode
	op_tok     string
	right_node LiteralNode
}
type TypeCastNode struct {
	tcast LiteralNode
	dtype LiteralNode
}

// Third class nodes
type ListSliceNode struct {
	name  string
	pos   LiteralNode
	ptrs  int
	deref bool
}
type VarNameNode struct {
	name  string
	ptrs  int
	deref bool
}
type DataTypeNode struct {
	dtype string
	ptrs  int
}
type FuncCallNode struct {
	call_name  string
	func_parse []LiteralNode // Subnode arr bcs -> call_name(expr(), expr(), ...n)
}

type DirectNode struct {
	type_ string
	value string
}

// -------------- 1 Interface implement --------------

// implements is_1_node for UniversalNone
func (unon UniversalNone) is_1_node() bool {
	return true
}
func (unon UniversalNone) what_type() string {
	return "UinversalNone"
}

// implements is_1_node for Return
func (ret ReturnNode) is_1_node() bool {
	return true
}
func (ret ReturnNode) what_type() string {
	return "ReturnNode"
}

// implements is_1_node for Assignement
func (ass AssignemntNode) is_1_node() bool {
	return true
}
func (ass AssignemntNode) what_type() string {
	return "AssignementNode"
}

// implements is_1_node for ReAssignement
func (reass ReAssignmentNode) is_1_node() bool {
	return true
}
func (reass ReAssignmentNode) what_type() string {
	return "ReAssignementNode"
}

// implements is_1_node for ArrAssignement
func (arrassgn ArrAssignementNode) is_1_node() bool {
	return true
}
func (arrassgn ArrAssignementNode) what_type() string {
	return "ArrAssignementNode"
}

// implements is_1_node for ArrReAssignement
func (arrreassgn ArrReAssignementNode) is_1_node() bool {
	return true
}
func (arrreassgn ArrReAssignementNode) what_type() string {
	return "ArrReAssignementNode"
}

// implements is_1_node for function
func (foo FunctionNode) is_1_node() bool {
	return true
}
func (foo FunctionNode) what_type() string {
	return "FunctionNode"
}

// implements is_1_node for function
func (asm AsmFunctionNode) is_1_node() bool {
	return true
}
func (asm AsmFunctionNode) what_type() string {
	return "AsmFunctionNode"
}
func arg_len(fn Node) int {
	if fn.what_type() == "AsmFunctionNode" {
		temp := fn.(AsmFunctionNode)
		return len(temp.arg_parse)
	} else {
		temp := fn.(FunctionNode)
		return len(temp.arg_parse)
	}
	return 1
}

// implement is_1_node for debug
func (debg DebugNode) is_1_node() bool {
	return true
}
func (debg DebugNode) what_type() string {
	return "DegubNode"
}

// implement is_1_node for struct
func (strct StructNode) is_1_node() bool {
	return true
}
func (strct StructNode) what_type() string {
	return "StructNode"
}

// -------------- 2 Interface implement --------------

// implements is_1_node for UniversalNone
func (unon UniversalNone) is_2_node() bool {
	return true
}

// implements is_2_node for BinOp
func (bio BinOpNode) is_2_node() bool {
	return true
}

// -------------- 3 Interface implement --------------

// implements is_1_node for UniversalNone
func (unon UniversalNone) is_3_node() bool {
	return true
}

// implementrs is_3_node for DataTypeNode
func (dtype DataTypeNode) is_3_node() bool {
	return true
}

// implements is_3_node for FuncCallNode
func (f_call FuncCallNode) is_3_node() bool {
	return true
}

// implements is_3_node for BinOpNode
func (binOp BinOpNode) is_3_node() bool {
	return true
}

// implements is_3_node for VarNameNode
func (vname VarNameNode) is_3_node() bool {
	return true
}

// implements is_3_node for ListSliceNode
func (lslice ListSliceNode) is_3_node() bool {
	return true
}

func (tcst TypeCastNode) is_3_node() bool {
	return true
}

func (drct DirectNode) is_3_node() bool {
	return true
}
