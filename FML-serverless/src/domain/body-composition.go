package domain

import (
	"enum/Gender"
	"enum/ObesityCategory"
	"enum/units/Distance"
	"enum/units/Mass"
)

type BodyStats struct {
	// Calculated
	Age                     int64
	RMR                     float64
	BMI                     float64
	MaxHeartRate            int64
	MifflinStJeor           float64
	AverageHRDuringWorkouts int64
	// Measurements
	Height          Distance.Unit
	Weight          Mass.Unit
	Gender          Gender.Enum
	ObesityCategory ObesityCategory.Enum
}
