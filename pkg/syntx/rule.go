package syntx

import (
	"fmt"
	"hash/maphash"
	"log"
)

type Rule struct {
	InnerRule Ruler
	hasher    maphash.Hash
}

var _ Ruler = (*Rule)(nil)

func NewRule() *Rule {
	return &Rule{
		InnerRule: nil,
	}
}

func (r *Rule) hash() uint64 {
	r.hasher.Reset()
	_, _ = r.hasher.WriteString(
		fmt.Sprintf("%p", r),
	)
	fmt.Printf("Returning hash: %p -> %d\n", r, r.hasher.Sum64())
	return r.hasher.Sum64()
}

func (r *Rule) Set(rule Ruler) *Rule {
	r.InnerRule = rule
	return r
}

func (r *Rule) Build(ctx Context) Context {
	if r == nil {
		log.Panic("Uninitialized rule")
	}

	ruleStartPos := len(ctx.Rules)
	ctx.RuleJumpTable[r.hash()] = ruleStartPos

	ctx = r.InnerRule.Build(ctx)
	ctx.Rules = append(ctx.Rules, int(ReturnType))

	fmt.Printf("Adding rule: %d\n", r.hash())

	return ctx
}

func (r *Rule) Cat(right Ruler) Ruler {
	return NewConcatenation(r, right)
}
