package syntx

import (
	"fmt"
	"testing"
)

func TestSimpleRule(t *testing.T) {
	var (
		text  = "ab"
		g     = NewGrammar()
		rule1 = NewRule()
		rule2 = NewRule()
	)

	g.Append(
		rule1.Set(
			NewConcatenation(
				NewCharacter("asd"),
				NewDeref(rule2),
			),
		),
		rule2.Set(
			NewCharacter("bnm"),
		),
	)

	result := g.Run(text)

	fmt.Printf("%s\n", g.Ctx)

	if !result {
		t.Errorf("Failed to parse %s", text)
	}
}

func TestComplexRule(t *testing.T) {
	var (
		text  = "abc"
		g     = NewGrammar()
		rule1 = NewRule()
		rule2 = NewRule()
	)

	g.Append(
		rule1.Set(
			NewConcatenation(
				NewCharacter("asd"),
				NewConcatenation(
					NewDeref(rule2),
					NewCharacter("xcv"),
				),
			),
		),
		rule2.Set(
			NewCharacter("bnm"),
		),
	)

	result := g.Run(text)

	fmt.Printf("%s\n", g.Ctx)

	if !result {
		t.Errorf("Failed to parse %s", text)
	}
}
