package services

import (
	"backend/internal/datasources/records"
	"backend/internal/datasources/repositories"
	"backend/internal/http/data_transfers"
	"backend/pkg/convert"
	"errors"
	"github.com/jinzhu/copier"
	"net/http"
)

type ActivitiesRepository interface {
	FindAll() ([]records.Activities, error)
	FindByID(id int) (records.Activities, error)
	Save(activity records.Activities) (int, error)
	Update(id int, activityMap map[string]interface{}) error
	Delete(id int) error
}

type ActivitiesService struct {
	repository ActivitiesRepository
}

func NewActivitiesService(repository ActivitiesRepository) *ActivitiesService {
	return &ActivitiesService{repository}
}

func (s *ActivitiesService) FindAll() ([]data_transfers.ActivityResponse, int, error) {
	var activitiesResponse []data_transfers.ActivityResponse

	activities, err := s.repository.FindAll()
	if err != nil {
		return activitiesResponse, http.StatusInternalServerError, err
	}

	err = copier.Copy(&activitiesResponse, &activities)
	if err != nil {
		return activitiesResponse, http.StatusInternalServerError, err
	}

	return activitiesResponse, http.StatusOK, nil
}

func (s *ActivitiesService) FindByID(id int) (data_transfers.ActivityResponse, int, error) {
	var activityResponse data_transfers.ActivityResponse

	activity, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return activityResponse, http.StatusNotFound, errors.New("activity not found")
		}
		return activityResponse, http.StatusInternalServerError, err
	}

	err = copier.Copy(&activityResponse, &activity)
	if err != nil {
		return activityResponse, http.StatusInternalServerError, err
	}

	return activityResponse, http.StatusOK, nil
}

func (s *ActivitiesService) Save(createActivityRequest data_transfers.CreateActivityRequest) (int, int, error) {
	var activity records.Activities

	err := copier.Copy(&activity, &createActivityRequest)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	id, err := s.repository.Save(activity)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return id, http.StatusCreated, nil
}

func (s *ActivitiesService) Update(id int, updateActivityRequest data_transfers.UpdateActivityRequest) (int, error) {
	activityMap, err := convert.StructToMap(updateActivityRequest)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = s.repository.Update(id, activityMap)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *ActivitiesService) Delete(id int) (int, error) {
	err := s.repository.Delete(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
