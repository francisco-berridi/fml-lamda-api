package NutritionAnalytics

import (
	"formula/GetDietHall2015"
	"domain"
	"enum/units/Mass"
)

func NewCalculator() *nutritionCalculator {
	return &nutritionCalculator{}
}

type nutritionCalculator struct {

}

func (n *nutritionCalculator) RecommendMealPlans(bodyStats domain.BodyStats, activityStats domain.ActivityStats, initialBW Mass.Unit) {


	resultY, err := GetDietHall2015.Calculate(bodyStats.Age, bodyStats.Gender, bodyStats.Height.Metric(), activityStats.PAL, activityStats.CardioPart, bodyStats.RMR, )

}

