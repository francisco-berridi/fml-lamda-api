package analytics

import (
	"domain/nutrition"
	"enums/Goal"
	"enums/Intensity"
	"time"
	"formulas/body"
)

type WeightLossScenario struct {
	TimeToGoal           time.Duration
	TimeToReachGoal      time.Time
	WillReachGoal        bool
	ProjectedFinalWeight float64
	Intensity            Intensity.Enum
	RecommendedPlan      nutrition.MealPlan
	ProjectedCurve       []float64
}

type DietOption struct {
	CurrentWeight   float64
	TargetWeight    float64
	PredictionRange time.Duration
	PrimaryGoal     Goal.Enum
	SecondaryGoal   Goal.Enum
	Scenarios       []WeightLossScenario
}


func RecommendMealPlansForWeightLoss(bodyComp *body.BodyComposition, targetWeight float64, activeCaloriesADay float64)([]nutrition.MealPlan, error) {

	plans := make([]nutrition.MealPlan, 0)




	return plans, nil
}

func RecommendMealPlansForANonWeightLoss() {

}