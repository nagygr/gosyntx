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
	JumpStack     *Stack[int]
	PositionStack *Stack[int]
	FillTable     []RulePosition
}

func NewContext() Context {
	return Context{
		RuleJumpTable: make(map[uint64]int),
		JumpStack:     NewStack[int](),
		PositionStack: NewStack[int](),
	}
}

func (ctx Context) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Literals: %v\nRules:\n\t[\n", ctx.Literals)

	var nextCommandIndex int = 0

	for i, v := range ctx.Rules {
		fmt.Fprintf(&b, "\t\t%3d: %3d", i, v)

		if i == nextCommandIndex && 0 <= v && v < len(CommandNames) {
			fmt.Fprintf(&b, " (%s)", CommandNames[v])
			nextCommandIndex += ArgNums[v] + 1
		}

		fmt.Fprintf(&b, "\n")
	}

	fmt.Fprintf(&b, "\t]\n")

	return b.String()
}
