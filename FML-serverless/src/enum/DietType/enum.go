package DietType

type Enum string

const (
	NoRestrictions Enum = "NoRestrictions"
	Vegetarian     Enum = "Vegetarian"
	Vegan          Enum = "Vegan"
	Paleo          Enum = "Paleo"
	Keto           Enum = "Keto"
)
