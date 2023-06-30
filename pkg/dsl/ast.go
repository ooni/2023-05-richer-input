package dsl

// ASTNode is an AST node of the DSL.
type ASTNode struct {
	// Func is the name of the DSL function to execute.
	Func string `json:"func"`

	// Arguments contains function-specific arguments.
	Arguments any `json:"arguments"`

	// Children contains function-specific children nodes.
	Children []*ASTNode `json:"children"`
}
