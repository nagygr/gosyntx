package syntx

import (
	"fmt"
	"strings"
)

type RulePosition struct {
	ruleHash uint64
	position int
}

type Context struct {
	Rules         []int
	Literals      []string
	RuleJumpTable map[uint64]int
	RuleNameTable map[int]string
	JumpStack     *Stack[int]
	PositionStack *Stack[int]
	RootNode      *AstNode
	CurrentNode   *Stack[*AstNode]
	CurrentRule   *Stack[string]
	FillTable     []RulePosition
}

func NewContext() Context {
	ctx := Context{
		RuleJumpTable: make(map[uint64]int),
		RuleNameTable: make(map[int]string),
		JumpStack:     NewStack[int](),
		PositionStack: NewStack[int](),
		RootNode:      NewAstNode("<root>", Range{0, 0}),
		CurrentNode:   NewStack[*AstNode](),
		CurrentRule:   NewStack[string](),
	}

	ctx.CurrentNode.Push(ctx.RootNode)

	return ctx
}

func (ctx Context) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Literals: %v\nRule names: %v\nRules:\n\t[\n", ctx.Literals, ctx.RuleNameTable)

	var nextCommandIndex int = 0

	for i, v := range ctx.Rules {
		fmt.Fprintf(&b, "\t\t%3d: %3d", i, v)

		if i == nextCommandIndex && 0 <= v && v < len(CommandNames) {
			fmt.Fprintf(&b, " (%s)", CommandNames[v])
			nextCommandIndex += ArgNums[v] + 1
		}

		fmt.Fprintf(&b, "\n")
	}

	fmt.Fprintf(&b, "\t]\n\n")

	fmt.Fprintf(&b, "AST:\n")

	printAst(&b, 1, ctx.RootNode)

	return b.String()
}

func printAst(b *strings.Builder, tabs int, node *AstNode) {
	printTabs(b, tabs)
	fmt.Fprintf(b, "Name: %s (%p)\n", node.Name, node)

	printTabs(b, tabs)
	fmt.Fprintf(b, "Range: [%d, %d]\n", node.CoveredRange.From, node.CoveredRange.To)

	for _, child := range node.Children {
		printAst(b, tabs+1, child)
	}
}

func printTabs(b *strings.Builder, tabs int) {
	for i := 0; i < tabs; i++ {
		fmt.Fprintf(b, "\t")
	}
}
