package services

import (
	"backend/internal/datasources/records"
	"backend/internal/datasources/repositories"
	"backend/internal/http/data_transfers"
	"backend/pkg/convert"
	"errors"
	"fmt"
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

func (s *WorkoutsService) FindAll() ([]data_transfers.WorkoutsResponse, int, error) {
	var workoutsResponse []data_transfers.WorkoutsResponse

	workouts, err := s.repository.FindAll()
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

func (s *WorkoutsService) FindByID(id int) (data_transfers.WorkoutsResponse, int, error) {
	var workoutResponse data_transfers.WorkoutsResponse

	workout, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return workoutResponse, http.StatusNotFound, errors.New("workout not found")
		}
		return workoutResponse, http.StatusInternalServerError, err
	}

	err = copier.Copy(&workoutResponse, &workout)
	if err != nil {
		return workoutResponse, http.StatusInternalServerError, err
	}

	return workoutResponse, http.StatusOK, nil
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

func (s *WorkoutsService) AddExercise(workoutID int, addExercisesRequest data_transfers.AddExercisesRequest) (int, error) {
	err := s.repository.AddExercise(workoutID, addExercisesRequest.ExerciseNames)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *WorkoutsService) RemoveExercise(workoutID int, deleteExerciseRequest data_transfers.DeleteExerciseRequest) (int, error) {
	err := s.repository.RemoveExercise(workoutID, deleteExerciseRequest.ExerciseIDs)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *WorkoutsService) Update(id int, workout data_transfers.UpdateWorkoutRequest) (int, error) {
	workoutMap, err := convert.StructToMap(workout)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("service - Update - convert.StructToMap: %w", err)
	}

	err = s.repository.Update(id, workoutMap)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *WorkoutsService) Delete(id int) (int, error) {
	err := s.repository.Delete(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
