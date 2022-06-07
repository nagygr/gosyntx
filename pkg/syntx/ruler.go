package syntx

type Ruler interface {
	Build(ctx Context) Context
	Cat(right Ruler) Ruler
}
