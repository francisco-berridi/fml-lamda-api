package formula

import (
	"math"
)
/** verified */
func RMR(weight, height, genderFactor float64, age int64) float64 {
	return 9.99*weight + 6.25*height - 4.92*float64(age) + genderFactor
}
/** verified */
func BMI(weight, height float64) float64 {
	return weight / math.Pow(height/100, 2)
}
/** verified */
func MifflinStJeor(age int64, weight, height, genderRMRFactor float64) float64 {
	return 9.99*weight + 6.25*height - 4.92*float64(age) + genderRMRFactor
}

// Using the Fitbit equation for Female
func MinuteOfActivityFemale(weight float64, heartRate int64) float64 {
	return -5.92 + 0.0577*float64(heartRate) - 0.0167*weight + 0.00052*float64(heartRate)*weight
}

// Using the Fitbit equation for Male
func MinuteOfActivityMale(weight float64, heartRate int64) float64 {
	return 3.56 - 0.0138*float64(heartRate) - 0.1358*weight + 0.00189*float64(heartRate)*weight
}
