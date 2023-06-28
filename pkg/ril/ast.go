package ril

// ASTNode is a JSON-serializable AST node.
type ASTNode struct {
	Func      string     `json:"func"`
	Arguments any        `json:"arguments"`
	Children  []*ASTNode `json:"children"`
}
