package parsehttp

import (
	"domain"
	"encoding/json"
	"enum/CardioOption"
	"enum/DayOfWeek"
	"enum/Experience"
	"enum/Goal"
	"enum/ObesityCategory"
	"enum/WorkoutCategory"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func GetWorkouts() ([]*domain.Workout, error) {
	client := &http.Client{}
	all := make([]*domain.Workout, 0)

	req, err := http.NewRequest("GET", "https://dev.api.fitmylife.net/parse/classes/Workout?c=999", nil)
	if err != nil {
		return all, err
	}

	req.Header.Add("X-Parse-Application-Id", "fitmylifeAppId")
	req.Header.Add("X-Parse-REST-API-Key", "flrcGF4lP4Lx6M2vq0rLo4MiAM2IN4BHRTr4QQKQ")

	resp, err := client.Do(req)
	if err != nil {
		return all, err
	}

	chars, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return all, err
	}

	result := ParseGeneralResult{}

	json.Unmarshal(chars, &result)

	for _, value := range result.Results {

		raw := value.(map[string]interface{})

		if raw["subcategory"] == nil {
			continue
		}
		summary, ok := raw["summary"].(string)
		if !ok {
			summary = ""
		}

		workout := domain.Workout{
			Id:         raw["objectId"].(string),
			Category:   parseWorkoutCategory(raw["subcategory"].(string)),
			Difficulty: int64(raw["difficulty"].(float64)),
			Name:       raw["name"].(string),
			Summary:    summary,
		}

		if raw["cardioPreferences"] != nil {
			workout.CardioOptions = parseCardioOptions(raw["cardioPreferences"].([]interface{}))
		}

		if raw["goalRelated"] != nil {
			workout.TargetGoals = parseTargetGoals(raw["goalRelated"].([]interface{}))
		}

		if raw["excludeObesityCategory"] != nil {
			workout.BMIRestriction = parseBMIRestriction(raw["excludeObesityCategory"].([]interface{}))
		} else {
			workout.BMIRestriction = ObesityCategory.UnderWeight
		}

		// estimatedDuration is in seconds
		workout.Duration = time.Duration(raw["estimatedDuration"].(float64)) * time.Second

		workout.TargetExperience = parseTargetExperience(raw["difficulty"].(float64))
		workout.Intensity = raw["estimatedHeartRatePercentage"].(float64) / 100

		workout.MetValue = int64(raw["metValue"].(float64))

		workout.ResistanceProportion = (raw["resistanceComponent"].(float64)) / 100
		workout.CardioProportion = 1 - workout.ResistanceProportion

		all = append(all, &workout)
	}

	return all, nil
}

func GetWorkoutSchedules() ([]*domain.WorkoutSchedule, error) {

	client := &http.Client{}
	schedules := make([]*domain.WorkoutSchedule, 0)

	req, err := http.NewRequest("GET", "https://dev.api.fitmylife.net/parse/classes/WorkoutSchedule?limit=999", nil)
	if err != nil {
		return schedules, err
	}

	req.Header.Add("X-Parse-Application-Id", "fitmylifeAppId")
	req.Header.Add("X-Parse-REST-API-Key", "flrcGF4lP4Lx6M2vq0rLo4MiAM2IN4BHRTr4QQKQ")

	resp, err := client.Do(req)
	if err != nil {
		return schedules, err
	}

	chars, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return schedules, err
	}

	result := ParseGeneralResult{}
	json.Unmarshal(chars, &result)

	for _, value := range result.Results {

		raw := value.(map[string]interface{})

		if raw["primaryGoal"] == nil || raw["secondaryGoal"] == nil {
			continue
		}

		primaryGoal := parseGoalWithParseName(raw["primaryGoal"].(string))
		secondaryGoal := parseGoalWithParseName(raw["secondaryGoal"].(string))

		if primaryGoal == Goal.Unknown || secondaryGoal == Goal.Unknown {
			panic(fmt.Sprintf("bad goal %s %s", raw["primaryGoal"], raw["secondaryGoal"]))
		}

		schedule := &domain.WorkoutSchedule{
			PrimaryGoal:     parseGoalWithParseName(raw["primaryGoal"].(string)),
			SecondaryGoal:   parseGoalWithParseName(raw["secondaryGoal"].(string)),
			ObesityCategory: parseObesityCategory(raw["obesityCategory"].(string)),
			Routine:         parseRoutine(raw),
		}

		schedules = append(schedules, schedule)
	}

	return schedules, nil

}

func parseObesityCategory(name string) ObesityCategory.Enum {

	if strings.Contains(name, "underWeight") {
		return ObesityCategory.UnderWeight
	}

	if strings.Contains(name, "normal") {
		return ObesityCategory.NormalWeight
	}

	if strings.Contains(name, "overWeight") {
		return ObesityCategory.OverWeight
	}

	if strings.Contains(name, "lass1Obesity") {
		return ObesityCategory.Class1Obesity
	}

	if strings.Contains(name, "lass2Obesity") {
		return ObesityCategory.Class2Obesity
	}

	if strings.Contains(name, "lass3Obesity") {
		return ObesityCategory.Class3Obesity
	}

	return ObesityCategory.Unknown

}

func parseGoalWithParseName(name string) Goal.Enum {

	if strings.Contains(name, "endurance") {
		return Goal.Endurance
	}

	if strings.Contains(name, "strength") {
		return Goal.Strength
	}

	if strings.Contains(name, "circuit") || strings.Contains(name, "cardioEndurance") {
		return Goal.Circuit
	}

	if strings.Contains(name, "aesthetics") || strings.Contains(name, "athleticStrength") {
		return Goal.Aesthetics
	}

	if strings.Contains(name, "Health") {
		return Goal.GetHealthy
	}

	if strings.Contains(name, "weight") {
		return Goal.WeightLoss
	}

	fmt.Println("======================================= goal name", name)
	return Goal.Unknown
}

func parseWorkoutCategory(name string) WorkoutCategory.Enum {

	if strings.Contains(name, "cardio") {
		return WorkoutCategory.Cardio
	}
	if strings.Contains(name, "fullbodyStrength") {
		return WorkoutCategory.FullbodyStrength
	}

	if strings.Contains(name, "lowerStrength") {
		return WorkoutCategory.LowerStrength
	}
	if strings.Contains(name, "upperStrength") {
		return WorkoutCategory.UpperStrength
	}
	if strings.Contains(name, "pullStrength") {
		return WorkoutCategory.PullStrength
	}
	if strings.Contains(name, "pushStrength") {
		return WorkoutCategory.PushStrength
	}
	if strings.Contains(name, "fullbodyCircuit") {
		return WorkoutCategory.FullbodyCircuit
	}

	if strings.Contains(name, "upperCircuit") {
		return WorkoutCategory.UpperCircuit
	}
	if strings.Contains(name, "lowerCircuit") {
		return WorkoutCategory.LowerCircuit
	}
	if strings.Contains(name, "otherStrength") {
		return WorkoutCategory.OtherStrength
	}
	if strings.Contains(name, "strength") {
		return WorkoutCategory.Strength
	}

	return WorkoutCategory.Other
}

func parseCardioOptions(raw []interface{}) []CardioOption.Enum {

	options := make([]CardioOption.Enum, 0)

	for _, value := range raw {

		name := value.(string)

		if strings.Contains(name, "rowingMachine") {
			options = append(options, CardioOption.RowingMachine)
		}
		if strings.Contains(name, "lliptical") {
			options = append(options, CardioOption.Elliptical)
		}
		if strings.Contains(name, "treadmill") {
			options = append(options, CardioOption.Treadmill)
		}
		if strings.Contains(name, "stairs") {
			options = append(options, CardioOption.Stairs)
		}
		if strings.Contains(name, "running") {
			options = append(options, CardioOption.Running)
		}
		if strings.Contains(name, "walking") {
			options = append(options, CardioOption.Walking)
		}
		if strings.Contains(name, "indoorSpinCycle") {
			options = append(options, CardioOption.IndoorSpinCycle)
		}
		if strings.Contains(name, "swimming") {
			options = append(options, CardioOption.Swimming)
		}
		if strings.Contains(name, "stairMaster") {
			options = append(options, CardioOption.StairMaster)
		}
		if strings.Contains(name, "indoorSpinBike") {
			options = append(options, CardioOption.IndoorSpinBike)
		}
	}

	return options
}

func parseTargetGoals(raw []interface{}) []Goal.Enum {
	targetGoals := make([]Goal.Enum, 0)

	for _, value := range raw {
		targetGoals = append(targetGoals, parseGoalWithParseName(value.(string)))
	}

	return targetGoals
}

func parseBMIRestriction(raw []interface{}) ObesityCategory.Enum {

	all := make([]ObesityCategory.Enum, 0)

	for _, value := range raw {
		all = append(all, parseObesityCategory(value.(string)))
	}

	ranked := ObesityCategory.Class3Obesity

	for _, category := range all {

		if category.Compare(ranked) <= 0 {
			ranked = category
		}
	}

	return ranked

}

func parseTargetExperience(number float64) Experience.Enum {

	switch {
	case number < 10:
		return Experience.None
	case number >= 10 && number <= 12:
		return Experience.Beginner
	case number >= 13 && number <= 16:
		return Experience.Intermediate
	case number >= 17:
		return Experience.Advanced
	}

	return Experience.Unknown
}

func parseRoutine(raw map[string]interface{}) map[DayOfWeek.Enum]WorkoutCategory.Enum {

	routine := make(map[DayOfWeek.Enum]WorkoutCategory.Enum)
	routine[DayOfWeek.Monday] = WorkoutCategory.Rest
	routine[DayOfWeek.Tuesday] = WorkoutCategory.Rest
	routine[DayOfWeek.Wednesday] = WorkoutCategory.Rest
	routine[DayOfWeek.Thursday] = WorkoutCategory.Rest
	routine[DayOfWeek.Friday] = WorkoutCategory.Rest
	routine[DayOfWeek.Saturday] = WorkoutCategory.Rest
	routine[DayOfWeek.Sunday] = WorkoutCategory.Rest

	if raw["monday"] != nil {
		routine[DayOfWeek.Monday] = parseWorkoutCategory(raw["monday"].(string))
	}
	if raw["tuesday"] != nil {
		routine[DayOfWeek.Tuesday] = parseWorkoutCategory(raw["tuesday"].(string))
	}
	if raw["wednesday"] != nil {
		routine[DayOfWeek.Wednesday] = parseWorkoutCategory(raw["wednesday"].(string))
	}
	if raw["thursday"] != nil {
		routine[DayOfWeek.Thursday] = parseWorkoutCategory(raw["thursday"].(string))
	}
	if raw["friday"] != nil {
		routine[DayOfWeek.Friday] = parseWorkoutCategory(raw["friday"].(string))
	}
	if raw["saturday"] != nil {
		routine[DayOfWeek.Saturday] = parseWorkoutCategory(raw["saturday"].(string))
	}
	if raw["sunday"] != nil {
		routine[DayOfWeek.Sunday] = parseWorkoutCategory(raw["sunday"].(string))
	}

	return routine

}

type ParseGeneralResult struct {
	Results []interface{} `json:"results"`
}
