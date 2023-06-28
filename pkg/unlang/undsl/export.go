package undsl

// ExportASTNode exports a [*Func] as an [*ASTNode].
func ExportASTNode(f *Func) *ASTNode {
	children := []*ASTNode{}
	for _, entry := range f.Children {
		children = append(children, ExportASTNode(entry))
	}
	return &ASTNode{
		Func:      f.Name,
		Arguments: f.Arguments,
		Children:  children,
	}
}
