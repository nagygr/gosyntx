package syntx

import (
	"testing"
)

func TestCharacter(t *testing.T) {
	var text = "b"
	var r1 Ruler = NewCharacter("abc")
	var g = NewGrammar()

	g.Append(r1)

	if !g.Run(text) {
		t.Errorf("Failed to parse %s", text)
	}
}
