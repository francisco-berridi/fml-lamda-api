package activity

import (
	"formulas/body"
	"domain/workouts"
	"math"
	"time"
	"enums/WorkoutCategory"
)

func  EstimateWeeklyCaloriesBurn(userBody *body.BodyComposition, w *workouts.WorkoutSchedule) (resting, active, total float64, err error) {

	if w.Sunday != nil {
		dayResting, dayActive, dayTotal, err := CalculateCaloriesExpendInWorkout(userBody, w.Sunday)

		if err != nil {
			return
		}

		resting += dayResting
		active += dayActive
		total += dayTotal
	}

	if w.Monday != nil {
		dayResting, dayActive, dayTotal, err := CalculateCaloriesExpendInWorkout(userBody, w.Monday)

		if err != nil {
			return
		}

		resting += dayResting
		active += dayActive
		total += dayTotal
	}

	if w.Tuesday != nil {
		dayResting, dayActive, dayTotal, err := CalculateCaloriesExpendInWorkout(userBody, w.Tuesday)

		if err != nil {
			return
		}

		resting += dayResting
		active += dayActive
		total += dayTotal
	}

	if w.Wednesday != nil {
		dayResting, dayActive, dayTotal, err := CalculateCaloriesExpendInWorkout(userBody, w.Wednesday)

		if err != nil {
			return
		}

		resting += dayResting
		active += dayActive
		total += dayTotal
	}

	if w.Thursday != nil {
		dayResting, dayActive, dayTotal, err := CalculateCaloriesExpendInWorkout(userBody, w.Thursday)

		if err != nil {
			return
		}

		resting += dayResting
		active += dayActive
		total += dayTotal
	}

	if w.Friday != nil {
		dayResting, dayActive, dayTotal, err := CalculateCaloriesExpendInWorkout(userBody, w.Friday)

		if err != nil {
			return
		}

		resting += dayResting
		active += dayActive
		total += dayTotal
	}

	if w.Saturday != nil {
		dayResting, dayActive, dayTotal, err := CalculateCaloriesExpendInWorkout(userBody, w.Saturday)

		if err != nil {
			return
		}

		resting += dayResting
		active += dayActive
		total += dayTotal
	}

	return

}


func CalculateCaloriesExpendInWorkout(body *body.BodyComposition, workout *workouts.Workout) (resting, active, total float64, err error) {

	mifflinStJeor := body.MifflinStJeor() // Resting calories in a day

	// The intensity of the workout is based on the users
	// ability to perform this specific workout, estimating
	// what the heart rate of the user will perform this
	// workout
	estimatedIntensity := int64(math.Round(workout.Intensity * float64(body.MaxHeartRate())))

	// Calories per minute the user
	// burns
	burnedCalsForMinute := mifflinStJeor / time.Duration(time.Hour*24).Minutes()

	// RMR per workout duration
	resting = burnedCalsForMinute * workout.Duration.Minutes()

	if estimatedIntensity > 52 && workout.Category.MajorCategory() != WorkoutCategory.Strength {
		total = body.CaloriesBurnDuringActivity(workout.Duration)
	} else {
		// use the MET value
		total = float64(float64(workout.MetValue) * resting)
	}

	active = total - resting

	return
}


/** verified */
func NoWorkoutDayActiveCalories(body *body.BodyComposition, pal float64) float64 {
	return body.MifflinStJeor() * (0.88*pal - 1.0)

}