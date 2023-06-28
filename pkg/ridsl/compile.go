package ridsl

// Compile converts a [*Func] to an [*ASTNode].
func Compile(f *Func) *ASTNode {
	children := []*ASTNode{}
	for _, entry := range f.Children {
		children = append(children, Compile(entry))
	}
	return &ASTNode{
		Func:      f.Name,
		Arguments: f.Arguments,
		Children:  children,
	}
}
