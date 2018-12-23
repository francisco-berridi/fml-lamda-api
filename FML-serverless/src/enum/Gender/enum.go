package Gender

type Enum string

const (
	Male   Enum = "Male"
	Female Enum = "Female"
)

func (e Enum) BodyFatFactor1() float64 {

	if e == Male {
		return 37.31
	}

	// Female
	return 39.96
}
func (e Enum) BodyFatFactor2() float64 {

	if e == Male {
		return 103.94
	}

	// Female
	return 102.01
}

func (e Enum) RMRFactor() float64 {

	if e == Male {
		return 5
	}

	// Female
	return -161
}

func (e Enum) Xi_BW() float64 {
	if e == Male {
		return 0.19
	}

	// Female
	return 0.17
}

func (e Enum) CoeffGender() float64 {
	if e == Male {
		return 0.563
	}

	//Female
	return 0.413
}
