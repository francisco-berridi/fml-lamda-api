package workouts

import (
	"enums/CardioOption"
	"enums/ExecutionStepType"
	"enums/ExecutionType"
	"enums/Experience"
	"enums/Goal"
	"enums/ObesityCategory"
	"enums/WorkoutCategory"
	"math"
	"time"
	"domain/behavior"
)

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

	Priority int32 `json:"priority"`
}

// Generates a rank (0.0 - 1.0) based on the following parameters
//
// 40% - for matching the primary and secondary categories (80/20)
// 30% - for desired workout length and level of expertise
// 20% - for match the preferred cardio category
// 10% - for Obesity restrictions
func (w Workout) GenerateRank(userHabits behavior.Habits, primaryGoal, secondaryGoal Goal.Enum, obesityCategory ObesityCategory.Enum) float64 {

	score := 0.0

	goalsFactor := 0.40
	preferencesFactor := 0.20
	workoutsFactor := 0.20
	obesityFactor := 0.20

	/*
	 * Goals
	 */
	goalsScore := 0.0
	if w.suitableForGoal(primaryGoal) {
		// its way more important the primary goal , so it receives 80% of the goal score
		goalsScore += 0.8
	}

	if w.suitableForGoal(secondaryGoal) {
		goalsScore += 0.2
	}

	score += goalsScore * goalsFactor

	/*
	 * Workouts experience
	 */
	lengthScore := 1.0
	// The difference between the workout duration and the requested duration
	// people will prefer shorted workouts so same or less gets the best score (1.0), more
	// than that will decrease the score by .1 by every 3 (180 sec) minutes that its length increases
	// until it gets to 0 to lose the score
	difference := w.Duration - userHabits.WorkoutsDuration

	if difference > 0 {
		// this means that the workout is longer that desired
		lengthScore = 1 - (0.1 * (difference.Minutes() / 3))
		if lengthScore < 0 {
			lengthScore = 0
		}
	}

	requiresGymScore := 0.0
	// If the workout requires a gym but the user doesnt attend one, no score , if attends, full score
	// If the workout doesnt require a gym then full score because home workouts
	// can be done at the gym too
	if w.RequiresGym && userHabits.AttendsGym {
		requiresGymScore = 1
	}
	experienceScore := 1.0

	// the farther the workout required experience and the users experience are
	// the less likely the user will want to do this workout, account for that
	experienceScore = 1 - math.Abs(float64(w.TargetExperience.Compare(userHabits.GymExperience))*0.5)
	if experienceScore < 0 {
		experienceScore = 0
	}

	workoutScore := (lengthScore + requiresGymScore + experienceScore) / 3
	score += workoutScore * workoutsFactor

	//fmt.Println("sum", (lengthScore + requiresGymScore + experienceScore), "workoutScore", workoutScore, "score", score)

	/*
	 * Preferred Cardio Exercises factor
	 */

	cardioScore := 0.0
	if len(w.CardioOptions) != 0 {
		for _, preferred := range userHabits.PreferredCardio {
			for _, option := range w.CardioOptions {
				if option == preferred {
					cardioScore += 0.5
				}
			}
		}
	}
	score += cardioScore * preferencesFactor

	obesityScore := 1.0
	/* obesity factor */
	if obesityCategory.Compare(w.BMIRestriction) > 0 {
		// Negative result means that the user can do the exercise
		// workout receives the full score (1.0), but if he user is not
		// too far from the recommended obesity restriction then the
		// workout will receive some score
		obesityScore = 1 - (0.5 * float64(obesityCategory.Compare(w.BMIRestriction)))
		if obesityScore < 0 {
			obesityScore = 0
		}
	}

	score += obesityScore * obesityFactor

	return score

}

func (w Workout) suitableForGoal(goal Goal.Enum) bool {

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

func Sequence() Execution {
	return Execution{
		Type:  ExecutionType.Sequence,
		Steps: make([]ExecutionStep, 0),
	}
}

func SuperSet() Execution {
	return Execution{
		Type:  ExecutionType.SuperSet,
		Steps: make([]ExecutionStep, 0),
	}
}

func Gym(repetitions, series int) ExecutionStep {
	return ExecutionStep{
		Type:        ExecutionStepType.Gym,
		Exercise:    Exercise{},
		Repetitions: repetitions,
		Series:      series,
	}
}
