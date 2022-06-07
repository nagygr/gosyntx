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
	CallType
	ReturnType
	EndType
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
	FillTable     []RulePosition
}

func NewContext() Context {
	return Context{
		RuleJumpTable: make(map[uint64]int),
		JumpStack:     NewStack[int](),
	}
}

func (ctx Context) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Literals: %v\nRules:\n\t[\n", ctx.Literals)

	for i, v := range ctx.Rules {
		fmt.Fprintf(&b, "\t\t%d: %d\n", i, v)
	}

	fmt.Fprintf(&b, "\t]\n")

	return b.String()
}

type Ruler interface {
	Build(ctx Context) Context
}

type End int

var _ Ruler = (*End)(nil)

func NewEnd() *End {
	return new(End)
}

func (e *End) Build(ctx Context) Context {
	ctx.Rules = append(ctx.Rules, int(EndType))
	return ctx
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
		fmt.Sprintf("%p", r),
	)
	fmt.Printf("Returning hash: %p -> %d\n", r, r.hasher.Sum64())
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

	fmt.Printf("Adding rule: %d\n", r.hash())

	return ctx
}

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

type Grammar struct {
	Ctx Context
}

func NewGrammar() Grammar {
	return Grammar{
		Ctx: NewContext(),
	}
}

func (g *Grammar) Append(rule Ruler) {
	g.Ctx = rule.Build(g.Ctx)
}

func (g *Grammar) resolveFillTable() {
	for _, fillElement := range g.Ctx.FillTable {
		fmt.Printf("Looking for Rule: %d\n", fillElement.ruleHash)
		g.Ctx.Rules[fillElement.position] = g.Ctx.RuleJumpTable[fillElement.ruleHash]
		fmt.Printf(
			"Adding values to jump table: [%d]: %d (hash: %d)\n",
			fillElement.position,
			g.Ctx.RuleJumpTable[fillElement.ruleHash],
			fillElement.ruleHash,
		)
	}

	g.Ctx.FillTable = g.Ctx.FillTable[:0]
}

func (g *Grammar) Run(text string) bool {
	var ruleIndex = 0
	var textIndex = 0
	var ok = false
	var endReached = false

	g.resolveFillTable()

	for !endReached && ruleIndex < len(g.Ctx.Rules) && textIndex < len(text) {
		switch RulerType(g.Ctx.Rules[ruleIndex]) {
		case CharacterType:
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

		case IfSuccessType:
			if ok {
				ruleIndex = g.Ctx.Rules[ruleIndex+1]
			} else {
				ruleIndex = g.Ctx.Rules[ruleIndex+2]
			}

		case ReturnType:
			if value, ok := g.Ctx.JumpStack.Pop(); ok {
				ruleIndex = value
			} else {
				ruleIndex++
			}

		case CallType:
			fmt.Printf(
				"Current pos: %d, Jumping to: %d, Return pos: %d\n",
				ruleIndex, g.Ctx.Rules[ruleIndex+1], ruleIndex+2,
			)

			ruleIndex = g.Ctx.Rules[ruleIndex+1]
			g.Ctx.JumpStack.Push(ruleIndex + 2)

		case EndType:
			endReached = true
		}
	}

	return ok
}

func TestCharacter() bool {
	var text = "b"
	var r1 Ruler = NewCharacter("abc")
	var g = NewGrammar()

	g.Append(r1)

	return g.Run(text)
}

func TestConcatenation() bool {
	var text = "ab"
	var g = NewGrammar()

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
	var g = NewGrammar()

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

func TestSimpleRule() bool {
	var (
		text = "ab"
		g    = NewGrammar()
		rule = NewRule()
	)

	rule.Set(
		NewCharacter("bnm"),
	)

	g.Append(
		NewConcatenation(
			NewCharacter("asd"),
			NewConcatenation(
				NewDeref(rule),
				NewEnd(),
			),
		),
	)

	g.Append(rule)

	result := g.Run(text)

	fmt.Printf("%s\n", g.Ctx)

	return result
}

func main() {
	var tests = []func() bool{
		TestCharacter,
		TestConcatenation,
		TestConcatenationOfConcatenation,
		TestSimpleRule,
	}

	for n, test := range tests {
		fmt.Printf("************************************\n")
		fmt.Printf("Test: %d\n", n)
		fmt.Printf("\tResult: %t\n", test())
		fmt.Printf("************************************\n")
	}
}
