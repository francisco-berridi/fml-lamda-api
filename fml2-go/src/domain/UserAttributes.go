package domain

import (
	"domain/workouts"
	"enums/ActivityLevel"
	"enums/Device"
	"enums/DietType"
	"enums/Experience"
	"enums/Gender"
	"enums/Goal"
	"time"
)

type UserAttributes struct {
	OnBoardingData    `json:",omitempty"`
	DerivedAttributes `json:",omitempty"`
	UserObjectives    `json:",omitempty"`
	Age               int64       `json:"age,omitempty"`
	DateOfBirth       *time.Time  `json:"dateOfBirth,omitempty"`
	BiologicalSex     Gender.Enum `json:"biologicalSex,omitempty"`
	HeightCM          int64       `json:"heightCm,omitempty"`
	WeightKG          float64     `json:"weightKg,omitempty"`
	Device            Device.Enum `json:"device,omitempty"`
	BodyFatPercentage float64     `json:"bodyFatPercentage,omitempty"`
}

// This are the initial parameters
type OnBoardingData struct {
	GeneralActivityLevel  ActivityLevel.Enum  `json:"generalActivityLevel,omitempty"`
	AllowToExercise       bool                `json:"allowToExercise,omitempty"`
	PriorReportedWeight   float64             `json:"priorReportedWeight,omitempty"`
	UsualExercisesPerWeek int64               `json:"usualExercisesPerWeek,omitempty"`
	UsualSessionLength    int64               `json:"usualSessionLength,omitempty"`
	GymExperience         Experience.Enum     `json:"gymExperience,omitempty"`
	PreferredExercises    []workouts.Exercise `json:"preferredExercises,omitempty"`
	CookingExperience     Experience.Enum     `json:"cookingExperience,omitempty"`
	DietRestriction       DietType.Enum       `json:"dietRestriction,omitempty"`
}

// This are calculated
type DerivedAttributes struct {
	FatPercentageCoefficient   float32 `json:"fatPercentageCoefficient,omitempty"`
	CurrentFatPercentage       float32 `json:"currentFatPercentage,omitempty"`
	RMR                        float32 `json:"rmr,omitempty"`
	CurrentFatMass             float32 `json:"currentFatMass,omitempty"`
	PAL                        float32 `json:"pal,omitempty"`
	ActiveCalories14DayAverage float32 `json:"activeCalories14DayAverage,omitempty"`
}

// This are where the user is expecting to be
type UserObjectives struct {
	PrimaryGoal   Goal.Enum `json:"primaryGoal,omitempty"`
	SecondaryGoal Goal.Enum `json:"secondaryGoal,omitempty"`
	TargetWeight  float64   `json:"targetWeight,omitempty"`
}
