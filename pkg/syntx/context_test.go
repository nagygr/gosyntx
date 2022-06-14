package syntx

import (
	"fmt"
	"strings"
	"testing"
)

type Data struct {
	CharField string
	IntField  int
	InnerRule *CharStruct
}

type CharStruct struct {
	DataField string
}

func (d *Data) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "CharField: %s\nIntField: %d\nInnerRule:\n", d.CharField, d.IntField)
	fmt.Fprintf(&b, "\tInnerRule.DataField: %s\n", d.InnerRule.DataField)

	return b.String()
}

func TestUnmarshalling(t *testing.T) {
	var (
		text = "a2b"
		g    = NewGrammar()
		r0   = NewRule()
		r1   = NewNamedRule("CharField")
		r2   = NewNamedRule("IntField")
		r3   = NewNamedRule("InnerRule")
		r4   = NewNamedRule("DataField")
	)

	g.Append(
		r0.Set(
			Der(r1).Cat(Der(r2)).Cat(Der(r3)),
		),
		r1.Set(
			Char("asd"),
		),
		r2.Set(
			Char("0123456789"),
		),
		r3.Set(
			Der(r4),
		),
		r4.Set(
			Char("bnm"),
		),
	)

	result := g.Run(text)

	fmt.Printf("%s\n", g.Ctx)

	if !result {
		t.Errorf("Failed to parse %s", text)
	}

	dataIface, err := Unmarshal[Data](g.Ctx.RootNode, text)
	data, ok := dataIface.(Data)

	if err != nil {
		t.Errorf("Error unmarshalling AST: %v", err)
	} else if !ok {
		t.Error("Returned interface is not of type Data")
	} else {
		fmt.Printf("%s\n", &data)
	}
}
