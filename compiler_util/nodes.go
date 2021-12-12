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
func (rnode RootNode) AddNodeToRoot(node Node) {
	rnode.nodes = append(rnode.nodes, node)
}

/* First class nodes */
type ReAssignmentNode struct {
	re_type string
	content LiteralNode
}
type AssignemntNode struct {
	asgn_type string
	ptrs      int
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
	array_name string
	arr_len    LiteralNode
}
type FunctionNode struct {
	func_name  string
	arg_parse  SubNode
	ret_type   LiteralNode
	code_block Node
}
type AsmFunctionNode struct {
	func_name string
	arg_parse SubNode
	ret_type  LiteralNode
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
type FuncParseNode struct {
	parsed []AssignemntNode // bcs -> mikf name(int i, int b) {}
}
type TypeCastNode struct {
	tcast LiteralNode
	dtype string
}

// Third class nodes
type ListSliceNode struct {
	name string
	pos  LiteralNode
}
type VarNameNode struct {
	name string
}
type DataTypeNode struct {
	dtype string
}
type FuncCallNode struct {
	call_name  string
	func_parse []SubNode // Subnode arr bcs -> call_name(expr(), expr(), ...n)
}

// -------------- 1 Interface implement --------------

// implements is_1_node for UniversalNone
func (unon UniversalNone) is_1_node() bool {
	return true
}
func (unon UniversalNone) what_type() string {
	return "UinversalNone"
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

// implements is_2_node for FuncParseNode
func (fparse FuncParseNode) is_2_node() bool {
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
