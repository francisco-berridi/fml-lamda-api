package test

import (
	"encoding/json"
	"enums/ActivityLevel"
	"enums/Device"
	"enums/DietType"
	"enums/Experience"
	"enums/Gender"
	"enums/Goal"
	"enums/Unit"
	"fmt"
	"services"
	"testing"
	"domain/shared"
)

func TestOnBoardingUser(t *testing.T) {

	onBoarding, err := services.NewOnBoardingService()
	if err != nil {
		t.Error(err)
	}

	user, err := onBoarding.SignUpUser("francisco.berridi@gmail.com", "Francisco", "Berridi", "321321")
	if err != nil {
		t.Error(err)
	}

	/*// On-boarding 1 (goals & fitness device)
	user.PrimaryGoal = GoalName.WeightLoss
	user.SecondaryGoal = GoalName.Endurance
	user.Device = Device.AppleWatch

	// on-boarding 2 (body measurements)
	user.BiologicalSex = BiologicalSex.Female
	dateOfBirth := time.Date(1981, 12, 11, 0, 0, 0, 0, nil) // UTC
	user.DateOfBirth = &dateOfBirth
	user.SetHeight(161, Unit.Cm)
	user.SetWeight(64, Unit.Kg)
	user.BodyFatPercentage = 0.331

	// on-boarding 3 (behaviour)
	user.GeneralActivityLevel = ActivityLevel.Moderately
	user.UsualExercisesPerWeek = 3
	user.UsualSessionLength = 35
	user.GymExperience = Experience.Intermediate
	user.PreferredExercises = []domain.Exercise{{
		Type: ExerciseType.Cardio, Name: "Cycling",
	},
		{
			Type: ExerciseType.Cardio, Name: "Stairs",
		},
	}

	// on-boarding (medical conditions) TODO

	// on-boarding (cooking skills)
	user.CookingExperience = Experience.Beginner
	user.DietRestriction = DietType.NoRestrictions  //TODO allow multiple restrictions

	user.TargetWeight = 50 //TODO, set units*/

	questionnaire := shared.Questionnaire{}
	questionnaire.PrimaryGoal = Goal.WeightLoss
	questionnaire.SecondaryGoal = Goal.Endurance
	questionnaire.Device = Device.AppleWatch

	questionnaire.BiologicalSex = Gender.Female
	//dateOfBirth := time.Date(1981, 12, 11, 0, 0, 0, 0, nil) // UTC
	//questionnaire.DateOfBirth = &dateOfBirth
	questionnaire.Height = 161
	questionnaire.HeightUnits = Unit.Cm
	questionnaire.Weight = 64
	questionnaire.WeightUnits = Unit.Kg
	questionnaire.BodyFatPercentage = 0.331

	questionnaire.ActivityLevel = ActivityLevel.Moderately
	questionnaire.WorkoutsPerWeek = 3
	questionnaire.WorkoutSessionLength = 20
	questionnaire.GymExperienced = Experience.Intermediate
	questionnaire.AttendsGym = false
	questionnaire.CookingSkills = Experience.Beginner
	questionnaire.DietType = DietType.NoRestrictions

	//workoutPrescription, err:= onBoarding.GenerateWorkoutPrescription(questionnaire)

	//work := workouts.Workout{}

	bytes, _ := json.MarshalIndent(&user, "", "\t")

	fmt.Println(string(bytes))

}
