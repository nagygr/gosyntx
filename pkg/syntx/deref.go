package syntx

import (
	"fmt"
)

type Deref struct {
	RulePtr *Rule
}

var _ Ruler = (*Deref)(nil)

func NewDeref(ptr *Rule) *Deref {
	return &Deref{ptr}
}

func (d *Deref) Build(ctx Context) Context {
	currentPos := len(ctx.Rules)
	ctx.Rules = append(ctx.Rules, int(CallType))
	ctx.Rules = append(ctx.Rules, 0)
	ctx.FillTable = append(
		ctx.FillTable,
		RulePosition{
			ruleHash: d.RulePtr.hash(),
			position: currentPos + 1,
		},
	)

	fmt.Printf("Appending rule to FillTable: %d\n", d.RulePtr.hash())

	return ctx
}
