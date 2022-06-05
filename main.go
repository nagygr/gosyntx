package main

import (
	"fmt"
	"hash/maphash"
	"log"
	"strings"
)

type RulerType int

const (
	CharacterType RulerType = iota
	IfSuccessType
	ReturnType
)

type Context struct {
	Rules         []int
	Literals      []string
	RuleJumpTable map[uint64]int
	FillTable     []struct {
		ruleHash uint64
		position int
	}
}

func (ctx Context) String() string {
	return fmt.Sprintf(
		"Literals: %v\nRules: %v\n",
		ctx.Literals, ctx.Rules,
	)
}

type Ruler interface {
	Build(ctx Context) Context
}

type Character struct {
	charSet string
}

var _ Ruler = (*Character)(nil)

func NewCharacter(charSet string) *Character {
	return &Character{
		charSet: charSet,
	}
}

func (c *Character) Build(ctx Context) Context {
	ctx.Literals = append(ctx.Literals, c.charSet)
	ctx.Rules = append(ctx.Rules, int(CharacterType))
	ctx.Rules = append(ctx.Rules, len(ctx.Literals)-1)

	return ctx
}

type Concatenation struct {
	leftRuler  Ruler
	rightRuler Ruler
}

var _ Ruler = (*Concatenation)(nil)

func NewConcatenation(leftRuler Ruler, rightRuler Ruler) *Concatenation {
	return &Concatenation{
		leftRuler:  leftRuler,
		rightRuler: rightRuler,
	}
}

func (c *Concatenation) Build(ctx Context) Context {
	ctx = c.leftRuler.Build(ctx)
	ctx.Rules = append(ctx.Rules, int(IfSuccessType))
	truePos := len(ctx.Rules)
	ctx.Rules = append(ctx.Rules, truePos+2)
	ctx.Rules = append(ctx.Rules, 0)
	ctx = c.rightRuler.Build(ctx)
	afterRightPos := len(ctx.Rules)
	ctx.Rules[truePos+1] = afterRightPos

	return ctx
}

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
	r.hasher.WriteString(
		fmt.Sprintf("%v", r),
	)
	return r.hasher.Sum64()
}

func (r *Rule) Set(rule Ruler) {
	r.InnerRule = rule
}

func (r *Rule) Build(ctx Context) Context {
	if r == nil {
		log.Panic("Uninitialized rule")
	}

	ruleStartPos := len(ctx.Rules)
	ctx.RuleJumpTable[r.hash()] = ruleStartPos

	ctx = r.InnerRule.Build(ctx)
	ctx.Rules = append(ctx.Rules, int(ReturnType))

	return ctx
}

type Grammar struct {
	Ctx Context
}

func (g *Grammar) Append(rule Ruler) {
	g.Ctx = rule.Build(g.Ctx)
}

func (g *Grammar) Run(text string) bool {
	var ruleIndex = 0
	var textIndex = 0
	var ok = false

	for ruleIndex < len(g.Ctx.Rules) && textIndex < len(text) {
		if RulerType(g.Ctx.Rules[ruleIndex]) == CharacterType {
			var (
				literalIndex = g.Ctx.Rules[ruleIndex+1]
				literal      = g.Ctx.Literals[literalIndex]
			)

			if strings.ContainsAny(string(text[textIndex]), literal) {
				textIndex++
				ruleIndex += 2
				ok = true
			} else {
				ruleIndex += 2
				ok = false
			}
		} else if RulerType(g.Ctx.Rules[ruleIndex]) == IfSuccessType {
			if ok {
				ruleIndex = g.Ctx.Rules[ruleIndex+1]
			} else {
				ruleIndex = g.Ctx.Rules[ruleIndex+2]
			}
		}
	}

	return ok
}

func TestCharacter() bool {
	var text = "b"
	var r1 Ruler = NewCharacter("abc")
	var g Grammar

	g.Append(r1)

	return g.Run(text)
}

func TestConcatenation() bool {
	var text = "ab"
	var g Grammar

	g.Append(
		NewConcatenation(
			NewCharacter("asd"),
			NewCharacter("nbm"),
		),
	)

	return g.Run(text)
}

func TestConcatenationOfConcatenation() bool {
	var text = "abc"
	var g Grammar

	g.Append(
		NewConcatenation(
			NewCharacter("asd"),
			NewConcatenation(
				NewCharacter("nbm"),
				NewCharacter("cvb"),
			),
		),
	)

	result := g.Run(text)

	fmt.Printf("%s\n", g.Ctx)

	return result
}

func main() {
	var tests = []func() bool{
		TestCharacter,
		TestConcatenation,
		TestConcatenationOfConcatenation,
	}

	for n, test := range tests {
		fmt.Printf("************************************\n")
		fmt.Printf("Test: %d\n", n)
		fmt.Printf("\tResult: %t\n", test())
		fmt.Printf("************************************\n")
	}
}
