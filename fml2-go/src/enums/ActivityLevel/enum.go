package ActivityLevel

type Enum string

const (
	NotActive     Enum = "NotActive"
	VeryLight     Enum = "VeryLight"
	Lightly       Enum = "Lightly"
	Moderately    Enum = "Moderately"
	Active        Enum = "Active"
	HeavilyActive Enum = "HeavilyActive"
)


func (e Enum) NumericValue() float64 {

	switch e {
	case NotActive:
		return 1.2

	case VeryLight:
		return 1.3

	case Lightly:
		return 1.4

	case Moderately:
		return 1.5

	case Active:
		return 1.6

	case HeavilyActive:
		return 1.8
	}

	return -1.0
}
