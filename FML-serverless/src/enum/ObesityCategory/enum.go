package ObesityCategory

type Enum string

const (
	Unknown       Enum = "Unknown"
	UnderWeight   Enum = "UnderWeight"
	NormalWeight  Enum = "NormalWeight"
	OverWeight    Enum = "OverWeight"
	Class1Obesity Enum = "Class1Obesity"
	Class2Obesity Enum = "Class2Obesity"
	Class3Obesity Enum = "Class3Obesity"
)

func GetWithBmi(bmi float64) Enum {

	switch {
	case bmi >= 0 && bmi <= 18.5:
		return UnderWeight

	case bmi > 18.5 && bmi <= 25:
		return NormalWeight

	case bmi > 25 && bmi <= 30:
		return OverWeight

	case bmi > 30 && bmi <= 35:
		return Class1Obesity

	case bmi > 35 && bmi <= 40:
		return Class2Obesity

	default:
		return Class3Obesity
	}
}

func (e Enum) Compare(category Enum) int {
	return e.NumericValue() - category.NumericValue()
}

func (e Enum) NumericValue() int {

	switch e {
	case UnderWeight:
		return 1

	case NormalWeight:
		return 2

	case OverWeight:
		return 3

	case Class1Obesity:
		return 4

	case Class2Obesity:
		return 5

	case Class3Obesity:
		return 6

	default:
		return -1
	}
}
