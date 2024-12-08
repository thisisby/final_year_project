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

type ActivityGroupsRepository interface {
	FindAll() ([]records.ActivityGroups, error)
	FindByID(id int) (records.ActivityGroups, error)
	Save(activityGroup records.ActivityGroups) (int, error)
	Update(id int, activityGroupMap map[string]interface{}) error
	Delete(id int) error
}

type ActivityGroupsService struct {
	repository ActivityGroupsRepository
}

func NewActivityGroupsService(repository ActivityGroupsRepository) *ActivityGroupsService {
	return &ActivityGroupsService{repository}
}

func (s *ActivityGroupsService) FindAll() ([]data_transfers.ActivityGroupResponse, int, error) {
	var activityGroupsResponse []data_transfers.ActivityGroupResponse

	activityGroups, err := s.repository.FindAll()
	if err != nil {
		return activityGroupsResponse, http.StatusInternalServerError, err
	}

	err = copier.Copy(&activityGroupsResponse, &activityGroups)
	if err != nil {
		return activityGroupsResponse, http.StatusInternalServerError, err
	}

	return activityGroupsResponse, http.StatusOK, nil
}

func (s *ActivityGroupsService) FindByID(id int) (data_transfers.ActivityGroupResponse, int, error) {
	var activityGroupResponse data_transfers.ActivityGroupResponse

	activityGroup, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return activityGroupResponse, http.StatusNotFound, errors.New("activity group not found")
		}
		return activityGroupResponse, http.StatusInternalServerError, err
	}

	err = copier.Copy(&activityGroupResponse, &activityGroup)
	if err != nil {
		return activityGroupResponse, http.StatusInternalServerError, err
	}

	return activityGroupResponse, http.StatusOK, nil
}

func (s *ActivityGroupsService) Save(createActivityGroupRequest data_transfers.CreateActivityGroupRequest) (int, int, error) {
	var activityGroup records.ActivityGroups

	err := copier.Copy(&activityGroup, &createActivityGroupRequest)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	id, err := s.repository.Save(activityGroup)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return id, http.StatusCreated, nil
}

func (s *ActivityGroupsService) Update(id int, updateActivityGroupRequest data_transfers.UpdateActivityGroupRequest) (int, error) {
	activityGroupMap, err := convert.StructToMap(updateActivityGroupRequest)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = s.repository.Update(id, activityGroupMap)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *ActivityGroupsService) Delete(id int) (int, error) {
	err := s.repository.Delete(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
