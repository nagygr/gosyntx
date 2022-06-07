package syntx

type RulerType int

const (
	CharacterType RulerType = iota
	IfSuccessType
	CallType
	ReturnType
)

var (
	CommandNames = []string{
		"Character",
		"IfSuccess",
		"Call",
		"Return",
	}

	ArgNums = []int{
		1,
		2,
		1,
		0,
	}
)
