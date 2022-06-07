package syntx

import (
	"fmt"
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
	var ok = false
	var endReached = false

	g.resolveFillTable()

	for !endReached && ruleIndex < len(g.Ctx.Rules) {
		fmt.Printf("Executing: %s\n", CommandNames[g.Ctx.Rules[ruleIndex]])
		switch RulerType(g.Ctx.Rules[ruleIndex]) {
		case CharacterType:
			var (
				literalIndex = g.Ctx.Rules[ruleIndex+1]
				literal      = g.Ctx.Literals[literalIndex]
			)

			if ok = textIndex < len(text) && strings.ContainsAny(string(text[textIndex]), literal); ok {
				textIndex++
			}

			ruleIndex += 2

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
				endReached = true
			}

			fmt.Printf("Covered range: ")
			var startPos int = 0
			if pos, ok := g.Ctx.PositionStack.Pop(); ok {
				startPos = pos
			}
			fmt.Printf("[%d:%d]\n", startPos, textIndex)

		case CallType:
			fmt.Printf(
				"Current pos: %d, Jumping to: %d, Return pos: %d\n",
				ruleIndex, g.Ctx.Rules[ruleIndex+1], ruleIndex+2,
			)

			ruleIndex = g.Ctx.Rules[ruleIndex+1]
			g.Ctx.JumpStack.Push(ruleIndex + 2)
			g.Ctx.PositionStack.Push(textIndex)
		}

		fmt.Printf(
			"Status [ok: %t, ruleIndex: %d, textIndex: %d, endReached: %t]\n",
			ok, ruleIndex, textIndex, endReached,
		)
	}

	return ok
}
