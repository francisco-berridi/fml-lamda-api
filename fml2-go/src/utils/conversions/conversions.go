package conversions

const CentimetersInFeet = 30.48
const PoundsInAKilogram = 2.20462

func FeetToCentimeters(feet float64) float64 {
	return feet / CentimetersInFeet
}

func LbToKg(pounds float64) float64 {
	return pounds / PoundsInAKilogram
}
