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
	"fmt"
	"formulas/body"
	"math"
	"services"
	"services/parsehttp"
	"sort"
	"testing"
	"time"
)

func TestSchedulesRoutines(t *testing.T) {
	t.SkipNow()

	ds := services.NewDataService()
	schedules, err := ds.FindWorkoutSchedules(ObesityCategory.NormalWeight, Goal.WeightLoss, Goal.GetHealthy)
	if err != nil {
		t.Error(err)
	}

	for intensity, schedule := range schedules {
		fmt.Println(intensity, schedule.Routine)
	}
}

func TestParseWorkout(t *testing.T) {
	t.SkipNow()

	all, err := parsehttp.GetWorkouts()
	if err != nil {
		t.Error(err)
	}

	answers := shared.Questionnaire{}

	answers.ObesityCategory = ObesityCategory.OverWeight
	answers.PreferredCardio = []CardioOption.Enum{
		CardioOption.Elliptical, CardioOption.StairMaster, CardioOption.RowingMachine,
	}
	answers.WorkoutSessionLength = 30 * time.Minute
	answers.AttendsGym = true
	answers.GymExperienced = Experience.Intermediate
	answers.PrimaryGoal = Goal.WeightLoss
	answers.SecondaryGoal = Goal.GetHealthy

	sort.Slice(all, func(i, j int) bool {
		return all[i].GenerateRank(answers) > all[j].GenerateRank(answers)
	})

	for key, value := range all {
		fmt.Println(key, value.GenerateRank(answers), value.Name, value.TargetExperience, value.Duration, value.BMIRestriction, value.TargetGoals)
	}
}

func TestSchedulesFull(t *testing.T) {

	answers := shared.Questionnaire{}

	answers.ObesityCategory = ObesityCategory.OverWeight
	answers.PreferredCardio = []CardioOption.Enum{
		CardioOption.Elliptical, CardioOption.StairMaster, CardioOption.RowingMachine,
	}
	answers.WorkoutSessionLength = 30 * time.Minute
	answers.AttendsGym = true
	answers.GymExperienced = Experience.Intermediate
	answers.PrimaryGoal = Goal.WeightLoss
	answers.SecondaryGoal = Goal.GetHealthy

	answers.CookingSkills = Experience.Intermediate
	answers.DietType = DietType.NoRestrictions
	answers.TargetWeight = 80
	answers.BiologicalSex = Gender.Male
	answers.Age = 36
	answers.Weight = 89
	answers.Height = 161

	ds := services.NewDataService()
	schedules, err := ds.FindWorkoutSchedules(answers.ObesityCategory, answers.PrimaryGoal, answers.SecondaryGoal)
	if err != nil {
		t.Error(err)
	}

	allWorkouts, err := parsehttp.GetWorkouts()
	if err != nil {
		t.Error(err)
	}

	// pick intense
	schedule := schedules[Intensity.Intense]
	schedule.FillWorkouts(answers, allWorkouts)

	totalCommittedActive := 0.0

	if schedule.Monday != nil {
		resting, active, total, _ := body.CalculateCaloriesExpenedInWorkout(schedule.Monday, answers.BiologicalSex, answers.Age, answers.Weight, answers.Height)
		totalCommittedActive += active
		fmt.Println("resting", resting, "active", active, "total", total)
	}

	if schedule.Tuesday != nil {
		resting, active, total, _ := body.CalculateCaloriesExpenedInWorkout(schedule.Tuesday, answers.BiologicalSex, answers.Age, answers.Weight, answers.Height)
		fmt.Println("resting", resting, "active", active, "total", total)
		totalCommittedActive += active
	}

	if schedule.Wednesday != nil {
		resting, active, total, _ := body.CalculateCaloriesExpenedInWorkout(schedule.Wednesday, answers.BiologicalSex, answers.Age, answers.Weight, answers.Height)
		fmt.Println("resting", resting, "active", active, "total", total)
		totalCommittedActive += active
	}

	if schedule.Thursday != nil {
		resting, active, total, _ := body.CalculateCaloriesExpenedInWorkout(schedule.Thursday, answers.BiologicalSex, answers.Age, answers.Weight, answers.Height)
		fmt.Println("resting", resting, "active", active, "total", total)
		totalCommittedActive += active
	}

	if schedule.Friday != nil {
		resting, active, total, _ := body.CalculateCaloriesExpenedInWorkout(schedule.Friday, answers.BiologicalSex, answers.Age, answers.Weight, answers.Height)
		fmt.Println("resting", resting, "active", active, "total", total)
		totalCommittedActive += active
	}

	if schedule.Saturday != nil {
		resting, active, total, _ := body.CalculateCaloriesExpenedInWorkout(schedule.Saturday, answers.BiologicalSex, answers.Age, answers.Weight, answers.Height)
		fmt.Println("resting", resting, "active", active, "total", total)
		totalCommittedActive += active
	}

	if schedule.Sunday != nil {
		resting, active, total, _ := body.CalculateCaloriesExpenedInWorkout(schedule.Sunday, answers.BiologicalSex, answers.Age, answers.Weight, answers.Height)
		fmt.Println("resting", resting, "active", active, "total", total)
		totalCommittedActive += active
	}

	fmt.Println("=-======== ", math.Round(totalCommittedActive))

	rmr := body.CalculateRMR(answers.BiologicalSex, answers.Age, answers.Weight, answers.Height)

	pal := body.EstimateInitialPAL(body.InitialPALInput{
		Age:                      answers.Age,
		Gender:                   answers.BiologicalSex,
		ActivityLevel:            ActivityLevel.Moderately,
		RMR:                      rmr,
		WeightInKg:               answers.Weight,
		ReportedWorkoutsPerWeek:  answers.WorkoutsPerWeek,
		ReportedWorkoutsDuration: answers.WorkoutSessionLength,
	})

	propPi := body.EstimateCardioPart(ActivityLevel.Moderately, pal, rmr, totalCommittedActive)

	fmt.Println("propPi", propPi)

	// ???
	revicedCML := rmr*pal + totalCommittedActive/2
	fmt.Println("revicedCML", math.Round(revicedCML))






}
