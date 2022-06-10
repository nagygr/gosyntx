package syntx

type Alternation struct {
	leftRuler  Ruler
	rightRuler Ruler
}

var _ Ruler = (*Alternation)(nil)

func NewAlternation(leftRuler Ruler, rightRuler Ruler) *Alternation {
	return &Alternation{
		leftRuler:  leftRuler,
		rightRuler: rightRuler,
	}
}

func (c *Alternation) Build(ctx Context) Context {
	ctx.Rules = append(ctx.Rules, int(PushTextPosType))
	ctx = c.leftRuler.Build(ctx)
	ctx.Rules = append(ctx.Rules, int(IfSuccessType))
	truePos := len(ctx.Rules)
	ctx.Rules = append(ctx.Rules, 0)
	ctx.Rules = append(ctx.Rules, truePos+2)
	ctx.Rules = append(ctx.Rules, int(PopTextPosType))
	ctx = c.rightRuler.Build(ctx)
	afterRightPos := len(ctx.Rules)
	ctx.Rules[truePos] = afterRightPos

	return ctx
}

func (c *Alternation) Cat(right Ruler) Ruler {
	return NewConcatenation(c, right)
}

func (c *Alternation) Or(right Ruler) Ruler {
	return NewAlternation(c, right)
}
