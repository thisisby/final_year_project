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

type SessionDetailsRepository interface {
	FindAll() ([]records.SessionDetails, error)
	FindByID(id int) (records.SessionDetails, error)
	Save(sessionDetail records.SessionDetails) (int, error)
	Update(id int, sessionDetail map[string]interface{}) error
	Delete(id int) error
	FindAllBySessionID(sessionID int) ([]records.SessionDetails, error)
}

type SessionDetailsService struct {
	repository SessionDetailsRepository
}

func NewSessionDetailsService(repository SessionDetailsRepository) *SessionDetailsService {
	return &SessionDetailsService{repository}
}

func (s *SessionDetailsService) FindAll() ([]data_transfers.CreateSessionDetailsResponse, int, error) {
	var sessionDetailsResponse []data_transfers.CreateSessionDetailsResponse

	sessionDetails, err := s.repository.FindAll()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	err = copier.Copy(&sessionDetailsResponse, &sessionDetails)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return sessionDetailsResponse, http.StatusOK, nil
}

func (s *SessionDetailsService) FindByID(id int) (data_transfers.CreateSessionDetailsResponse, int, error) {
	var sessionDetailResponse data_transfers.CreateSessionDetailsResponse

	sessionDetail, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return sessionDetailResponse, http.StatusNotFound, errors.New("session detail not found")
		}
		return sessionDetailResponse, http.StatusInternalServerError, err
	}

	err = copier.Copy(&sessionDetailResponse, &sessionDetail)
	if err != nil {
		return sessionDetailResponse, http.StatusInternalServerError, err
	}

	return sessionDetailResponse, http.StatusOK, nil
}

func (s *SessionDetailsService) Save(sessionDetail data_transfers.CreateSessionDetailsRequest) (int, int, error) {
	var sessionDetailRecord records.SessionDetails

	err := copier.Copy(&sessionDetailRecord, &sessionDetail)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	id, err := s.repository.Save(sessionDetailRecord)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return id, http.StatusCreated, nil
}

func (s *SessionDetailsService) Update(id int, sessionDetail data_transfers.UpdateSessionDetailsRequest) (int, error) {
	sessionDetailMap, err := convert.StructToMap(sessionDetail)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("service - Update - convert.StructToMap: %w", err)
	}

	err = s.repository.Update(id, sessionDetailMap)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *SessionDetailsService) Delete(id int) (int, error) {
	err := s.repository.Delete(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *SessionDetailsService) FindAllBySessionID(sessionID int) ([]data_transfers.CreateSessionDetailsResponse, int, error) {
	var sessionDetailsResponse []data_transfers.CreateSessionDetailsResponse

	sessionDetails, err := s.repository.FindAllBySessionID(sessionID)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return nil, http.StatusNotFound, errors.New("session details not found")
		}
		return nil, http.StatusInternalServerError, err
	}

	err = copier.Copy(&sessionDetailsResponse, &sessionDetails)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return sessionDetailsResponse, http.StatusOK, nil
}
