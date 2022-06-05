package main

import (
	"fmt"
	"strings"
)

type RuleType int

const (
	CharacterType RuleType = iota
	IfSuccessType
)

type Rule interface {
	Build(rules []int, literals []string) ([]int, []string)
}

type Character struct {
	charSet string
}

var _ Rule = (*Character)(nil)

func NewCharacter(charSet string) *Character {
	return &Character {
		charSet: charSet,
	}
}

func (c *Character) Build(rules []int, literals []string) ([]int, []string) {
	literals = append(literals, c.charSet)
	rules = append(rules, int(CharacterType))
	rules = append(rules, len(literals) - 1)

	return rules, literals
}

type Concatenation struct {
	leftRule Rule
	rightRule Rule
}

var _ Rule = (*Concatenation)(nil)

func NewConcatenation(leftRule Rule, rightRule Rule) *Concatenation {
	return &Concatenation{
		leftRule: leftRule,
		rightRule: rightRule,
	}
}

func (c* Concatenation) Build(rules []int, literals []string) ([]int, []string) {
	rules, literals = c.leftRule.Build(rules, literals)
	rules = append(rules, int(IfSuccessType))
	truePos := len(rules)
	rules = append(rules, truePos + 2)
	rules = append(rules, 0)
	rules, literals = c.rightRule.Build(rules, literals)
	afterRightPos := len(rules)
	rules[truePos + 1] = afterRightPos

	return rules, literals
}

type Grammar struct {
	rules []int
	literals []string
}

func (g *Grammar) Append(rule Rule) {
	g.rules, g.literals = rule.Build(g.rules, g.literals)
}

func (g *Grammar) Run(text string) bool {
	var ruleIndex = 0
	var textIndex = 0
	var ok = false

	for ruleIndex < len(g.rules) && textIndex < len(text) {
		if RuleType(g.rules[ruleIndex]) == CharacterType {
			var (
				literalIndex = g.rules[ruleIndex + 1]
				literal = g.literals[literalIndex]
			)

			if strings.ContainsAny(string(text[textIndex]), literal) {
				textIndex++
				ruleIndex += 2
				ok = true
			} else {
				ruleIndex += 2
				ok = false
			}
		} else if RuleType(g.rules[ruleIndex]) == IfSuccessType {
			if ok {
				ruleIndex = g.rules[ruleIndex + 1]
			} else {
				ruleIndex = g.rules[ruleIndex + 2]
			}
		}
	}

	return ok
}

func TestCharacter() bool {
	var text = "f"
	var r1 Rule = NewCharacter("abc")
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
