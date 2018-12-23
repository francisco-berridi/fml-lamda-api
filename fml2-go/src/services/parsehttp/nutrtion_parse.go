package parsehttp

import (
	"domain/nutrition"
	"encoding/json"
	"enums/DietType"
	"io/ioutil"
	"net/http"
	"fmt"
	"enums/Goal"
	"enums/MealCategory"
	"strings"
)

func GetMealPlans(dietType DietType.Enum) ([]*nutrition.MealPlan, error) {

	client := &http.Client{}
	all := make([]*nutrition.MealPlan, 0)

	req, err := http.NewRequest("GET", "https://dev.api.fitmylife.net/parse/classes/MealPlan?limit=999", nil)
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

		all = append(all, parseMealPlan(value.(map[string]interface{})))
	}

	filtered := make([]*nutrition.MealPlan, 0)

	if dietType != DietType.NoRestrictions {

		for _, plan := range all {

			for _, dt := range plan.Restrictions {
				if dt == dietType {
					filtered = append(filtered, plan)
					break
				}
			}
		}

	} else {
		filtered = all
	}

	return filtered, nil
}

func parseMealPlan(raw map[string]interface{}) *nutrition.MealPlan {

	plan := nutrition.MealPlan{
		Id:          raw["objectId"].(string),
		Calories:    raw["caloricValue"].(float64),
		Goal:        parseGoalWithParseName(raw["primaryGoalAlignment"].(string)),
		Name:        raw["name"].(string),
		Summary:     raw["summary"].(string),
		FatPart:     raw["totalFatProportionValue"].(float64),
		CarbsPart:   raw["carbohydrateProportionValue"].(float64),
		ProteinPart: raw["proteinProportionValue"].(float64),
	}

	plan.CaloriesFromCarbs = plan.Calories * plan.CarbsPart
	plan.CaloriesFromFat = plan.Calories * plan.FatPart
	plan.CaloriesFromProtein = plan.Calories * plan.ProteinPart

	restrictions := make([]DietType.Enum, 0)

	if raw["tags"] != nil {
		allTags := raw["tags"].([]interface{})

		for _, tag := range allTags {

			switch tag.(string) {
			case "Vegetarian":
				restrictions = append(restrictions, DietType.Vegetarian)
				break

				// TODO: complete the rest of the other diest
				// when they are available
			}
		}
	}

	plan.Restrictions = restrictions

	return &plan
}

func GetMeals(mealPlan *nutrition.MealPlan) (map[MealCategory.Enum][]*nutrition.Meal, error) {

	client := &http.Client{}
	//all := make([]*nutrition.Meal, 0)
	mealsByType := make(map[MealCategory.Enum][]*nutrition.Meal)

	fmt.Println("========", mealPlan.Goal, mealPlan.Calories)

	searchString := fmt.Sprintf("%.0f", mealPlan.Calories)
	goalString := fmt.Sprintf("%s", mealPlan.Goal)

	if mealPlan.Goal == Goal.WeightLoss {
		goalString = "WeightManagement"
	}

	searchString += goalString


	for _, dietType := range mealPlan.Restrictions {

		if dietType == DietType.NoRestrictions {
			continue
		}

		searchString += fmt.Sprintf("%s", dietType)
	}

	fmt.Println("------------------- searchString", searchString)


	req, err := http.NewRequest("GET", "https://dev.api.fitmylife.net/parse/classes/Meal?limit=999&where={\"mealPlanImportIdentifier\":\""+searchString+"\"}", nil)
	if err != nil {
		return mealsByType, err
	}

	req.Header.Add("X-Parse-Application-Id", "fitmylifeAppId")
	req.Header.Add("X-Parse-REST-API-Key", "flrcGF4lP4Lx6M2vq0rLo4MiAM2IN4BHRTr4QQKQ")

	resp, err := client.Do(req)
	if err != nil {
		return mealsByType, err
	}

	chars, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return mealsByType, err
	}


	result := ParseGeneralResult{}
	json.Unmarshal(chars, &result)

	fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$ MEAL")



	for _, raw := range result.Results {

		meal := parseMeal(raw.(map[string]interface{}))

		if _, ok := mealsByType[meal.MealCategory]; !ok {
			mealsByType[meal.MealCategory] = make([]*nutrition.Meal, 0)
		}

		mealsByType[meal.MealCategory] = append(mealsByType[meal.MealCategory], meal)
	}


	fmt.Println(MealCategory.Breakfast, len(mealsByType[MealCategory.Breakfast]))
	fmt.Println(MealCategory.Lunch, len(mealsByType[MealCategory.Lunch]))
	fmt.Println(MealCategory.Dinner, len(mealsByType[MealCategory.Dinner]))
	fmt.Println(MealCategory.Snack, len(mealsByType[MealCategory.Snack]))

	
	return mealsByType, nil

}


func parseMeal(raw map[string]interface{}) *nutrition.Meal {

	return &nutrition.Meal{
		MealCategory: parseMealCategory(raw["mealtimeIdentifier"].(string)),
	}
}


func parseMealCategory(name string) MealCategory.Enum {

	if strings.Contains(name, "Breakfast") {
		return MealCategory.Breakfast
	}

	if strings.Contains(name, "Lunch") {
		return MealCategory.Lunch
	}

	if strings.Contains(name, "Dinner") {
		return MealCategory.Dinner
	}

	return MealCategory.Snack
}