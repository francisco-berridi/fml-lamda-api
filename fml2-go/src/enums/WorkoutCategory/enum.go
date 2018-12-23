package WorkoutCategory

type Enum string

const (
	Rest             Enum = "Rest"
	Cardio           Enum = "Cardio"
	Strength         Enum = "Strength"
	Circuit          Enum = "Cicuit"
	FullbodyStrength Enum = "FullbodyStrength"
	FullbodyCircuit  Enum = "FullbodyCircuit"
	LowerStrength    Enum = "LowerStrength"
	UpperStrength    Enum = "UpperStrength"
	PullStrength     Enum = "PullStrength"
	PushStrength     Enum = "PushStrength"
	UpperCircuit     Enum = "UpperCircuit"
	LowerCircuit     Enum = "LowerCircuit"
	Other            Enum = "Other"
	OtherStrength    Enum = "OtherStrength"
)

func (e Enum) MajorCategory() Enum {

	switch e {
	case Cardio:
		fallthrough
	case Other:
		return Cardio

	case FullbodyCircuit:
		fallthrough
	case Circuit:
		fallthrough
	case UpperCircuit:
		fallthrough
	case LowerCircuit:
		return Circuit

	case FullbodyStrength:
		fallthrough
	case LowerStrength:
		fallthrough
	case UpperStrength:
		fallthrough
	case PullStrength:
		fallthrough
	case PushStrength:
		fallthrough
	case OtherStrength:
		return Strength

	default:
		return Rest
	}

	return Other
}
