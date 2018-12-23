package workouts

import (
	"enums/DayOfWeek"
	"enums/Goal"
	"enums/ObesityCategory"
	"enums/WorkoutCategory"
	"sort"
	"domain/behavior"
)

type WorkoutSchedule struct {
	Routine         map[DayOfWeek.Enum]WorkoutCategory.Enum `json:"routine"`
	Monday          *Workout                                `json:"monday"`
	Tuesday         *Workout                                `json:"tuesday"`
	Wednesday       *Workout                                `json:"wednesday"`
	Thursday        *Workout                                `json:"thursday"`
	Friday          *Workout                                `json:"friday"`
	Saturday        *Workout                                `json:"saturday"`
	Sunday          *Workout                                `json:"sunday"`
	PrimaryGoal     Goal.Enum                               `json:"primaryGoal"`
	SecondaryGoal   Goal.Enum                               `json:"secondaryGoal"`
	ObesityCategory ObesityCategory.Enum                    `json:"obesityCategory"`
}

func (w *WorkoutSchedule) ActivityDays() []*Workout {

	days := make([]*Workout, 0)
	for day, category := range w.Routine {
		if day == DayOfWeek.Monday && category != WorkoutCategory.Rest {
			days = append(days, w.Monday)
		}
		if day == DayOfWeek.Tuesday && category != WorkoutCategory.Rest {
			days = append(days, w.Tuesday)
		}
		if day == DayOfWeek.Wednesday && category != WorkoutCategory.Rest {
			days = append(days, w.Wednesday)
		}
		if day == DayOfWeek.Thursday && category != WorkoutCategory.Rest {
			days = append(days, w.Thursday)
		}
		if day == DayOfWeek.Friday && category != WorkoutCategory.Rest {
			days = append(days, w.Friday)
		}
		if day == DayOfWeek.Saturday && category != WorkoutCategory.Rest {
			days = append(days, w.Saturday)
		}
		if day == DayOfWeek.Sunday && category != WorkoutCategory.Rest {
			days = append(days, w.Sunday)
		}
	}

	return days
}

func (w *WorkoutSchedule) FillWorkouts(userHabits behavior.Habits, primaryGoal, secondaryGoal Goal.Enum, obesityCategory ObesityCategory.Enum, allWorkouts []Workout) error {

	//Sort workouts by rank
	// TODO: Implement sort
	sort.Slice(allWorkouts, func(i, j int) bool {
		return allWorkouts[i].GenerateRank(userHabits, primaryGoal, secondaryGoal, obesityCategory) > allWorkouts[j].GenerateRank(userHabits, primaryGoal, secondaryGoal, obesityCategory)
	})

	for day, category := range w.Routine {

		if category == WorkoutCategory.Rest {
			continue
		}
		if day == DayOfWeek.Monday {
			w.Monday = w.findWorkoutForCategory(category, allWorkouts)
		} else if day == DayOfWeek.Tuesday {
			w.Tuesday = w.findWorkoutForCategory(category, allWorkouts)
		} else if day == DayOfWeek.Wednesday {
			w.Wednesday = w.findWorkoutForCategory(category, allWorkouts)
		} else if day == DayOfWeek.Thursday {
			w.Thursday = w.findWorkoutForCategory(category, allWorkouts)
		} else if day == DayOfWeek.Friday {
			w.Friday = w.findWorkoutForCategory(category, allWorkouts)
		} else if day == DayOfWeek.Saturday {
			w.Saturday = w.findWorkoutForCategory(category, allWorkouts)
		} else if day == DayOfWeek.Sunday {
			w.Sunday = w.findWorkoutForCategory(category, allWorkouts)
		}
	}

	return nil
}

func (w *WorkoutSchedule) findWorkoutForCategory(category WorkoutCategory.Enum, rankedWorkouts []Workout) *Workout {

	var found *Workout

	for _, ranked := range rankedWorkouts {
		// We give preference to workouts that are unique∂
		// thinking that the user would prefer not to repeat
		// workouts during the week
		if ranked.Category == category && w.isWorkoutUnique(ranked) {
			found = &ranked
			break
		}
	}

	if found == nil {
		// Cannot find a workout that is in the category
		// AND also unique, so repeat the workout
		for _, ranked := range rankedWorkouts {
			if ranked.Category == category {
				found = &ranked
				break
			}
		}
	}

	return found
}

func (w *WorkoutSchedule) isWorkoutUnique(workout Workout) bool {

	if w.Monday != nil && w.Monday.Id == workout.Id {
		return false
	} else if w.Tuesday != nil && w.Tuesday.Id == workout.Id {
		return false
	} else if w.Wednesday != nil && w.Wednesday.Id == workout.Id {
		return false
	} else if w.Thursday != nil && w.Thursday.Id == workout.Id {
		return false
	} else if w.Friday != nil && w.Friday.Id == workout.Id {
		return false
	} else if w.Saturday != nil && w.Saturday.Id == workout.Id {
		return false
	} else if w.Sunday != nil && w.Sunday.Id == workout.Id {
		return false
	}

	return true
}