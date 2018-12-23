package ActivityAnalytics

import (
	"domain"
	"domain/onboarding"
	"enum/ActivityLevel"
	"enum/DayOfWeek"
	"enum/Gender"
	"enum/Goal"
	"enum/Intensity"
	"enum/ObesityCategory"
	"enum/WorkoutCategory"
	"formula"
	"math"
	"service/DataService"
	"sort"
	"time"
)

func NewCalculator(stats domain.BodyStats, dailyActivityLevel ActivityLevel.Enum) *activityAnalyticsCalculator {
	return &activityAnalyticsCalculator{
		stats, dailyActivityLevel,
	}
}

type activityAnalyticsCalculator struct {
	bodyStats     domain.BodyStats
	activityLevel ActivityLevel.Enum
}

func (a *activityAnalyticsCalculator) RecommendWorkoutSchedule(obesityCategory ObesityCategory.Enum, primaryGoal, secondaryGoal Goal.Enum) (map[Intensity.Enum]*domain.WorkoutSchedule, error) {

	schedules := make(map[Intensity.Enum]*domain.WorkoutSchedule)

	ds := DataService.New()
	allSchedules, err := ds.AllWorkoutSchedules()
	if err != nil {
		return schedules, err
	}

	availableOptions := make([]*domain.WorkoutSchedule, 0)

	for _, value := range allSchedules {
		if value.ObesityCategory == obesityCategory &&
			value.PrimaryGoal == primaryGoal &&
			value.SecondaryGoal == secondaryGoal {

			availableOptions = append(availableOptions, value)
		}
	}
	// sort by number of activity days
	sort.Slice(availableOptions, func(i, j int) bool {

		schedule1 := availableOptions[i]
		schedule2 := availableOptions[j]

		return len(schedule1.ActivityDays()) < len(schedule2.ActivityDays())
	})

	if len(availableOptions) == 1 {
		schedules[Intensity.Recommended] = availableOptions[0]
	} else if len(availableOptions) == 2 {
		schedules[Intensity.Easy] = availableOptions[0]
		schedules[Intensity.Recommended] = availableOptions[1]
	} else if len(availableOptions) == 3 {
		schedules[Intensity.Easy] = availableOptions[0]
		schedules[Intensity.Recommended] = availableOptions[1]
		schedules[Intensity.Intense] = availableOptions[2]
	}

	return schedules, nil

}

func (a *activityAnalyticsCalculator) SelectWorkoutsForSchedule(schedule *domain.WorkoutSchedule, primaryGoal, secondaryGoal Goal.Enum, userHabits onboarding.Habits, body domain.BodyStats) error {

	ds := DataService.New()
	allWorkouts, err := ds.AllWorkouts()
	if err != nil {
		return err
	}

	scores := make(map[*domain.Workout]float64)

	for _, workout := range allWorkouts {

		score, err := a.generateWorkoutScore(workout, primaryGoal, secondaryGoal, userHabits, body)
		if err != nil {
			return err
		}

		scores[workout] = score
	}

	// Sort them
	sort.Slice(allWorkouts, func(i, j int) bool {
		return scores[allWorkouts[i]] > scores[allWorkouts[j]]
	})

	for day, category := range schedule.Routine {

		if category == WorkoutCategory.Rest {
			continue
		}
		if day == DayOfWeek.Monday {
			schedule.Monday = a.findWorkoutForCategory(category, allWorkouts, schedule)
		} else if day == DayOfWeek.Tuesday {
			schedule.Tuesday = a.findWorkoutForCategory(category, allWorkouts, schedule)
		} else if day == DayOfWeek.Wednesday {
			schedule.Wednesday = a.findWorkoutForCategory(category, allWorkouts, schedule)
		} else if day == DayOfWeek.Thursday {
			schedule.Thursday = a.findWorkoutForCategory(category, allWorkouts, schedule)
		} else if day == DayOfWeek.Friday {
			schedule.Friday = a.findWorkoutForCategory(category, allWorkouts, schedule)
		} else if day == DayOfWeek.Saturday {
			schedule.Saturday = a.findWorkoutForCategory(category, allWorkouts, schedule)
		} else if day == DayOfWeek.Sunday {
			schedule.Sunday = a.findWorkoutForCategory(category, allWorkouts, schedule)
		}
	}

	return nil

}

func (a *activityAnalyticsCalculator) generateWorkoutScore(workout *domain.Workout, primaryGoal, secondaryGoal Goal.Enum, userHabits onboarding.Habits, body domain.BodyStats) (float64, error) {

	//TODO: Validate data

	score := 0.0

	// We can adjust how weighted each factor
	// is depending of its importance (make sure the total is 1.0)
	goalsFactor := 0.40
	preferencesFactor := 0.20
	workoutsFactor := 0.20
	obesityFactor := 0.20

	/*
	 * Goals
	 */
	goalsScore := 0.0
	if workout.IsSuitableForGoal(primaryGoal) {
		// its way more important the primary goal , so it receives 80% of the goal score
		goalsScore += 0.8
	}

	if workout.IsSuitableForGoal(secondaryGoal) {
		goalsScore += 0.2
	}

	score += goalsScore * goalsFactor

	/*
	 * Workouts experience
	 */
	lengthScore := 1.0
	// The difference between the workout duration and the requested duration
	// people will prefer shorted workouts so same or less gets the best score (1.0), more
	// than that will decrease the score by .1 by every 3 (180 sec) minutes that its length increases
	// until it gets to 0 to lose the score
	difference := workout.Duration - userHabits.WorkoutSessionLength

	if difference > 0 {
		// this means that the workout is longer that desired
		lengthScore = 1 - (0.1 * (difference.Minutes() / 3))
		if lengthScore < 0 {
			lengthScore = 0
		}
	}

	requiresGymScore := 0.0
	// If the workout requires a gym but the user doesnt attend one, no score , if attends, full score
	// If the workout doesnt require a gym then full score because home workouts
	// can be done at the gym too
	if workout.RequiresGym && userHabits.AttendsGym {
		requiresGymScore = 1
	}
	experienceScore := 1.0

	// the farther the workout required experience and the users experience are
	// the less likely the user will want to do this workout, account for that
	experienceScore = 1 - math.Abs(float64(workout.TargetExperience.Compare(userHabits.GymExperience))*0.5)
	if experienceScore < 0 {
		experienceScore = 0
	}

	workoutScore := (lengthScore + requiresGymScore + experienceScore) / 3
	score += workoutScore * workoutsFactor

	/*
	 * Preferred Cardio Exercises factor
	 */
	cardioScore := 0.0
	if len(workout.CardioOptions) != 0 {
		for _, preferred := range userHabits.PreferredCardio {
			for _, option := range workout.CardioOptions {
				if option == preferred {
					cardioScore += 0.5
				}
			}
		}
	}
	score += cardioScore * preferencesFactor

	obesityScore := 1.0
	/* obesity factor */
	if body.ObesityCategory.Compare(workout.BMIRestriction) > 0 {
		// Negative result means that the user can do the exercise
		// workout receives the full score (1.0), but if he user is not
		// too far from the recommended obesity restriction then the
		// workout will receive some score
		obesityScore = 1 - (0.5 * float64(body.ObesityCategory.Compare(workout.BMIRestriction)))
		if obesityScore < 0 {
			obesityScore = 0
		}
	}

	score += obesityScore * obesityFactor

	return score, nil
}

func (a *activityAnalyticsCalculator) findWorkoutForCategory(category WorkoutCategory.Enum, rankedWorkouts []*domain.Workout, schedule *domain.WorkoutSchedule) *domain.Workout {

	var found *domain.Workout

	for _, ranked := range rankedWorkouts {
		// We give preference to workouts that are unique∂
		// thinking that the user would prefer not to repeat
		// workouts during the week
		if ranked.Category == category && schedule.IsWorkoutUnique(ranked) {
			found = ranked
			break
		}
	}

	if found == nil {
		// Cannot find a workout that is in the category
		// AND also unique, so repeat the workout
		for _, ranked := range rankedWorkouts {
			if ranked.Category == category {
				found = ranked
				break
			}
		}
	}

	return found
}

// Calculates pop_pi
func (a *activityAnalyticsCalculator) CalculateCardioPart(averageActiveCalories, pal float64) float64 {

	cardioPart := ((a.activityLevel.NumericValue() - (pal * 0.12) - 1) * pal) / averageActiveCalories
	return math.Min(math.Max(cardioPart, 0), 1)
}

func (a *activityAnalyticsCalculator) estimateWeeklyCaloriesBurn(bodyStats domain.BodyStats, w *domain.WorkoutSchedule) (resting, active, total float64, error error) {

	if w.Sunday != nil {
		dayResting, dayActive, dayTotal, err := a.calculateCaloriesExpendInWorkout(w.Sunday)

		if err != nil {
			error = err
			return
		}

		resting += dayResting
		active += dayActive
		total += dayTotal
	}

	if w.Monday != nil {
		dayResting, dayActive, dayTotal, err := a.calculateCaloriesExpendInWorkout(w.Monday)

		if err != nil {
			error = err
			return
		}

		resting += dayResting
		active += dayActive
		total += dayTotal
	}

	if w.Tuesday != nil {
		dayResting, dayActive, dayTotal, err := a.calculateCaloriesExpendInWorkout(w.Tuesday)

		if err != nil {
			error = err
			return
		}

		resting += dayResting
		active += dayActive
		total += dayTotal
	}

	if w.Wednesday != nil {
		dayResting, dayActive, dayTotal, err := a.calculateCaloriesExpendInWorkout(w.Wednesday)

		if err != nil {
			error = err
			return
		}

		resting += dayResting
		active += dayActive
		total += dayTotal
	}

	if w.Thursday != nil {
		dayResting, dayActive, dayTotal, err := a.calculateCaloriesExpendInWorkout(w.Thursday)

		if err != nil {
			error = err
			return
		}

		resting += dayResting
		active += dayActive
		total += dayTotal
	}

	if w.Friday != nil {
		dayResting, dayActive, dayTotal, err := a.calculateCaloriesExpendInWorkout(w.Friday)

		if err != nil {
			error = err
			return
		}

		resting += dayResting
		active += dayActive
		total += dayTotal
	}

	if w.Saturday != nil {
		dayResting, dayActive, dayTotal, err := a.calculateCaloriesExpendInWorkout( w.Saturday)

		if err != nil {
			error = err
			return
		}

		resting += dayResting
		active += dayActive
		total += dayTotal
	}

	return

}

func (a *activityAnalyticsCalculator) calculateCaloriesExpendInWorkout(workout *domain.Workout) (resting, active, total float64, err error) {

	// The intensity of the workout is based on the users
	// ability to perform this specific workout, estimating
	// what the heart rate of the user will perform this
	// workout
	estimatedIntensity := int64(math.Round(workout.Intensity * float64(a.bodyStats.MaxHeartRate)))

	// Calories per minute the user
	// burns
	burnedRestingCaloriesForMinute := a.bodyStats.MifflinStJeor / time.Duration(time.Hour*24).Minutes()

	// RMR per workout duration
	resting = burnedRestingCaloriesForMinute * workout.Duration.Minutes()

	if estimatedIntensity > 52 && workout.Category.MajorCategory() != WorkoutCategory.Strength {
		total = a.CaloriesBurnDuringActivity(workout.Duration)
	} else {
		// use the MET value
		total = float64(float64(workout.MetValue) * resting)
	}

	active = total - resting

	return
}

func (a *activityAnalyticsCalculator) CaloriesBurnPerMinuteOfActivity() float64 {

	if a.bodyStats.Gender == Gender.Female {
		return formula.MinuteOfActivityFemale(a.bodyStats.Weight.Metric(), a.bodyStats.AverageHRDuringWorkouts)
	} else {
		return formula.MinuteOfActivityMale(a.bodyStats.Weight.Metric(), a.bodyStats.AverageHRDuringWorkouts)
	}
}

func (a *activityAnalyticsCalculator) CaloriesBurnDuringActivity(activityDuration time.Duration) float64 {
	return a.CaloriesBurnPerMinuteOfActivity() * activityDuration.Minutes()
}

func (a *activityAnalyticsCalculator) CalculatePAL(averageWorkoutLength time.Duration, workoutsPerWeek int64) float64 {

	weeklyTimeWorkingOut := averageWorkoutLength.Minutes() * float64(workoutsPerWeek)
	averageTimeOfActivityPerDay := time.Duration(weeklyTimeWorkingOut / 7.0 * float64(time.Minute.Nanoseconds()))
	averageWeeklyActivityCalories := a.CaloriesBurnDuringActivity(averageTimeOfActivityPerDay)

	return a.activityLevel.NumericValue() + averageWeeklyActivityCalories/a.bodyStats.RMR

}
