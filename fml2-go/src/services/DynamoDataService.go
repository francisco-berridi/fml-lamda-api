package services

import (
	"domain/workouts"
	"enums/Goal"
	"enums/ObesityCategory"
	"services/parsehttp"
	"sort"
	"enums/Intensity"
	"fmt"
)

type DynamoDataService struct {
}

func (*DynamoDataService) FindWorkoutSchedules(obesityCategory ObesityCategory.Enum, primaryGoal, secondaryGoal Goal.Enum) (map[Intensity.Enum]*workouts.WorkoutSchedule, error) {

	schedules := make(map[Intensity.Enum]*workouts.WorkoutSchedule)

	// Parse
	all, err := parsehttp.GetWorkoutSchedules()
	if err != nil {
		return schedules, err
	}

	availableOptions := make([]workouts.WorkoutSchedule, 0)

	for _, value := range all {
		if value.ObesityCategory == obesityCategory &&
			value.PrimaryGoal == primaryGoal &&
			value.SecondaryGoal == secondaryGoal {

			availableOptions = append(availableOptions, value)
		}
	}
	// sort by number of activity days
	sort.Slice(availableOptions, func(i, j int) bool {

		schedule1 := availableOptions[i]
		schedule2 := availableOptions[j]

		return len(schedule1.ActivityDays()) < len(schedule2.ActivityDays())
	})

	fmt.Println("availableOptions", availableOptions)

	if len(availableOptions) == 1 {
		schedules[Intensity.Recommended] = &availableOptions[0]
	} else if len(availableOptions) == 2 {
		schedules[Intensity.Easy] = &availableOptions[0]
		schedules[Intensity.Recommended] = &availableOptions[1]
	} else if len(availableOptions) == 3 {
		schedules[Intensity.Easy] = &availableOptions[0]
		schedules[Intensity.Recommended] = &availableOptions[1]
		schedules[Intensity.Intense] = &availableOptions[2]
	}

	return schedules, nil
}
