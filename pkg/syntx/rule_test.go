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

func TestSimpleRuleWithCat(t *testing.T) {
	var (
		text  = "ab"
		g     = NewGrammar()
		rule1 = NewRule()
		rule2 = NewRule()
	)

	g.Append(
		rule1.Set(
			Char("asd").
				Cat(Der(rule2)),
		),
		rule2.Set(
			Char("bnm"),
		),
	)

	result := g.Run(text)

	fmt.Printf("%s\n", g.Ctx)

	if !result {
		t.Errorf("Failed to parse %s", text)
	}
}

func TestRulesWithNames(t *testing.T) {
	var (
		text = "abc"
		g    = NewGrammar()
		r1   = NewNamedRule("One")
		r2   = NewNamedRule("Two")
		r3   = NewNamedRule("Three")
	)

	g.Append(
		r1.Set(
			Char("asd").
				Cat(Der(r2)),
		),
		r2.Set(
			Char("bnm").
				Cat(Der(r3)),
		),
		r3.Set(
			Char("xcv"),
		),
	)

	result := g.Run(text)

	fmt.Printf("%s\n", g.Ctx)

	if !result {
		t.Errorf("Failed to parse: %s", text)
	}
}
