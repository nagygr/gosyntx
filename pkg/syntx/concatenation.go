package syntx

type Concatenation struct {
	leftRuler  Ruler
	rightRuler Ruler
}

var _ Ruler = (*Concatenation)(nil)

func NewConcatenation(leftRuler Ruler, rightRuler Ruler) *Concatenation {
	return &Concatenation{
		leftRuler:  leftRuler,
		rightRuler: rightRuler,
	}
}

func (c *Concatenation) Build(ctx Context) Context {
	ctx = c.leftRuler.Build(ctx)
	ctx.Rules = append(ctx.Rules, int(IfSuccessType))
	truePos := len(ctx.Rules)
	ctx.Rules = append(ctx.Rules, truePos+2)
	ctx.Rules = append(ctx.Rules, 0)
	ctx = c.rightRuler.Build(ctx)
	afterRightPos := len(ctx.Rules)
	ctx.Rules[truePos+1] = afterRightPos

	return ctx
}

func (c *Concatenation) Cat(right Ruler) Ruler {
	return NewConcatenation(c, right)
}

func (c *Concatenation) Or(right Ruler) Ruler {
	return NewAlternation(c, right)
}
