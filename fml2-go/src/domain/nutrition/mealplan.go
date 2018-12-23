package nutrition

import (
	"enums/DietType"
	"enums/Goal"
	"enums/MealCategory"
	"fmt"
)

type MealPlan struct {
	Id                   string
	Name                 string
	Goal                 Goal.Enum
	Restrictions         []DietType.Enum
	Summary              string
	Calories             float64
	FatPart              float64
	CarbsPart            float64
	ProteinPart          float64
	CaloriesFromFat      float64
	CaloriesFromProtein  float64
	CaloriesFromCarbs    float64
	WeightLossProjection []float64
	MealSchedule
}

func (m *MealPlan) ProjectedLossFinalWeight() float64 {
	return m.WeightLossProjection[len(m.WeightLossProjection)-1]
}

func (m *MealPlan) ProjectedWeeklyWeightLoss(weight float64) float64 {
	return 2.20462 * 7.0 * (weight - m.ProjectedLossFinalWeight()) / 100.0
}

func (m *MealPlan) DaysToTargetWeight(weight float64) int64 {

	for day, value := range m.WeightLossProjection {
		if value <= weight {
			return int64(day)
		}
	}

	return -1
}

func (m *MealPlan) ProjectedWeightAfterDay(day int64) (float64, error) {

	if day < 0 || int(day) > len(m.WeightLossProjection)-1 {
		return -1.0, fmt.Errorf("invalid day number %d, for weight loss curve", day)
	}

	return m.WeightLossProjection[int(day)], nil
}

func (m *MealPlan) CanMeetTarget(weight float64) bool {
	return m.DaysToTargetWeight(weight) != -1
}

type MealSchedule struct {
	Sunday    []*Meal
	Monday    []*Meal
	Tuesday   []*Meal
	Wednesday []*Meal
	Thursday  []*Meal
	Friday    []*Meal
	Saturday  []*Meal
}

type Meal struct {
	MealCategory MealCategory.Enum
}
