package body

import (
	"domain/behavior"
	"enums/ActivityLevel"
	"enums/Gender"
	"enums/ObesityCategory"
	"math"
	"time"
)

func CalculateBodyFatPercentage(gender Gender.Enum, age int64, weightInKg, heightInCentimeters float64) float64 {
	//FIXME Sanity check ?

	return (0.14*float64(age) + 37.31*math.Log(10000*weightInKg/heightInCentimeters/heightInCentimeters) - 103.94) / 100

}

func CalculateMuscleMass(gender Gender.Enum, age int64, weightInKg, heightInCentimeters, leanBodyMass float64) float64 {

	var ecf float64
	var cm float64
	var ecp float64

	if gender == Gender.Male {
		ecf = 0.191*weightInKg + 9.570*heightInCentimeters/100 + 0.025*float64(age) - 12.424
	} else {
		ecf = 0.167*weightInKg + 5.987*heightInCentimeters/100 - 4.027
	}

	cm = leanBodyMass - 1.732*0.04*weightInKg - 1.01087*ecf
	ecp = 0.732*0.04*weightInKg + 0.01087*ecf

	return 0.45*(0.25*cm+ecp) + 0.8*0.5 + 0.35*(0.05*cm-0.5) + 0.413*(0.7*cm+ecf)
}

type InitialPALInput struct {
	ReportedWorkoutsDuration time.Duration
	ReportedWorkoutsPerWeek  int64
	Gender                   Gender.Enum
	Age                      int64
	WeightInKg               float64
	ActivityLevel            ActivityLevel.Enum
	RMR                      float64
}

/** verified */
// Calculates pop_pi
func EstimateCardioPart(activityLevel ActivityLevel.Enum, pal, rmr, averageActiveCalories float64) float64 {
	cardioPart := ((activityLevel.NumericValue() - (pal * 0.12) - 1) * rmr) / averageActiveCalories
	return math.Min(math.Max(cardioPart, 0), 1)

}

// BodyComposition composition, always current or most
// recent values
type BodyComposition struct {
	DateOfBirth time.Time   // constant
	Weight      float64     // KG, TODO Create Height units
	Height      float64     // cms, TODO
	Gender      Gender.Enum // constant
}

// Age is calculated
func (b *BodyComposition) Age() int64 {
	return int64(b.DateOfBirth.Add(time.Since(b.DateOfBirth)).Year())
}

/** verified */
func (b *BodyComposition) MaxHeartRate() int64 {
	return 220 - b.Age()
}

func (b *BodyComposition) AverageHeartRateDuringWorkouts() int64 {
	return int64(math.Round(0.62 * float64(b.MaxHeartRate())))
}

func (b *BodyComposition) RMR() float64 {
	return 9.99*b.Weight + 6.25*b.Height - 4.92*float64(b.Age()) + b.Gender.RMRFactor()
}

func (b *BodyComposition) BMI() float64 {
	return b.Weight / math.Pow(b.Height/100, 2)
}

func (b *BodyComposition) PAL(habits behavior.Habits) float64 {

	weeklyTimeWorkingOut := habits.WorkoutsDuration.Minutes() * float64(habits.WorkoutsPerWeek)
	averageTimeOfActivityPerDay := time.Duration(weeklyTimeWorkingOut / 7.0 * float64(time.Minute.Nanoseconds()))
	averageWeeklyActivityCalories := b.CaloriesBurnDuringActivity(averageTimeOfActivityPerDay)

	return habits.ActivityLevel.NumericValue() + averageWeeklyActivityCalories/b.RMR()
}

func (b *BodyComposition) ObesityCategory() ObesityCategory.Enum {
	return ObesityCategory.GetWithBmi(b.BMI())
}

/** verified */
func (b *BodyComposition) MifflinStJeor() float64 {
	return 9.99*b.Weight + 6.25*b.Height - 4.92*float64(b.Age()) + b.Gender.RMRFactor()
}

//TODO: Maybe move to formulas

func (b *BodyComposition) CaloriesBurnDuringActivity(activityDuration time.Duration) float64 {
	return b.CaloriesBurnPerMinuteOfActivity() * activityDuration.Minutes()
}

// Fitbit equation
func (b *BodyComposition) CaloriesBurnPerMinuteOfActivity() float64 {

	heartRate := b.AverageHeartRateDuringWorkouts()

	if b.Gender == Gender.Female {
		return -5.92 + 0.0577*float64(heartRate) - 0.0167*b.Weight + 0.00052*float64(heartRate)*b.Weight
	} else {
		return 3.56 - 0.0138*float64(heartRate) - 0.1358*b.Weight + 0.00189*float64(heartRate)*b.Weight
	}
}
