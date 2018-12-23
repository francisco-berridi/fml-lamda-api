package Distance

import "enum/UnitsSystem"



type Unit struct {
	System UnitsSystem.Enum
	Value  float64
}

func (h *Unit) Metric() float64 {

	if h.System == UnitsSystem.Metric {
		return h.Value
	} else {
		return inchesToCentimeters(h.Value)
	}
}

func (h *Unit) Imperial() float64 {

	if h.System == UnitsSystem.Imperial {
		return h.Value
	} else {
		return centimetersToInches(h.Value)
	}
}

func inchesToCentimeters(in float64) (cm float64) {
	return in * 2.54
}

func centimetersToInches (cm float64) (in float64) {
	return cm / 2.54
}
