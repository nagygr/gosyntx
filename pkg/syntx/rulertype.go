package syntx

type RulerType int

const (
	CharacterType RulerType = iota
	IfSuccessType
	CallType
	ReturnType
	PushTextPosType
	PopTextPosType
)

var (
	CommandNames = map[RulerType]string{
		CharacterType:   "Character",
		IfSuccessType:   "IfSuccess",
		CallType:        "Call",
		ReturnType:      "Return",
		PushTextPosType: "PushTextPos",
		PopTextPosType:  "PopTextPos",
	}

	ArgNums = map[RulerType]int{
		CharacterType:   1,
		IfSuccessType:   2,
		CallType:        1,
		ReturnType:      0,
		PushTextPosType: 0,
		PopTextPosType:  0,
	}
)
