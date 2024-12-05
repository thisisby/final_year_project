package services

import (
	"backend/internal/constants"
	"backend/internal/datasources/records"
	"backend/internal/http/data_transfers"
	"backend/pkg/convert"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"net/http"
)

type UsersRepository interface {
	FindAll() ([]records.Users, error)
	FindByID(id int) (records.Users, error)
	Save(user records.Users) error
	Update(id int, user map[string]interface{}) error
	Delete(id int) error
}

type UsersService struct {
	repository UsersRepository
}

func NewUsersService(repository UsersRepository) *UsersService {
	return &UsersService{repository}
}

func (s *UsersService) FindAll() ([]data_transfers.UsersResponse, int, error) {
	var usersResponse []data_transfers.UsersResponse
	users, err := s.repository.FindAll()
	if err != nil {
		if errors.Is(err, constants.ErrorRowNotFound) {
			return nil, http.StatusNotFound, errors.New("users not found")
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("service - FindAll - repository.FindAll: %w", err)
	}

	err = copier.Copy(&usersResponse, &users)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("service - FindAll - copier.Copy: %w", err)
	}

	return usersResponse, http.StatusOK, nil
}

func (s *UsersService) FindByID(id int) (data_transfers.UsersResponse, int, error) {
	var userResponse data_transfers.UsersResponse
	user, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, constants.ErrorRowNotFound) {
			return userResponse, http.StatusNotFound, errors.New("user not found")
		}
		return userResponse, http.StatusInternalServerError, fmt.Errorf("service - FindByID - repository.FindByID: %w", err)
	}

	err = copier.Copy(&userResponse, &user)

	return userResponse, http.StatusOK, nil
}

func (s *UsersService) Save(userRequest data_transfers.CreateUsersRequest) (int, error) {
	var userRecord records.Users

	err := copier.Copy(&userRecord, &userRequest)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("service - Save - copier.Copy: %w", err)
	}

	err = s.repository.Save(userRecord)
	if err != nil {
		if errors.Is(err, constants.ErrorRowExists) {
			return http.StatusConflict, errors.New("user with this email already exists")
		}
		return http.StatusInternalServerError, fmt.Errorf("service - Save - repository.Save: %w", err)
	}

	return http.StatusCreated, nil
}

func (s *UsersService) Update(id int, user data_transfers.UpdateUsersRequest) (int, error) {
	userMap, err := convert.StructToMap(user)

	err = s.repository.Update(id, userMap)
	if err != nil {
		if errors.Is(err, constants.ErrorRowNotFound) {
			return http.StatusNotFound, errors.New("user not found")
		}
		return http.StatusInternalServerError, fmt.Errorf("service - Update - repository.Update: %w", err)
	}

	return http.StatusOK, nil
}

func (s *UsersService) Delete(id int) error {
	err := s.repository.Delete(id)
	if err != nil {
		if errors.Is(err, constants.ErrorRowNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("service - Delete - repository.Delete: %w", err)
	}

	return nil
}
