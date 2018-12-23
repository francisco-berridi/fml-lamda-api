package shared

import (
	"enums/ActivityLevel"
	"enums/CardioOption"
	"enums/Device"
	"enums/DietType"
	"enums/Experience"
	"enums/Gender"
	"enums/Goal"
	"enums/ObesityCategory"
	"enums/Unit"
	"time"
)

type Questionnaire struct {
	AboutYou
	Configuration
	HealthStatus
	DailyRoutine
	Objectives
}

type AboutYou struct {
	BiologicalSex     Gender.Enum
	DateOfBirth       *time.Time
	Age               int64
	Height            float64
	Weight            float64
	PreviousWeight    float64
	BodyFatPercentage float64
	ObesityCategory   ObesityCategory.Enum
}

type Objectives struct {
	TargetWeight  float64
	PrimaryGoal   Goal.Enum
	SecondaryGoal Goal.Enum
}

type DailyRoutine struct {
	ActivityLevel        ActivityLevel.Enum
	WorkoutsPerWeek      int64
	WorkoutSessionLength time.Duration
	GymExperienced       Experience.Enum
	AttendsGym           bool
	CookingSkills        Experience.Enum
	DietType             DietType.Enum
	PreferredCardio      []CardioOption.Enum
}

type HealthStatus struct {
	HeartCondition                bool
	ChestPainWhileActive          bool
	ChestPainWhileResting         bool
	GetsDizzy                     bool
	BloodMedication               bool
	AnyOtherReasonToAvoidActivity bool
	DoctorConsent                 bool
}

type Configuration struct {
	HeightUnits Unit.Enum
	WeightUnits Unit.Enum
	Device      Device.Enum
	Location    *time.Location
}
