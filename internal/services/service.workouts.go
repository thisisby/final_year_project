package services

import (
	"backend/internal/datasources/records"
	"backend/internal/datasources/repositories"
	"backend/internal/http/data_transfers"
	"errors"
	"github.com/jinzhu/copier"
	"net/http"
)

type WorkoutsRepository interface {
	FindAll() ([]records.Workouts, error)
	FindByID(id int) (records.Workouts, error)
	FindAllByOwnerID(ownerID int) ([]records.Workouts, error)
	Save(workout records.Workouts, exerciseNames []string) error
	AddExercise(workoutID int, exerciseNames []string) error
	RemoveExercise(workoutID int, exerciseIDs []int) error
	Update(id int, workout map[string]interface{}) error
	Delete(id int) error
}

type WorkoutsService struct {
	repository WorkoutsRepository
}

func NewWorkoutsService(repository WorkoutsRepository) *WorkoutsService {
	return &WorkoutsService{repository}
}

func (s *WorkoutsService) FindAllByOwnerID(ownerID int) ([]data_transfers.WorkoutsResponse, int, error) {
	var workoutsResponse []data_transfers.WorkoutsResponse

	workouts, err := s.repository.FindAllByOwnerID(ownerID)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return nil, http.StatusNotFound, errors.New("workouts not found")
		}
		return nil, http.StatusInternalServerError, err
	}

	err = copier.Copy(&workoutsResponse, &workouts)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return workoutsResponse, http.StatusOK, nil
}

func (s *WorkoutsService) Save(workout data_transfers.CreateWorkoutRequest) (int, error) {
	var workoutRecord records.Workouts
	err := copier.Copy(&workoutRecord, &workout)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = s.repository.Save(workoutRecord, workout.ExerciseNames)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}
