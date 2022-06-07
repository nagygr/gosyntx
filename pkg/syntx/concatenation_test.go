package syntx

import (
	"fmt"
	"testing"
)

func TestConcatenation(t *testing.T) {
	var text = "ab"
	var g = NewGrammar()

	g.Append(
		NewConcatenation(
			NewCharacter("asd"),
			NewCharacter("nbm"),
		),
	)

	result := g.Run(text)

	fmt.Printf("%s\n", g.Ctx)

	if !result {
		t.Errorf("Failed to parse %s", text)
	}
}

func TestConcatenationOfConcatenation(t *testing.T) {
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

	if !result {
		t.Errorf("Failed to parse %s", text)
	}
}
