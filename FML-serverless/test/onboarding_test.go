package test

import (
	"analytics/ActivityAnalytics"
	"analytics/BodyAnalytics"
	"domain/onboarding"
	"enum/ActivityLevel"
	"enum/CardioOption"
	"enum/DietType"
	"enum/Experience"
	"enum/Gender"
	"enum/Goal"
	"enum/Intensity"
	"enum/ObesityCategory"
	"enum/units"
	"fmt"
	"testing"
	"time"
)

func TestWeightLossOptions(t *testing.T) {

	loc, _ := time.LoadLocation("UTC")

	answers := onboarding.Questionnaire{}

	answers.ObesityCategory = ObesityCategory.OverWeight
	answers.PreferredCardio = []CardioOption.Enum{
		CardioOption.Elliptical, CardioOption.StairMaster, CardioOption.RowingMachine,
	}
	answers.WorkoutSessionLength = 30 * time.Minute
	answers.AttendsGym = true
	answers.GymExperience = Experience.Intermediate
	answers.PrimaryGoal = Goal.WeightLoss
	answers.SecondaryGoal = Goal.GetHealthy
	answers.WorkoutsPerWeek = 4
	answers.ActivityLevel = ActivityLevel.Lightly

	answers.CookingSkills = Experience.Intermediate
	answers.DietType = DietType.NoRestrictions
	answers.TargetWeight = units.Kgs(80)
	answers.BiologicalSex = Gender.Male
	answers.Age = 36
	answers.Weight = units.Kgs(89)
	answers.Height = units.Centimeters(181)

	answers.Location = loc
	answers.DateOfBirth = time.Date(1982, time.August, 26, 0, 0, 0, 0, answers.Location)

	// Body Analytics
	bodyAnalysis := BodyAnalytics.NewCalculator(answers.BiologicalSex, answers.DateOfBirth, answers.Height, answers.Weight)
	bodyStats := bodyAnalysis.GenerateStats()

	fmt.Println(bodyStats.RMR)


	// Activity Analytics
	activityAnalytics := ActivityAnalytics.NewCalculator(bodyStats, answers.ActivityLevel)
	usersPAL := activityAnalytics.CalculatePAL(answers.WorkoutSessionLength, answers.WorkoutsPerWeek)
	fmt.Println(usersPAL)

	recommended, err := activityAnalytics.RecommendWorkoutSchedule(bodyStats.ObesityCategory, answers.PrimaryGoal, answers.SecondaryGoal)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(recommended[Intensity.Easy])
	fmt.Println(recommended[Intensity.Recommended])
	fmt.Println(recommended[Intensity.Intense])

	chosenSchedule := recommended[Intensity.Recommended] // chose an schedule

	activityAnalytics.SelectWorkoutsForSchedule(chosenSchedule, answers.PrimaryGoal, answers.SecondaryGoal, answers.Habits, bodyStats)

	fmt.Println(chosenSchedule.Monday)
	fmt.Println(chosenSchedule.Tuesday)
	fmt.Println(chosenSchedule.Wednesday)
	fmt.Println(chosenSchedule.Tuesday)
	fmt.Println(chosenSchedule.Friday)
	fmt.Println(chosenSchedule.Saturday)
	fmt.Println(chosenSchedule.Sunday)

	// Nutrition Analytics
	//nutritionAnalysis := NutritionAnalytics.NewCalculator()
	//nutritionAnalysis.RecommendMealPlans()

}
