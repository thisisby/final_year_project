package services

import (
	"backend/internal/datasources/records"
	"backend/internal/datasources/repositories"
	"backend/internal/http/data_transfers"
	"backend/pkg/convert"
	"errors"
	"github.com/jinzhu/copier"
	"net/http"
	"time"
)

type ExerciseSetsRepository interface {
	Save(exerciseSet records.ExerciseSets) (int, error)
	FindAllByWorkoutExerciseID(workoutExerciseID int) ([]records.ExerciseSets, error)
	FindByID(id int) (records.ExerciseSets, error)
	Update(id int, exerciseSetMap map[string]interface{}) error
	Delete(id int) error
	FindAllByCreatedAt(ownerID int, createdAt time.Time) ([]records.ExerciseSets, error)
	FindAllInDateRange(ownerID int, startDate time.Time, endDate time.Time) ([]records.ExerciseSets, error)
	FindTotalSetsByDate(ownerID int, date time.Time) (int, error)
	FindTotalRepsByDate(ownerID int, date time.Time) (int, error)
	FindUniqueWorkoutExercisesByDate(ownerID int, date time.Time) (int, error)
	FindExercisesDetailsByDate(ownerID int, date time.Time) ([]records.ExerciseDetails, error)
}

type ExerciseSetsService struct {
	repository ExerciseSetsRepository
}

func NewExerciseSetsService(repository ExerciseSetsRepository) *ExerciseSetsService {
	return &ExerciseSetsService{repository}
}

func (s *ExerciseSetsService) Save(createExerciseSetsRequest data_transfers.CreateExerciseSetsRequest) (int, int, error) {
	var exerciseSet records.ExerciseSets

	err := copier.Copy(&exerciseSet, &createExerciseSetsRequest)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	id, err := s.repository.Save(exerciseSet)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return id, http.StatusCreated, nil
}

func (s *ExerciseSetsService) FindByWorkoutExerciseID(workoutExerciseID int) ([]data_transfers.ExerciseSetsResponse, int, error) {
	var exerciseSetsResponse []data_transfers.ExerciseSetsResponse

	exerciseSets, err := s.repository.FindAllByWorkoutExerciseID(workoutExerciseID)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return exerciseSetsResponse, http.StatusNotFound, errors.New("exercise sets not found")
		}

		return exerciseSetsResponse, http.StatusInternalServerError, err
	}

	err = copier.Copy(&exerciseSetsResponse, &exerciseSets)
	if err != nil {
		return exerciseSetsResponse, http.StatusInternalServerError, err
	}

	return exerciseSetsResponse, http.StatusOK, nil
}

func (s *ExerciseSetsService) FindByID(exerciseSetID int) (data_transfers.ExerciseSetsResponse, int, error) {
	var exerciseSetResponse data_transfers.ExerciseSetsResponse

	exerciseSet, err := s.repository.FindByID(exerciseSetID)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return exerciseSetResponse, http.StatusNotFound, errors.New("exercise set not found")
		}

		return exerciseSetResponse, http.StatusInternalServerError, err
	}

	err = copier.Copy(&exerciseSetResponse, &exerciseSet)
	if err != nil {
		return exerciseSetResponse, http.StatusInternalServerError, err
	}

	return exerciseSetResponse, http.StatusOK, nil
}

func (s *ExerciseSetsService) Update(id int, updateExerciseSetsRequest data_transfers.UpdateExerciseSetsRequest) (int, error) {
	exerciseSetMap, err := convert.StructToMap(updateExerciseSetsRequest)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = s.repository.Update(id, exerciseSetMap)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *ExerciseSetsService) Delete(exerciseSetID int) (int, error) {
	err := s.repository.Delete(exerciseSetID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
