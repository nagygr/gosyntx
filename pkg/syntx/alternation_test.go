package syntx

import (
	"fmt"
	"testing"
)

func TestAlternation(t *testing.T) {
	var (
		text = "a"
		g    = NewGrammar()
		r1   = NewNamedRule("r1")
		r2   = NewNamedRule("r2")
	)

	g.Append(
		r1.Set(
			Char("bnm").
				Or(Der(r2)),
		),
		r2.Set(
			Char("asd"),
		),
	)

	result := g.Run(text)

	fmt.Printf("%s\n", g.Ctx)

	if !result {
		t.Errorf("Failed to parse %s", text)
	}
}

func TestAlternationOfRules(t *testing.T) {
	var (
		text = "ab"
		g    = NewGrammar()
		r1   = NewNamedRule("r1")
		r2   = NewNamedRule("r2")
		r3   = NewNamedRule("r3")
	)

	g.Append(
		r1.Set(
			Der(r3).
				Or(Der(r2)),
		),
		r2.Set(
			Char("asd").
				Cat(Char("bnm")),
		),
		r3.Set(
			Char("asd").
				Cat(Char("uio")),
		),
	)

	result := g.Run(text)

	fmt.Printf("%s\n", g.Ctx)

	if !result {
		t.Errorf("Failed to parse %s", text)
	}
}
