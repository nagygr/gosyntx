package syntx

type Range struct {
	From int
	To   int
}

type AstNode struct {
	Name         string
	CoveredRange Range
	Children     []*AstNode
}

func NewAstNode(name string, covered Range) *AstNode {
	return &AstNode{
		Name:         name,
		CoveredRange: covered,
		Children:     make([]*AstNode, 0),
	}
}

func (a *AstNode) AddChild(node *AstNode) {
	a.Children = append(a.Children, node)
}
