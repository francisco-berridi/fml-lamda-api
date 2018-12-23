package units

import (
	"enum/UnitsSystem"
	"enum/units/Distance"
	"enum/units/Mass"
)

func Kgs(value float64) Mass.Unit {
	return Mass.Unit{
		Value:  value,
		System: UnitsSystem.Metric,
	}
}

func Lbs(value float64) Mass.Unit {
	return Mass.Unit{
		Value:  value,
		System: UnitsSystem.Imperial,
	}
}

func Centimeters(value float64) Distance.Unit {
	return Distance.Unit{
		Value:  value,
		System: UnitsSystem.Metric,
	}
}

func Inches(value float64) Distance.Unit {
	return Distance.Unit{
		Value:  value,
		System: UnitsSystem.Imperial,
	}
}
