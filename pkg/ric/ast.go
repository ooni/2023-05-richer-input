package ric

import "encoding/json"

// ASTNode is a JSON-deserializable AST node.
//
// The Func is the name of the function to execute. You should use [ridsl.IfFuncExists]
// if a function may not exist on all deployed OONI Probe versions.
//
// The Arguments contain function-specific arguments, which may possibly be `null`
// in the serialized JSON. This field is always a possibly-empty structure. Because it is
// a structure, we can implement extensions in new OONI Probe versions as long as the
// zero value of new fields implements the previous behavior.
//
// The Children contains the children nodes in the AST.
type ASTNode struct {
	Func      string          `json:"func"`
	Arguments json.RawMessage `json:"arguments"`
	Children  []*ASTNode      `json:"children"`
}
