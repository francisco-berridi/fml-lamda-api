package Mass

import "enum/UnitsSystem"

type Unit struct {
	System UnitsSystem.Enum
	Value  float64
}


func (h *Unit) Metric() float64 {

	if h.System == UnitsSystem.Metric {
		return h.Value
	} else {
		return lbToKg(h.Value)
	}
}

func (h *Unit) Imperial() float64 {

	if h.System == UnitsSystem.Imperial {
		return h.Value
	} else {
		return kgToLb(h.Value)
	}
}


func lbToKg(lb float64) float64 {
	return lb * 0.453592
}

func kgToLb(kg float64) float64 {
	return kg / 0.453592
}
