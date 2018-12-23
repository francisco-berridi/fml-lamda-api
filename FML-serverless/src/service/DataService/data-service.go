package DataService

import (
	"domain"
	"service/parsehttp"
)

func New() *dataService {
	return &dataService{

	}
}

type dataService struct {

}


func (d *dataService) AllWorkoutSchedules()([]*domain.WorkoutSchedule, error) {
	return parsehttp.GetWorkoutSchedules()
}

func (d *dataService) AllWorkouts() ([]*domain.Workout, error) {
	return parsehttp.GetWorkouts()
}
