package undsl

// ExportASTNode exports a [*Func] as an [*ASTNode].
func ExportASTNode(f *Func) *ASTNode {
	children := []*ASTNode{}
	for _, entry := range f.Children {
		children = append(children, ExportASTNode(entry))
	}
	arguments := f.Arguments
	if arguments == nil {
		arguments = &Empty{}
	}
	return &ASTNode{
		Func:      f.Name,
		Arguments: arguments,
		Children:  children,
	}
}
