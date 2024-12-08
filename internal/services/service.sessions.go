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

type SessionsRepository interface {
	FindAll() ([]records.Sessions, error)
	FindByID(id int) (records.Sessions, error)
	Save(session records.Sessions) (int, error)
	Update(id int, sessionMap map[string]interface{}) error
	Delete(id int) error
}

type SessionsService struct {
	repository SessionsRepository
}

func NewSessionsService(repository SessionsRepository) *SessionsService {
	return &SessionsService{repository}
}

func (s *SessionsService) FindAll() ([]data_transfers.SessionResponse, int, error) {
	var sessionResponses []data_transfers.SessionResponse

	sessions, err := s.repository.FindAll()
	if err != nil {
		return sessionResponses, http.StatusInternalServerError, err
	}

	err = copier.Copy(&sessionResponses, &sessions)
	if err != nil {
		return sessionResponses, http.StatusInternalServerError, err
	}

	return sessionResponses, http.StatusOK, nil
}

func (s *SessionsService) FindByID(id int) (data_transfers.SessionResponse, int, error) {
	var sessionResponse data_transfers.SessionResponse

	session, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return sessionResponse, http.StatusNotFound, errors.New("session not found")
		}
		return sessionResponse, http.StatusInternalServerError, err
	}

	err = copier.Copy(&sessionResponse, &session)
	if err != nil {
		return sessionResponse, http.StatusInternalServerError, err
	}

	return sessionResponse, http.StatusOK, nil
}

func (s *SessionsService) Save(createSessionRequest data_transfers.CreateSessionRequest) (int, int, error) {
	var session records.Sessions

	err := copier.Copy(&session, &createSessionRequest)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	id, err := s.repository.Save(session)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return id, http.StatusCreated, nil
}

func (s *SessionsService) Update(id int, updateSessionRequest data_transfers.UpdateSessionRequest) (int, error) {
	sessionMap, err := convert.StructToMap(updateSessionRequest)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = s.repository.Update(id, sessionMap)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *SessionsService) Delete(id int) (int, error) {
	err := s.repository.Delete(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
