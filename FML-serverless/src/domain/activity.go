package domain

import (
	"enum/CardioOption"
	"enum/DayOfWeek"
	"enum/ExecutionStepType"
	"enum/ExecutionType"
	"enum/ExerciseType"
	"enum/Experience"
	"enum/Goal"
	"enum/ObesityCategory"
	"enum/WorkoutCategory"
	"time"
)

type ActivityStats struct {
	PAL        float64
	CardioPart float64
}

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

func (w *WorkoutSchedule) IsWorkoutUnique(workout *Workout) bool {

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

type Workout struct {
	Id                   string               `json:"id"`
	Name                 string               `json:"name"`
	Summary              string               `json:"summary"`
	RequiresGym          bool                 `json:"requiresGym"`
	MetValue             int64                `json:"metValue"`
	Category             WorkoutCategory.Enum `json:"category"`
	TargetGoals          []Goal.Enum          `json:"targetGoals"`
	Difficulty           int64                `json:"difficulty"`
	AerobicCoefficient   float64              `json:"aerobicCoefficient"`
	Execution            Execution            `json:"execution"`
	TargetExperience     Experience.Enum      `json:"targetExperience"`
	BMIRestriction       ObesityCategory.Enum `json:"bmiRestriction"`
	Duration             time.Duration        `json:"duration"`
	CardioOptions        []CardioOption.Enum  `json:"cardioOptions,omitempty"`
	Intensity            float64              `json:"estimatedHeartRatePercentage"`
	ResistanceProportion float64              `json:"resistanceProportion"`
	CardioProportion     float64              `json:"cardioProportion"`
	Priority             int32                `json:"priority"`
}

func (w *Workout) IsSuitableForGoal(goal Goal.Enum) bool {

	for _, targetGoal := range w.TargetGoals {
		if targetGoal == goal {
			return true
		}
	}

	return false
}

type Execution struct {
	Type       ExecutionType.Enum `json:"type"`
	Steps      []ExecutionStep    `json:"steps"`
	Repetition int32              `json:"repetition"`
}

type ExecutionStep struct {
	Type        ExecutionStepType.Enum `json:"type"`
	Duration    time.Duration          `json:"duration"`
	Exercise    Exercise               `json:"exercise"`
	Series      int                    `json:"series"`
	Repetitions int                    `json:"repetitions"`
	Distance    float32                `json:"distance"`
}

type Exercise struct {
	Type ExerciseType.Enum
	Name string
}
