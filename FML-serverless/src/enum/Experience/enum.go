package Experience

type Enum string

const (
	Unknown      Enum = "Unknown"
	None         Enum = "None"
	Beginner     Enum = "Beginner"
	Intermediate Enum = "Intermediate"
	Advanced     Enum = "Advanced"
)

func (e Enum) Compare(experience Enum) int64 {
	return e.Value() - experience.Value()
}

func (e Enum) Value() int64 {
	switch e {
	case None:
		return 0

	case Beginner:
		return 1

	case Intermediate:
		return 2

	case Advanced:
		return 3
	}

	return -1
}
