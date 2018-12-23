package onboarding

import (
	"enum/ActivityLevel"
	"enum/CardioOption"
	"enum/Device"
	"enum/DietType"
	"enum/Experience"
	"enum/Gender"
	"enum/Goal"
	"enum/ObesityCategory"
	"enum/UnitsSystem"
	"enum/units/Distance"
	"enum/units/Mass"
	"time"
)

type Questionnaire struct {
	aboutYou
	configuration
	healthStatus
	Habits
	objectives
}

type aboutYou struct {
	BiologicalSex     Gender.Enum
	DateOfBirth       time.Time
	Age               int64
	Height            Distance.Unit
	Weight            Mass.Unit
	PreviousWeight    Mass.Unit
	BodyFatPercentage float64
	ObesityCategory   ObesityCategory.Enum
}

type objectives struct {
	TargetWeight  Mass.Unit
	PrimaryGoal   Goal.Enum
	SecondaryGoal Goal.Enum
}

type Habits struct {
	ActivityLevel        ActivityLevel.Enum
	WorkoutsPerWeek      int64
	WorkoutSessionLength time.Duration
	GymExperience        Experience.Enum
	AttendsGym           bool
	CookingSkills        Experience.Enum
	DietType             DietType.Enum
	PreferredCardio      []CardioOption.Enum
}

type healthStatus struct {
	HeartCondition                bool
	ChestPainWhileActive          bool
	ChestPainWhileResting         bool
	GetsDizzy                     bool
	BloodMedication               bool
	AnyOtherReasonToAvoidActivity bool
	DoctorConsent                 bool
}

type configuration struct {
	HeightUnits UnitsSystem.Enum
	WeightUnits UnitsSystem.Enum
	Device      Device.Enum
	Location    *time.Location
}
