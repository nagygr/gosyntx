package syntx

type Ruler interface {
	Build(ctx Context) Context
}
