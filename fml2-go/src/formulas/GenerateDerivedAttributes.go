package formulas

import "domain"

func GenerateDerivedAttributes(data1, data2, data3 float64) (domain.DerivedAttributes, error) {

	return domain.DerivedAttributes{
		FatPercentageCoefficient:   0.21,
		ActiveCalories14DayAverage: 0,
		CurrentFatMass:             25,
		CurrentFatPercentage:       0.14,
		PAL:                        1.20,
		RMR:                        2222,
	}, nil
}
