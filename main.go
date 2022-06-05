package main

import (
	"fmt"
	"strings"
)

type RulerType int

const (
	CharacterType RulerType = iota
	IfSuccessType
)

type Context struct {
	Rules []int
	Literals []string
	RuleJumpTable map[string]int
	FillTable []struct{
		ruleName string
		position int
	}
}

type Ruler interface {
	Build(ctx Context) Context
}

type Character struct {
	charSet string
}

var _ Ruler = (*Character)(nil)

func NewCharacter(charSet string) *Character {
	return &Character {
		charSet: charSet,
	}
}

func (c *Character) Build(ctx Context) Context {
	ctx.Literals = append(ctx.Literals, c.charSet)
	ctx.Rules = append(ctx.Rules, int(CharacterType))
	ctx.Rules = append(ctx.Rules, len(ctx.Literals) - 1)

	return ctx
}

type Concatenation struct {
	leftRuler Ruler
	rightRuler Ruler
}

var _ Ruler = (*Concatenation)(nil)

func NewConcatenation(leftRuler Ruler, rightRuler Ruler) *Concatenation {
	return &Concatenation{
		leftRuler: leftRuler,
		rightRuler: rightRuler,
	}
}

func (c* Concatenation) Build(ctx Context) Context {
	ctx = c.leftRuler.Build(ctx)
	ctx.Rules = append(ctx.Rules, int(IfSuccessType))
	truePos := len(ctx.Rules)
	ctx.Rules = append(ctx.Rules, truePos + 2)
	ctx.Rules = append(ctx.Rules, 0)
	ctx = c.rightRuler.Build(ctx)
	afterRightPos := len(ctx.Rules)
	ctx.Rules[truePos + 1] = afterRightPos

	return ctx
}

type Grammar struct {
	ctx Context
}

func (g *Grammar) Append(rule Ruler) {
	g.ctx = rule.Build(g.ctx)
}

func (g *Grammar) Run(text string) bool {
	var ruleIndex = 0
	var textIndex = 0
	var ok = false

	for ruleIndex < len(g.ctx.Rules) && textIndex < len(text) {
		if RulerType(g.ctx.Rules[ruleIndex]) == CharacterType {
			var (
				literalIndex = g.ctx.Rules[ruleIndex + 1]
				literal = g.ctx.Literals[literalIndex]
			)

			if strings.ContainsAny(string(text[textIndex]), literal) {
				textIndex++
				ruleIndex += 2
				ok = true
			} else {
				ruleIndex += 2
				ok = false
			}
		} else if RulerType(g.ctx.Rules[ruleIndex]) == IfSuccessType {
			if ok {
				ruleIndex = g.ctx.Rules[ruleIndex + 1]
			} else {
				ruleIndex = g.ctx.Rules[ruleIndex + 2]
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

func main() {
	var tests = []func()bool{
		TestCharacter,
		TestConcatenation,
	}

	for n, test := range tests {
		fmt.Printf("************************************\n")
		fmt.Printf("Test: %d\n", n)
		fmt.Printf("\tResult: %t\n", test())
		fmt.Printf("************************************\n")
	}
}
