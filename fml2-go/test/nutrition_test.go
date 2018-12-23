package test

import (
	"domain/shared"
	"enums/ActivityLevel"
	"enums/CardioOption"
	"enums/DietType"
	"enums/Experience"
	"enums/Gender"
	"enums/Goal"
	"enums/Intensity"
	"enums/ObesityCategory"
	"enums/Unit"
	"fmt"
	"formulas/body"
	"testing"
	"time"
	"formulas/analytics"
	"services/parsehttp"
	"domain/behavior"
	"formulas/analytics/activity"
	"formulas"
	"encoding/json"
	"domain/nutrition"
	"sort"
	"math"
)

func TestWeightLossDietOptions(t *testing.T) {

	answers := shared.Questionnaire{}

	answers.PrimaryGoal = Goal.WeightLoss
	answers.SecondaryGoal = Goal.GetHealthy

	answers.ObesityCategory = ObesityCategory.OverWeight
	answers.PreferredCardio = []CardioOption.Enum{
		CardioOption.Elliptical, CardioOption.StairMaster, CardioOption.RowingMachine,
	}
	answers.WorkoutsPerWeek = 4
	answers.WorkoutSessionLength = 30 * time.Minute
	answers.AttendsGym = true
	answers.GymExperienced = Experience.Intermediate

	answers.CookingSkills = Experience.Intermediate
	answers.DietType = DietType.NoRestrictions
	answers.Height = 181
	answers.HeightUnits = Unit.Cm

	answers.PreviousWeight = 89.0
	answers.Weight = 75.0
	answers.WeightUnits = Unit.Kg
	answers.BiologicalSex = Gender.Male
	answers.Age = 36
	answers.BodyFatPercentage = 0.30
	answers.ActivityLevel = ActivityLevel.VeryLight

	// TARGET WEIGHT
	answers.TargetWeight = 80

	location, _ := time.LoadLocation("UTC")

	userBody := &body.BodyComposition{
		DateOfBirth: time.Date(1982, time.August, 26, 0,0,0,0, location),
		Weight: 90,
		Height: 181,
		Gender: Gender.Male,
	}

	userHabits := behavior.Habits{
		WorkoutsPerWeek: 4,
		WorkoutsDuration: 30 * time.Minute,
		ActivityLevel: ActivityLevel.Lightly,
	}

	//pal := userBody.PAL(userHabits)
	rmr := userBody.RMR()
	//bmi := userBody.BMI()

	fmt.Println(" ------------------------------------ rmr",rmr)
	fmt.Println(" ------------------------------------ pal", userBody.PAL(userHabits))


	// Select Workout Routine
	// ------------------------
	/*ds := services.NewDataService()
	schedules, err := ds.FindWorkoutSchedules(ObesityCategory.GetWithBmi(userBody.BMI()), answers.PrimaryGoal, answers.SecondaryGoal)
	if err != nil {
		t.Error(err)
	}
*/
	/**/

	schedules, err := analytics.AdviceWorkoutSchedules(userBody, answers.PrimaryGoal, answers.SecondaryGoal)
	if err != nil {
		t.Error(err)
	}

	// pick intense
	schedule := schedules[Intensity.Recommended]


	allWorkouts, err := parsehttp.GetWorkouts()

	schedule.FillWorkouts(userHabits, answers.PrimaryGoal, answers.SecondaryGoal, userBody.ObesityCategory(),  allWorkouts) // ranks
	// --------------------------------------------------------

	fmt.Println(schedule.Routine)

	resting, active, total, err := activity.EstimateWeeklyCaloriesBurn(userBody, schedule)

	if err != nil {
		t.Error(err)
	}

	//fmt.Println("resting", resting, "active", active, "total", total)


	all := 0.0


	//-----------------------------------

	activeCalories14DaysAverage := all / 7.0

	fmt.Println(" ------------------------------------ subscribedExersice", subscribedExersice)
	subscribedExersice = subscribedExersice / 7.0

	prop_pi := body.EstimateCardioPart(answers.ActivityLevel, pal, rmr, activeCalories14DaysAverage)

	fmt.Println(" ------------------------------------ subscribedExersice per day", subscribedExersice)
	fmt.Println(" ------------------------------------ WeeklyActivity", all)
	fmt.Println(" ------------------------------------ prop_pi", prop_pi)

	step := 0.2

	//bfp := body.CalculateBodyFatPercentage(answers.BiologicalSex, answers.Age, answers.Weight, answers.Height)
	fatMass := answers.Weight * answers.BodyFatPercentage

	//FIXME , why days 13 and prediction 100 ?
	resultY, err := formulas.GetDietHall2015(answers.Age, answers.BiologicalSex, answers.Height, pal, prop_pi, rmr,
		answers.PreviousWeight, fatMass, answers.Weight, 0.45, 0.22, 0.33, 13, 0, step)

	_, err = json.MarshalIndent(resultY, "", "\t")
	if err != nil {
		t.Error(err)
	}

	dietPlans, err := parsehttp.GetMealPlans(DietType.NoRestrictions)
	if err != nil {
		t.Error(err)
	}

	//// This is only a duplication of the code
	// subscribedExersice = average active PER DAY (total active in the week / 7)
	// newExercise = Extra exercise that we need to recommend
	newExercise := subscribedExersice - activeCalories14DaysAverage*(1-prop_pi)
	//newExercise 82.54267484952385
	fmt.Println("========================newExercise", newExercise)

	// Calories Mantainance Level
	revisedCML := rmr*pal + newExercise/2
	usablePlans := make([]*nutrition.MealPlan, 0)
	for _, plan := range dietPlans {

		if plan.Calories-revisedCML < 0 {
			usablePlans = append(usablePlans, plan)
		}
	}

	sort.Slice(usablePlans, func(i, j int) bool {
		return usablePlans[i].Calories-usablePlans[j].Calories > 0
	})

	// Weight loss Projection
	// ----------------------------------------------------------
	var allProjectedCurves [][]float64
	var paramZArray []formulas.ResultZInterpretation

	for _, plan := range usablePlans {

		var bodyWeightFMLKg []float64
		if resultY.AverageEnergyIntake < plan.Calories {
			//TODO: complete?
		} else {

			//This is On-Boarding (?)
			//TODO: review
			resultZ := formulas.Rk4ContinueTrajectoryHall2015(plan.CaloriesFromCarbs, plan.CaloriesFromFat, plan.CaloriesFromProtein, newExercise,
				resultY.LastValuesODEFittedModel.G, resultY.LastValuesODEFittedModel.Psig, resultY.LastValuesODEFittedModel.P, resultY.LastValuesODEFittedModel.L_diet,
				resultY.LastValuesODEFittedModel.Fa, resultY.LastValuesODEFittedModel.ECF_F, resultY.LastValuesODEFittedModel.ECF_S, resultY.LastValuesODEFittedModel.Therm, 0.4, 0, 13,
				113, step, resultY.UserDependparamaters, answers.BiologicalSex)

			paramZArray = append(paramZArray, formulas.ResultZInterpretation{
				Diet:     []float64{plan.CaloriesFromCarbs, plan.CaloriesFromFat, plan.CaloriesFromProtein},
				Exercise: newExercise,
				LastODEs: []float64{resultY.LastValuesODEFittedModel.G, resultY.LastValuesODEFittedModel.Psig,
					resultY.LastValuesODEFittedModel.P, resultY.LastValuesODEFittedModel.L_diet, resultY.LastValuesODEFittedModel.Fa,
					resultY.LastValuesODEFittedModel.ECF_F, resultY.LastValuesODEFittedModel.ECF_S, resultY.LastValuesODEFittedModel.Therm},
				Step: step,
			})

			scaleBWZ := answers.Weight - resultZ.BW[0]/1000.0

			for _, value := range resultZ.BW {
				bodyWeightFMLKg = append(bodyWeightFMLKg, value/1000.0+scaleBWZ)
			}
		}

		allProjectedCurves = append(allProjectedCurves, bodyWeightFMLKg)
		plan.WeightLossProjection = bodyWeightFMLKg
	}

	// Select plan based on BMI
	var selectedDiets []*nutrition.MealPlan
	var weightLossLimits []float64

	for _, plan := range usablePlans {

		switch {
		case bmi < 20:
			weightLossLimits = []float64{0.0, 0.5}
			if plan.ProjectedWeeklyWeightLoss(answers.Weight) <= 0.45 {
				selectedDiets = append(selectedDiets, plan)
			}
			break

		case bmi >= 20 && bmi < 23:
			weightLossLimits = []float64{0.5, 1.0}
			if plan.ProjectedWeeklyWeightLoss(answers.Weight) > 0.45 && plan.ProjectedWeeklyWeightLoss(answers.Weight) < 1.05 {
				selectedDiets = append(selectedDiets, plan)
			}
			break

		case bmi >= 23 && bmi < 25:
			weightLossLimits = []float64{0.5, 1.5}
			if plan.ProjectedWeeklyWeightLoss(answers.Weight) > 0.45 && plan.ProjectedWeeklyWeightLoss(answers.Weight) < 1.55 {
				selectedDiets = append(selectedDiets, plan)
			}
			break

		default:
			weightLossLimits = []float64{1.0, 2.0}
			if plan.ProjectedWeeklyWeightLoss(answers.Weight) > 0.95 && plan.ProjectedWeeklyWeightLoss(answers.Weight) < 2.05 {
				selectedDiets = append(selectedDiets, plan)
			}
		}
	}

	fmt.Println("====================== bmi", bmi, weightLossLimits)

	// Correct the number of selected diets
	// remember that the selected diets will be sorted with the
	// Higher allowed consumed calories first
	// and saved in the "diets to recommend"
	dietsRecommendation := make(map[Intensity.Enum]*nutrition.MealPlan)

	totalDiets := len(selectedDiets)

	switch {
	case totalDiets == 0 || totalDiets == 1:
		dietsRecommendation[Intensity.Recommended] = selectedDiets[0]
		break

	case totalDiets == 2:
		dietsRecommendation[Intensity.Easy] = selectedDiets[0]
		dietsRecommendation[Intensity.Recommended] = selectedDiets[1]
		break

	case totalDiets >= 3:
		// select the first one (top)
		dietsRecommendation[Intensity.Easy] = selectedDiets[0]

		// select the "middle one"
		middleIndex := int(math.Floor(float64(len(selectedDiets) / 2.0)))
		fmt.Println("================= middleIndex", middleIndex, len(selectedDiets))
		dietsRecommendation[Intensity.Recommended] = selectedDiets[middleIndex]

		// select the last one
		dietsRecommendation[Intensity.Intense] = selectedDiets[len(selectedDiets)-1]
		break
	}

	if plan, ok := dietsRecommendation[Intensity.Easy]; ok {
		if plan.CanMeetTarget(answers.TargetWeight) {
			fmt.Println("Easy", plan.DaysToTargetWeight(answers.TargetWeight))
		} else {

			totalLoss := answers.Weight - plan.ProjectedLossFinalWeight()
			howFar := totalLoss / (answers.Weight - answers.TargetWeight)
			fmt.Println(answers.TargetWeight, plan.ProjectedLossFinalWeight())
			fmt.Println(fmt.Sprintf("Easy will take you %f%% of you goal, lose %f KG", math.Round(howFar*100), totalLoss), plan.Calories)
		}
	}

	if plan, ok := dietsRecommendation[Intensity.Recommended]; ok {
		fmt.Println("Recommended", plan.DaysToTargetWeight(answers.TargetWeight), plan.Calories)
	}

	if plan, ok := dietsRecommendation[Intensity.Intense]; ok {
		fmt.Println("Intense", plan.DaysToTargetWeight(answers.TargetWeight), plan.Calories)
	}

	// Chose recommended
	chosenPlan := dietsRecommendation[Intensity.Recommended]

	// Get meals
	_, err = parsehttp.GetMeals(chosenPlan)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}



	fmt.Println("User chose", chosenPlan.Summary)

	// workout schedule
	//schedule

}
