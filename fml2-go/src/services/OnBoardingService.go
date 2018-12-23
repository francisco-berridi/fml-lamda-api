package services

import (
	"domain"
	"domain/nutrition"
	"domain/shared"
	"domain/workouts"
	"enums/Gender"
	"formulas"
	"services/authentication"
)

type onBoardingService struct {
}

func (o *onBoardingService) SignUpUser(email, firstName, lastName, password string) (domain.User, error) {

	// create user using authentication services
	user, err := authentication.CreateNewUser(email, firstName, lastName, password)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil

}

func (o *onBoardingService) GenerateWorkoutPrescription(questionnaire shared.Questionnaire) (WorkoutPrescription, error) {

	//TODO: this is assumed to be Centimeters
	//dataService := NewDataService()
	//dataService.FindWorkoutSchedules(formulas.CalculateBMI(questionnaire.Height, questionnaire.Weight))

	return WorkoutPrescription{}, nil
}

func (o *onBoardingService) OnBoardUser(questionnaire shared.Questionnaire) (domain.User, error) {

	// NOTE: at this point the user has an non-initialized account, this means that we have already
	// verified the email that is providing as user name and so far is has been marked as "never on-boarded"

	// Generate DerivedAttributes

	// Gather on-boarding data
	// Generate goals
	// validate data
	// Generate user
	// return

	return domain.User{}, nil
}

func (o *onBoardingService) generateUserAttributes(questionnaire shared.Questionnaire,
	onBoardingData domain.OnBoardingData, objectives domain.UserObjectives) (domain.UserAttributes, error) {

	derivedAttributes, err := formulas.GenerateDerivedAttributes(22, 22, 22)

	if err != nil {
		return domain.UserAttributes{}, err
	}

	return domain.UserAttributes{
		OnBoardingData:    onBoardingData,
		DerivedAttributes: derivedAttributes,
		UserObjectives: domain.UserObjectives{
			PrimaryGoal:   "WeightLoss",
			SecondaryGoal: "Aesthetics",
			TargetWeight:  123,
		},
		Age:           questionnaire.Age,
		DateOfBirth:   questionnaire.DateOfBirth,
		BiologicalSex: Gender.Female,
		//HeightCM:       questionnaire.HeightCM,
	}, nil

}

type WorkoutPrescription struct {
	Easy        workouts.WorkoutSchedule
	Recommended workouts.WorkoutSchedule
	Intense     workouts.WorkoutSchedule
}

type DietPrescription struct {
	Easy        nutrition.MealPlan
	Recommended nutrition.MealPlan
	Intense     nutrition.MealPlan
}
