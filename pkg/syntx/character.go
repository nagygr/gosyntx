package syntx

type Character struct {
	charSet string
}

var _ Ruler = (*Character)(nil)

func NewCharacter(charSet string) *Character {
	return &Character{
		charSet: charSet,
	}
}

func (c *Character) Build(ctx Context) Context {
	ctx.Literals = append(ctx.Literals, c.charSet)
	ctx.Rules = append(ctx.Rules, int(CharacterType))
	ctx.Rules = append(ctx.Rules, len(ctx.Literals)-1)

	return ctx
}
