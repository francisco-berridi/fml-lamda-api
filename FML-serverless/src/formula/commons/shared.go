package commons

import (
	"math"
	"math/rand"
	"fmt"
)

func GetRandomInt(min, max float64) int64 {
	return int64(math.Round(rand.Float64()*(max-min) + min))
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