package analytics

import (
	"formulas/body"
	"enums/Goal"
	"services"
	"enums/Intensity"
	"domain/workouts"
)

func AdviceWorkoutSchedules(userBody *body.BodyComposition, primaryGoal, secondaryGoal Goal.Enum) (map[Intensity.Enum]*workouts.WorkoutSchedule, error) {

	schedules := make(map[Intensity.Enum]*workouts.WorkoutSchedule)

	ds := services.NewDataService()
	schedules, err := ds.FindWorkoutSchedules(userBody.ObesityCategory(), primaryGoal, secondaryGoal)

	if err != nil {
		return schedules, err
	}

	return schedules, nil
}