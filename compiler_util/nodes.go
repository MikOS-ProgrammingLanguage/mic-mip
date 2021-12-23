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
	Is_1_node() bool
	What_type() string
}
type SubNode interface {
	Is_2_node() bool
	What_type() string
}
type LiteralNode interface {
	Is_3_node() bool
	What_type() string
}

type UniversalNone struct {
}
type RootNode struct {
	Nodes []Node
}

// adds a given node to the root node stem
func (rnode RootNode) AddNodeToRoot(node Node) RootNode {
	rnode.Nodes = append(rnode.Nodes, node)
	return rnode
}

/* First class nodes */
type IfNode struct {
	elif      bool
	bool_     LiteralNode
	codeblock []Node
}
type ElseNode struct {
	codeblock []Node
}
type ForNode struct {
	// how the fuck
	// codeblock
}
type WhileNode struct {
	// bool statement
	// codeblock
}
type ReturnNode struct {
	Return_val LiteralNode
}
type ReAssignmentNode struct {
	Reassgn_t string
	Re_type   string
	Ptrs      int
	Content   LiteralNode
}
type AssignemntNode struct {
	Asgn_type string
	Ptrs      int
	Global    bool
	Var_name  string
	Content   LiteralNode
}
type ArrReAssignementNode struct {
	Reassgn_t string
	Re_type   string
	Ptrs      int
	Arr_idx   LiteralNode
	Content   LiteralNode
}
type ArrAssignementNode struct {
	Asgn_type  string
	Ptrs       int
	Global     bool
	Array_name string
	Arr_len    LiteralNode
}
type FunctionNode struct {
	Decl       bool
	Func_name  string
	Arg_parse  []Node
	Ret_type   string
	Code_block []Node
}
type AsmFunctionNode struct {
	Func_name string
	Arg_parse []Node
	Ret_type  string
	Asm_block string
}
type StructNode struct {
	Name    string
	Typedef bool
	Estruct bool
	Vars    []Node
}
type DebugNode struct {
}
type RefStructNode struct {
	struct_name string
	ptrs        int
	var_name    string
}

/* Second class nodes */
type BinOpNode struct {
	Left_node  LiteralNode
	Op_tok     string
	Right_node LiteralNode
}
type TypeCastNode struct {
	Tcast LiteralNode
	Dtype LiteralNode
}
type BoolOpNode struct {
	left   LiteralNode
	op_tok string
	right  LiteralNode
}

// Third class nodes
type ListSliceNode struct {
	Name  string
	Pos   LiteralNode
	Ptrs  int
	Deref bool
	Not   bool
}
type VarNameNode struct {
	Name  string
	Ptrs  int
	Deref bool
	Not   bool
}
type DataTypeNode struct {
	Dtype string
	Ptrs  int
}
type FuncCallNode struct {
	Call_name  string
	Func_parse []LiteralNode // Subnode arr bcs -> call_name(expr(), expr(), ...n)
}

type DirectNode struct {
	Type_ string
	Value string
}

// -------------- 1 Interface implement --------------

// implements is_1_node for UniversalNone
func (unon UniversalNone) Is_1_node() bool {
	return true
}
func (unon UniversalNone) What_type() string {
	return "UinversalNone"
}

// implements is_1_node
func (ifn IfNode) Is_1_node() bool {
	return true
}
func (ifn IfNode) What_type() string {
	return "IfNode"
}

// implements is_1_node
func (elsen ElseNode) Is_1_node() bool {
	return true
}
func (elsen ElseNode) What_type() string {
	return "ElseNode"
}

// implements is_1_node
func (forn ForNode) Is_1_node() bool {
	return true
}
func (forn ForNode) What_type() string {
	return "ForNode"
}

// implements is_1_node
func (whilen WhileNode) Is_1_node() bool {
	return true
}
func (whilen WhileNode) What_type() string {
	return "WhileNode"
}

// implements is_1_node
func (struct_ref RefStructNode) Is_1_node() bool {
	return true
}
func (struct_red RefStructNode) What_type() string {
	return "RefStructNode"
}

// implements is_1_node for FuncCall
func (fcall FuncCallNode) Is_1_node() bool {
	return true
}
func (fcall FuncCallNode) What_type() string {
	return "FuncCallNode"
}

// implements is_1_node for Return
func (ret ReturnNode) Is_1_node() bool {
	return true
}
func (ret ReturnNode) What_type() string {
	return "ReturnNode"
}

// implements is_1_node for Assignement
func (ass AssignemntNode) Is_1_node() bool {
	return true
}
func (ass AssignemntNode) What_type() string {
	return "AssignementNode"
}

// implements is_1_node for ReAssignement
func (reass ReAssignmentNode) Is_1_node() bool {
	return true
}
func (reass ReAssignmentNode) What_type() string {
	return "ReAssignementNode"
}

// implements is_1_node for ArrAssignement
func (arrassgn ArrAssignementNode) Is_1_node() bool {
	return true
}
func (arrassgn ArrAssignementNode) What_type() string {
	return "ArrAssignementNode"
}

// implements is_1_node for ArrReAssignement
func (arrreassgn ArrReAssignementNode) Is_1_node() bool {
	return true
}
func (arrreassgn ArrReAssignementNode) What_type() string {
	return "ArrReAssignementNode"
}

// implements is_1_node for function
func (foo FunctionNode) Is_1_node() bool {
	return true
}
func (foo FunctionNode) What_type() string {
	return "FunctionNode"
}

// implements is_1_node for function
func (asm AsmFunctionNode) Is_1_node() bool {
	return true
}
func (asm AsmFunctionNode) What_type() string {
	return "AsmFunctionNode"
}
func Arg_len(fn Node) int {
	if fn.What_type() == "AsmFunctionNode" {
		temp := fn.(AsmFunctionNode)
		return len(temp.Arg_parse)
	} else {
		temp := fn.(FunctionNode)
		return len(temp.Arg_parse)
	}
}

// implement is_1_node for debug
func (debg DebugNode) Is_1_node() bool {
	return true
}
func (debg DebugNode) What_type() string {
	return "DegubNode"
}

// implement is_1_node for struct
func (strct StructNode) Is_1_node() bool {
	return true
}
func (strct StructNode) What_type() string {
	return "StructNode"
}

// -------------- 2 Interface implement --------------

// implements is_1_node for UniversalNone
func (unon UniversalNone) Is_2_node() bool {
	return true
}

// implements is_2_node for BinOp
func (bio BinOpNode) Is_2_node() bool {
	return true
}
func (bio BinOpNode) What_type() string {
	return "BinOpNode"
}

// implements is_2_node for BoolOp
func (boo BoolOpNode) Is_2_node() bool {
	return true
}
func (boo BoolOpNode) What_type() string {
	return "BoolOpNode"
}

// -------------- 3 Interface implement --------------

// implements is_1_node for UniversalNone
func (unon UniversalNone) Is_3_node() bool {
	return true
}

// implementrs is_3_node for DataTypeNode
func (dtype DataTypeNode) Is_3_node() bool {
	return true
}
func (dtype DataTypeNode) What_type() string {
	return "DataTypeNode"
}

// implements is_3_node for bool op
func (boo BoolOpNode) Is_3_node() bool {
	return true
}

// implements is_3_node for FuncCallNode
func (f_call FuncCallNode) Is_3_node() bool {
	return true
}

// implements is_3_node for BinOpNode
func (binOp BinOpNode) Is_3_node() bool {
	return true
}

// implements is_3_node for VarNameNode
func (vname VarNameNode) Is_3_node() bool {
	return true
}
func (vname VarNameNode) What_type() string {
	return "VarNameNode"
}

// implements is_3_node for ListSliceNode
func (lslice ListSliceNode) Is_3_node() bool {
	return true
}
func (lslice ListSliceNode) What_type() string {
	return "ListSliceNode"
}

func (tcst TypeCastNode) Is_3_node() bool {
	return true
}
func (tcst TypeCastNode) What_type() string {
	return "TypeCastNode"
}

func (drct DirectNode) Is_3_node() bool {
	return true
}
func (drct DirectNode) What_type() string {
	return "DirectNode"
}
