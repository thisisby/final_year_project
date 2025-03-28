package services

import (
	"backend/internal/datasources/records"
	"backend/internal/datasources/repositories"
	"backend/internal/helpers"
	"backend/internal/http/data_transfers"
	"backend/pkg/convert"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"net/http"
	"strings"

	"github.com/Pallinder/go-randomdata"
)

type UsersRepository interface {
	FindAll() ([]records.Users, error)
	FindByID(id int) (records.Users, error)
	FindByEmail(email string) (records.Users, error)
	FindByUsername(username string) (records.Users, error)
	Save(user records.Users) error
	Update(id int, user map[string]interface{}) error
	ChangeAvatar(id int, avatar string) error
	Delete(id int) error
	UsernameExists(username string) (bool, error)
	FindAllWithFilters(params repositories.QueryParams) ([]records.Users, int, error)
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
		if errors.Is(err, repositories.ErrorRowNotFound) {
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

func (s *UsersService) FindAllWithFilters(params repositories.QueryParams) ([]data_transfers.UsersResponse, int, int, error) {
	var usersResponse []data_transfers.UsersResponse
	users, total, err := s.repository.FindAllWithFilters(params)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return nil, 0, http.StatusNotFound, errors.New("users not found")
		}
		return nil, 0, http.StatusInternalServerError, fmt.Errorf("service - FindAllWithFilters - repository.FindAllWithFilters: %w", err)
	}

	err = copier.Copy(&usersResponse, &users)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, fmt.Errorf("service - FindAllWithFilters - copier.Copy: %w", err)
	}

	return usersResponse, total, http.StatusOK, nil
}

func (s *UsersService) FindByID(id int) (data_transfers.UsersResponse, int, error) {
	var userResponse data_transfers.UsersResponse
	user, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return userResponse, http.StatusNotFound, errors.New("user not found")
		}
		return userResponse, http.StatusInternalServerError, fmt.Errorf("service - FindByID - repository.FindByID: %w", err)
	}

	err = copier.Copy(&userResponse, &user)

	return userResponse, http.StatusOK, nil
}

func (s *UsersService) FindByEmail(email string) (data_transfers.UsersResponse, int, error) {
	var userResponse data_transfers.UsersResponse
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return userResponse, http.StatusNotFound, errors.New("user not found")
		}
		return userResponse, http.StatusInternalServerError, fmt.Errorf("service - FindByEmail - repository.FindByEmail: %w", err)
	}

	err = copier.Copy(&userResponse, &user)

	return userResponse, http.StatusOK, nil
}

func (s *UsersService) FindByUsername(username string) (data_transfers.UsersResponse, int, error) {
	var userResponse data_transfers.UsersResponse
	user, err := s.repository.FindByUsername(username)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return userResponse, http.StatusNotFound, errors.New("user not found")
		}
		return userResponse, http.StatusInternalServerError, fmt.Errorf("service - FindByUsername - repository.FindByUsername: %w", err)
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

	userRecord.Username, err = s.generateRealishUniqueUsername()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("service - Save - generateRealishUniqueUsername: %w", err)
	}

	userRecord.Password, err = helpers.GenerateHash(userRecord.Password)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("service - Save - helpers.GenerateHash: %w", err)
	}
	err = s.repository.Save(userRecord)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowExists) {
			return http.StatusConflict, errors.New("user with this email already exists")
		}
		return http.StatusInternalServerError, fmt.Errorf("service - Save - repository.Save: %w", err)
	}

	return http.StatusCreated, nil
}

func (s *UsersService) Update(id int, user data_transfers.UpdateUsersRequest) (int, error) {
	userMap, err := convert.StructToMap(user)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("service - Update - convert.StructToMap: %w", err)
	}

	err = s.repository.Update(id, userMap)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return http.StatusNotFound, errors.New("user not found")
		}
		return http.StatusInternalServerError, fmt.Errorf("service - Update - repository.Update: %w", err)
	}

	return http.StatusOK, nil
}

func (s *UsersService) Delete(id int) error {
	err := s.repository.Delete(id)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("service - Delete - repository.Delete: %w", err)
	}

	return nil
}

func (s *UsersService) ChangeAvatar(id int, avatar string) (int, error) {
	err := s.repository.ChangeAvatar(id, avatar)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return http.StatusNotFound, errors.New("user not found")
		}
		return http.StatusInternalServerError, fmt.Errorf("service - ChangeAvatar - repository.ChangeAvatar: %w", err)
	}

	return http.StatusOK, nil
}

func (s *UsersService) generateRealishUniqueUsername() (string, error) {
	for i := 0; i < 10; i++ { // Attempt 10 times to generate a unique username
		firstName := strings.ToLower(randomdata.FirstName(randomdata.RandomGender))
		lastName := strings.ToLower(randomdata.LastName())
		randomNumber := randomdata.Number(1000) // Add a random number to increase uniqueness

		username := fmt.Sprintf("%s%s%d", firstName, lastName, randomNumber)

		// Check if the generated username already exists.
		if exists, err := s.repository.UsernameExists(username); err != nil {
			return "", fmt.Errorf("failed to check username existence: %w", err)
		} else if !exists {
			return username, nil
		}
	}

	return "", errors.New("failed to generate a unique username after multiple attempts") //Handle the case that after 10 tries, a unique username wasn't found.
}
