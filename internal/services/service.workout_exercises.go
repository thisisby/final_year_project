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

type WorkoutExercisesRepository interface {
	FindAll() ([]records.WorkoutExercises, error)
	FindByID(id int) (records.WorkoutExercises, error)
	FindAllByWorkoutID(workoutID int) ([]records.WorkoutExercises, error)
	Save(workoutExercise records.WorkoutExercises) (int, error)
	Update(id int, workoutExercise map[string]interface{}) error
	Delete(id int) error
}

type WorkoutExercisesService struct {
	repository WorkoutExercisesRepository
}

func NewWorkoutExercisesService(repository WorkoutExercisesRepository) *WorkoutExercisesService {
	return &WorkoutExercisesService{repository}
}

func (s *WorkoutExercisesService) FindAll() ([]data_transfers.WorkoutExercisesResponse, int, error) {
	var workoutExercisesResponse []data_transfers.WorkoutExercisesResponse

	workoutExercises, err := s.repository.FindAll()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	err = copier.Copy(&workoutExercisesResponse, &workoutExercises)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return workoutExercisesResponse, http.StatusOK, nil
}

func (s *WorkoutExercisesService) FindByID(id int) (data_transfers.WorkoutExercisesResponse, int, error) {
	var workoutExerciseResponse data_transfers.WorkoutExercisesResponse

	workoutExercise, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return workoutExerciseResponse, http.StatusNotFound, errors.New("workout exercise not found")
		}
		return workoutExerciseResponse, http.StatusInternalServerError, err
	}

	err = copier.Copy(&workoutExerciseResponse, &workoutExercise)
	if err != nil {
		return workoutExerciseResponse, http.StatusInternalServerError, err
	}

	return workoutExerciseResponse, http.StatusOK, nil
}

func (s *WorkoutExercisesService) FindAllByWorkoutID(workoutID int) ([]data_transfers.WorkoutExercisesResponse, int, error) {
	var workoutExercisesResponse []data_transfers.WorkoutExercisesResponse

	workoutExercises, err := s.repository.FindAllByWorkoutID(workoutID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	err = copier.Copy(&workoutExercisesResponse, &workoutExercises)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return workoutExercisesResponse, http.StatusOK, nil
}

func (s *WorkoutExercisesService) Save(workoutExercise data_transfers.CreateWorkoutExercisesRequest) (int, int, error) {
	var workoutExercises records.WorkoutExercises

	err := copier.Copy(&workoutExercises, &workoutExercise)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	id, err := s.repository.Save(workoutExercises)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return id, http.StatusCreated, nil
}

func (s *WorkoutExercisesService) Update(id int, workoutExercise data_transfers.UpdateWorkoutExercisesRequest) (int, error) {
	workoutExerciseMap, err := convert.StructToMap(workoutExercise)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("service - Update - convert.StructToMap: %w", err)
	}

	err = s.repository.Update(id, workoutExerciseMap)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *WorkoutExercisesService) Delete(id int) (int, error) {
	err := s.repository.Delete(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
