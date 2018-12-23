package behavior

import (
	"time"
	"enums/ActivityLevel"
	"enums/Experience"
	"enums/DietType"
	"enums/CardioOption"
)

type Habits struct {
	WorkoutsDuration time.Duration
	WorkoutsPerWeek  int64
	ActivityLevel    ActivityLevel.Enum
	GymExperience    Experience.Enum
	AttendsGym       bool
	CookingSkills    Experience.Enum
	DietType         DietType.Enum
	PreferredCardio  []CardioOption.Enum
}