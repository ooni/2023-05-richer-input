package uncompiler

import "encoding/json"

// ASTNode is a JSON-deserializable AST node.
type ASTNode struct {
	// Func is the name of the [unruntime] function to execute.
	Func string `json:"func"`

	// Arguments contains function-specific arguments.
	Arguments json.RawMessage `json:"arguments"`

	// Children contains function-specific children nodes.
	Children []*ASTNode `json:"children"`
}
