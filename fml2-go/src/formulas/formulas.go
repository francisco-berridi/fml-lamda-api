package formulas

import (
	"enums/Gender"
	"fmt"
	"formulas/body"
	"math"
	"math/rand"
	"utils/inputs"
)

func GeneratePredictionReportInput(gender Gender.Enum, age int64, weightInKg, preGoalWeightKg, heightInCentimeters float64) inputs.PredictionReportInput {

	preGoalWeightArray := []float64{
		preGoalWeightKg, weightInKg,
	}

	// TODO: do not calculate, use input
	preGoalFatPercentage := body.CalculateBodyFatPercentage(gender, age, preGoalWeightKg, heightInCentimeters)
	fatPercentage := body.CalculateBodyFatPercentage(gender, age, weightInKg, heightInCentimeters)

	fatMass := weightInKg * fatPercentage

	preGoalFatMassArray := []float64{
		preGoalWeightKg * preGoalFatPercentage,
		fatMass,
	}

	leanMass := weightInKg - fatMass
	muscleMass := body.CalculateMuscleMass(gender, age, weightInKg, heightInCentimeters, leanMass)

	return inputs.PredictionReportInput{
		PreGoalBWArray:      preGoalWeightArray,
		PreGoalFatMassArray: preGoalFatMassArray,
		FatMass:             fatMass,
		MuscleMass:          muscleMass,
		BodyWeight:          weightInKg,
	}
}

func GetRandomInt(min, max float64) int64 {
	return int64(math.Round(rand.Float64()*(max-min) + min))
}

func VectorTrend(vector []float64) ([]float64, error) {

	if len(vector) < 2 {
		return vector, fmt.Errorf("Vector should have at leat 2 values")
	}

	var numerator = 0.0
	var sumMyVector = 0.0

	vectorLength := float64(len(vector))

	for i := 0.0; i < vectorLength; i++ {
		numerator += (i + 1/2 - vectorLength/2) * vector[int(i)]
		sumMyVector += vector[int(i)]
	}

	var slope = 12 * numerator / ((vectorLength - 1) * vectorLength * (vectorLength + 1))
	var intercept = sumMyVector/vectorLength - slope*(vectorLength+1)/2
	return []float64{intercept + slope, intercept + slope*vectorLength}, nil

}

func Average(values []float64) float64 {

	if len(values) == 0 {
		return 0.0
	}

	if len(values) == 1 {
		return values[0]
	}

	sum := 0.0
	for i := 0; i < len(values); i++ {
		sum += values[i]
	}

	return sum / float64(len(values))

}


type ResultZInterpretation struct {
	BW       float64
	FM       float64
	RMR      float64
	Proport  []float64
	Food1    []float64
	Food2    []float64
	Workouts []float64
	Diet     []float64
	Exercise float64
	LastODEs []float64
	Step     float64
}
