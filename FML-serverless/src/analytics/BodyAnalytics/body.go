package BodyAnalytics

import (
	"constant"
	"domain"
	"enum/Gender"
	"enum/ObesityCategory"
	"enum/units/Distance"
	"enum/units/Mass"
	"formula"
	"math"
	"time"
)

func NewCalculator(gender Gender.Enum, dateOfBirth time.Time, height Distance.Unit, weigh Mass.Unit) *bodyAnalyticsCalculator {
	return &bodyAnalyticsCalculator{
		height:      height,
		weight:      weigh,
		dateOfBirth: dateOfBirth,
		gender:      gender,
	}
}

type bodyAnalyticsCalculator struct {
	height      Distance.Unit
	weight      Mass.Unit
	dateOfBirth time.Time
	gender      Gender.Enum
}

func (b *bodyAnalyticsCalculator) GenerateStats() domain.BodyStats {
	stats := domain.BodyStats{
		Height:                  b.height,
		Weight:                  b.weight,
		Gender:                  b.gender,
		Age:                     b.Age(),
		RMR:                     b.CalculateRMR(),
		BMI:                     b.CalculateBMI(),
		MaxHeartRate:            b.MaxHeartRate(),
		ObesityCategory:         ObesityCategory.GetWithBmi(b.CalculateBMI()),
		MifflinStJeor:           b.calculateMifflinStJeor(),
		AverageHRDuringWorkouts: b.averageHeartRateDuringWorkouts(),
	}

	return stats
}

func (b *bodyAnalyticsCalculator) CalculateRMR() float64 {
	return formula.RMR(b.weight.Metric(), b.height.Metric(), b.gender.RMRFactor(), b.Age())
}

func (b *bodyAnalyticsCalculator) CalculateBMI() float64 {
	return formula.BMI(b.weight.Metric(), b.height.Metric())
}

func (b *bodyAnalyticsCalculator) Age() int64 {
	return int64(math.Floor(time.Since(b.dateOfBirth).Hours() / constant.OneYear.Hours()))
}

func (b *bodyAnalyticsCalculator) MaxHeartRate() int64 {
	return 220 - b.Age()
}


func (b *bodyAnalyticsCalculator) averageHeartRateDuringWorkouts() int64 {
	return int64(math.Round(0.62 * float64(b.MaxHeartRate())))
}


func (b *bodyAnalyticsCalculator) calculateMifflinStJeor() float64 {
	return formula.MifflinStJeor(b.Age(), b.height.Metric(), b.weight.Metric(), b.gender.RMRFactor())
}
