package syntx

import (
	"fmt"
	"log"
	"strings"
)

type Grammar struct {
	Ctx Context
}

func NewGrammar() Grammar {
	return Grammar{
		Ctx: NewContext(),
	}
}

func (g *Grammar) Append(rules ...Ruler) {
	for _, rule := range rules {
		g.Ctx = rule.Build(g.Ctx)
	}
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
	var parseOk = false
	var endReached = false

	g.resolveFillTable()

	name, found := g.Ctx.RuleNameTable[ruleIndex]

	if found {
		if name != "" {
			newNode := NewAstNode(name, Range{0, 0})
			top, ok := g.Ctx.CurrentNode.Top()

			if !ok {
				log.Panic("No AstNode is stack of nodes when dereferencing rule.")
			}

			top.AddChild(newNode)
			g.Ctx.CurrentNode.Push(newNode)
		}
		g.Ctx.CurrentRule.Push(name)
	}

	for !endReached && ruleIndex < len(g.Ctx.Rules) {
		fmt.Printf("Executing: %s @ %d\n", CommandNames[g.Ctx.Rules[ruleIndex]], ruleIndex)
		switch RulerType(g.Ctx.Rules[ruleIndex]) {
		case CharacterType:
			var (
				literalIndex = g.Ctx.Rules[ruleIndex+1]
				literal      = g.Ctx.Literals[literalIndex]
			)

			if parseOk = textIndex < len(text) && strings.ContainsAny(string(text[textIndex]), literal); parseOk {
				textIndex++
			}

			ruleIndex += 2

		case IfSuccessType:
			if parseOk {
				ruleIndex = g.Ctx.Rules[ruleIndex+1]
			} else {
				ruleIndex = g.Ctx.Rules[ruleIndex+2]
			}

		case ReturnType:
			if value, ok := g.Ctx.JumpStack.Pop(); ok {
				ruleIndex = value
				fmt.Printf("Returning to %d\n", ruleIndex)
			} else {
				endReached = true
				fmt.Printf("Returning from parsing (endReached)\n")
			}

			fmt.Printf("Covered range: ")
			var startPos int = 0
			if pos, ok := g.Ctx.PositionStack.Pop(); ok {
				startPos = pos
			}
			fmt.Printf("[%d:%d]\n", startPos, textIndex)

			top, ok := g.Ctx.CurrentRule.Top()

			if !ok {
				log.Panic("No name in stack of RuleNames when returning from rule")
			}

			if top != "" {
				top, ok := g.Ctx.CurrentNode.Top()

				if !ok {
					log.Panic("No AstNode in stack of nodes when returning from rule")
				}

				top.CoveredRange = Range{startPos, textIndex}
				g.Ctx.CurrentNode.Pop()
			}

			g.Ctx.CurrentRule.Pop()

			topRule, ok := g.Ctx.CurrentNode.Top()

			if ok && topRule.Name == "<root>" {
				topRule.CoveredRange = Range{startPos, textIndex}
			}

		case CallType:
			fmt.Printf(
				"Current pos: %d, Jumping to: %d, Return pos: %d\n",
				ruleIndex, g.Ctx.Rules[ruleIndex+1], ruleIndex+2,
			)

			fmt.Printf("Adding address to JumpStack: %d\n", ruleIndex+2)
			g.Ctx.JumpStack.Push(ruleIndex + 2)
			ruleIndex = g.Ctx.Rules[ruleIndex+1]
			g.Ctx.PositionStack.Push(textIndex)

			ruleName := g.Ctx.RuleNameTable[ruleIndex]
			g.Ctx.CurrentRule.Push(ruleName)

			if ruleName != "" {
				newNode := NewAstNode(ruleName, Range{0, 0})
				top, ok := g.Ctx.CurrentNode.Top()

				if !ok {
					log.Panic("No AstNode is stack of nodes when dereferencing rule.")
				}

				fmt.Printf("ADDING TREE NODE: %s to %s\n", ruleName, top.Name)
				top.AddChild(newNode)
				g.Ctx.CurrentNode.Push(newNode)
			}
		}

		fmt.Printf(
			"Status [ok: %t, ruleIndex: %d, textIndex: %d, endReached: %t]\n",
			parseOk, ruleIndex, textIndex, endReached,
		)
	}

	fmt.Printf("Parsed successfully: %t\n", parseOk)
	return parseOk
}
