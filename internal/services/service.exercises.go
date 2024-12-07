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

type ExercisesRepository interface {
	FindAll() ([]records.Exercises, error)
	FindByID(id int) (records.Exercises, error)
	FindByName(name string) (records.Exercises, error)
	Save(exercise records.Exercises) error
	Update(id int, exercise map[string]interface{}) error
	Delete(id int) error
	CreateCustomExercise(exercise records.Exercises, userID int) (int, error)
	FindAllUserExercises(userID int) ([]records.Exercises, error)
	FindAllWithWorkoutCheck(workoutID int) ([]records.ExercisesWithWorkoutCheck, error)
}

type ExercisesService struct {
	repository ExercisesRepository
}

func NewExercisesService(repository ExercisesRepository) *ExercisesService {
	return &ExercisesService{repository}
}

func (s *ExercisesService) FindAll() ([]data_transfers.ExercisesResponse, int, error) {
	var exercisesResponse []data_transfers.ExercisesResponse
	exercises, err := s.repository.FindAll()
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return nil, http.StatusNotFound, errors.New("exercises not found")
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("service - FindAll - repository.FindAll: %w", err)
	}

	err = copier.Copy(&exercisesResponse, &exercises)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("service - FindAll - copier.Copy: %w", err)
	}

	return exercisesResponse, http.StatusOK, nil
}

func (s *ExercisesService) FindByID(id int) (data_transfers.ExercisesResponse, int, error) {
	var exerciseResponse data_transfers.ExercisesResponse
	exercise, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return exerciseResponse, http.StatusNotFound, errors.New("exercise not found")
		}
		return exerciseResponse, http.StatusInternalServerError, fmt.Errorf("service - FindByID - repository.FindByID: %w", err)
	}

	err = copier.Copy(&exerciseResponse, &exercise)
	if err != nil {
		return exerciseResponse, http.StatusInternalServerError, fmt.Errorf("service - FindByID - copier.Copy: %w", err)
	}

	return exerciseResponse, http.StatusOK, nil
}

func (s *ExercisesService) FindByName(name string) (data_transfers.ExercisesResponse, int, error) {
	var exerciseResponse data_transfers.ExercisesResponse
	exercise, err := s.repository.FindByName(name)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return exerciseResponse, http.StatusNotFound, errors.New("exercise not found")
		}
		return exerciseResponse, http.StatusInternalServerError, fmt.Errorf("service - FindByName - repository.FindByName: %w", err)
	}

	err = copier.Copy(&exerciseResponse, &exercise)
	if err != nil {
		return exerciseResponse, http.StatusInternalServerError, fmt.Errorf("service - FindByName - copier.Copy: %w", err)
	}

	return exerciseResponse, http.StatusOK, nil
}

func (s *ExercisesService) Save(exercise data_transfers.CreateExercisesRequest) (int, error) {
	var exerciseRecord records.Exercises
	err := copier.Copy(&exerciseRecord, &exercise)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("service - Save - copier.Copy: %w", err)
	}

	err = s.repository.Save(exerciseRecord)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("service - Save - repository.Save: %w", err)
	}

	return http.StatusCreated, nil
}

func (s *ExercisesService) Update(id int, exercise data_transfers.UpdateExercisesRequest) (int, error) {
	exerciseMap, err := convert.StructToMap(exercise)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("service - Update - convert.StructToMap: %w", err)
	}

	err = s.repository.Update(id, exerciseMap)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return http.StatusNotFound, errors.New("exercise not found")
		}
		return http.StatusInternalServerError, fmt.Errorf("service - Update - repository.Update: %w", err)
	}

	return http.StatusOK, nil
}

func (s *ExercisesService) Delete(id int) (int, error) {
	err := s.repository.Delete(id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("service - Delete - repository.Delete: %w", err)
	}

	return http.StatusOK, nil
}

func (s *ExercisesService) CreateCustomExercise(createExerciseRequest data_transfers.CreateExercisesRequest, userID int) (int, int, error) {
	var exerciseRecord records.Exercises
	err := copier.Copy(&exerciseRecord, &createExerciseRequest)
	if err != nil {
		return 0, http.StatusInternalServerError, fmt.Errorf("service - CreateCustomExercise - copier.Copy: %w", err)
	}

	id, err := s.repository.CreateCustomExercise(exerciseRecord, userID)
	if err != nil {
		return 0, http.StatusInternalServerError, fmt.Errorf("service - CreateCustomExercise - repository.Save: %w", err)
	}

	return id, http.StatusCreated, nil
}

func (s *ExercisesService) FindAllUserExercises(userID int) ([]data_transfers.ExercisesResponse, int, error) {
	var exercisesResponse []data_transfers.ExercisesResponse

	exercises, err := s.repository.FindAllUserExercises(userID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("service - FindAllUserExercises - repository.FindAllUserExercises: %w", err)
	}

	err = copier.Copy(&exercisesResponse, &exercises)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("service - FindAllUserExercises - copier.Copy: %w", err)
	}

	return exercisesResponse, http.StatusOK, nil
}

func (s *ExercisesService) FindAllWithWorkoutCheck(workoutID int) ([]data_transfers.ExercisesResponseWithWorkoutCheckResponse, int, error) {
	var exercisesResponse []data_transfers.ExercisesResponseWithWorkoutCheckResponse

	exercises, err := s.repository.FindAllWithWorkoutCheck(workoutID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("service - FindAllWithWorkoutCheck - repository.FindAllWithWorkoutCheck: %w", err)
	}

	err = copier.Copy(&exercisesResponse, &exercises)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("service - FindAllWithWorkoutCheck - copier.Copy: %w", err)
	}

	return exercisesResponse, http.StatusOK, nil
}
