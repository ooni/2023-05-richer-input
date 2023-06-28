package undsl

// ASTNode is a JSON-serializable AST node. The corresponding [uncompiler] data
// structure is [*uncompiler.ASTNode]. You SHOULD use factories exported by this
// package to create [*ASTNode] instances. If you need to manually create a new
// [*ASTNode], make sure you fill the fields marked as MANDATORY.
type ASTNode struct {
	// Func is the MANDATORY name of the [unruntime] function to execute.
	Func string `json:"func"`

	// Arguments contains Func-specific OPTIONAL arguments to construct the
	// corresponding [unruntime] function.
	Arguments any `json:"arguments"`

	// Children contains OPTIONAL child nodes.
	Children []*ASTNode `json:"children"`
}
