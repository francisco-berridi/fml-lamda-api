package domain

import (
	"domain/workouts"
	"domain/nutrition"
)

type Prescription struct {
	NutritionRx
	ExerciseRx
}

type NutritionRx struct {
	DailyCalories float32
	CarbsPart     float32
	ProteinPart   float32
	FatPart       float32
	MealPlan      nutrition.MealPlan
}


type ExerciseRx struct {
	WeeklyCalories float32
	CardioPart     float32
	StrengthPart   float32
	Schedule       ExerciseSchedule
}

type ExerciseSchedule struct {
	Sunday    []workouts.Workout
	Monday    []workouts.Workout
	Tuesday   []workouts.Workout
	Wednesday []workouts.Workout
	Thursday  []workouts.Workout
	Friday    []workouts.Workout
	Saturday  []workouts.Workout
}

