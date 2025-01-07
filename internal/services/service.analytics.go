package services

import (
	"backend/internal/http/data_transfers"
	"backend/pkg/format"
	"fmt"
	"github.com/jinzhu/copier"
	"net/http"
	"time"
)

type AnalyticsService struct {
	exerciseSetsRepository ExerciseSetsRepository
	sessionsRepository     SessionsRepository
}

func NewAnalyticsService(exerciseSetsRepository ExerciseSetsRepository, sessionsRepository SessionsRepository) *AnalyticsService {
	return &AnalyticsService{
		exerciseSetsRepository: exerciseSetsRepository,
		sessionsRepository:     sessionsRepository,
	}
}

func (s *AnalyticsService) FindTotalAnalytics(ownerID int, date time.Time) (data_transfers.DayWiseAnalyticsResponse, int, error) {
	var dayWiseAnalyticsResponse data_transfers.DayWiseAnalyticsResponse
	var sessionsResponse []data_transfers.SessionResponse
	exerciseSets, err := s.exerciseSetsRepository.FindAllByCreatedAt(ownerID, date)
	if err != nil {
		return data_transfers.DayWiseAnalyticsResponse{}, http.StatusInternalServerError, err
	}

	sessions, err := s.sessionsRepository.FindAllByStartTime(ownerID, date)
	if err != nil {
		return data_transfers.DayWiseAnalyticsResponse{}, http.StatusInternalServerError, err
	}

	var totalSessionsDuration time.Duration
	for _, session := range sessions {
		timeDiff := session.EndTime.Sub(session.StartTime)
		totalSessionsDuration += timeDiff
	}

	err = copier.Copy(&sessionsResponse, &sessions)
	if err != nil {
		return data_transfers.DayWiseAnalyticsResponse{}, http.StatusInternalServerError, err
	}

	details := make(map[string][]string)
	totalReps := 0
	for _, set := range exerciseSets {
		detailString := fmt.Sprintf("%d reps %.2f kg", set.Reps, set.Weight)
		totalReps += set.Reps
		details[set.WorkoutExercise.Exercise.Name] = append(details[set.WorkoutExercise.Exercise.Name], detailString)
	}
	totalSets := len(exerciseSets)

	dayWiseAnalyticsResponse.Reps = totalReps
	dayWiseAnalyticsResponse.Sets = totalSets
	dayWiseAnalyticsResponse.Exercises = len(details)
	dayWiseAnalyticsResponse.Details = details
	dayWiseAnalyticsResponse.Date = date
	dayWiseAnalyticsResponse.Sessions = sessionsResponse
	dayWiseAnalyticsResponse.SessionsTime = format.DurationToHoursMinutesSeconds(totalSessionsDuration)

	return dayWiseAnalyticsResponse, http.StatusOK, nil
}

func (s *AnalyticsService) FindTrainedDatesInDateRange(userID int, startDate time.Time, endDate time.Time) ([]data_transfers.ExerciseSetsResponse, []data_transfers.SessionResponse, int, error) {
	var exercisesSetResponse []data_transfers.ExerciseSetsResponse
	var sessionsResponse []data_transfers.SessionResponse

	exerciseSets, err := s.exerciseSetsRepository.FindAllInDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	sessions, err := s.sessionsRepository.FindAllInDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	err = copier.Copy(&sessionsResponse, &sessions)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	err = copier.Copy(&exercisesSetResponse, &exerciseSets)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	return exercisesSetResponse, sessionsResponse, http.StatusOK, nil
}
