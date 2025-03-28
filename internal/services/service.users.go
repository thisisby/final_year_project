package services

import (
	"backend/internal/datasources/records"
	"backend/internal/datasources/repositories"
	"backend/internal/helpers"
	"backend/internal/http/data_transfers"
	"backend/pkg/convert"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"net/http"
	"strings"
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

	userRecord.Username = generateUsernameFromEmail(userRecord.Email)

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

func generateUsernameFromEmail(email string) string {
	// 1. Normalize the email by lowercasing and trimming spaces.
	normalizedEmail := strings.TrimSpace(strings.ToLower(email))

	// 2. Extract the part before the '@' symbol.
	atIndex := strings.Index(normalizedEmail, "@")
	if atIndex == -1 {
		// If '@' is not found, use the entire email.
		return generateFallbackUsername(normalizedEmail)
	}

	localPart := normalizedEmail[:atIndex]

	// 3. Clean up the local part by removing non-alphanumeric characters.
	var cleanedUsername strings.Builder
	for _, char := range localPart {
		if ('a' <= char && char <= 'z') || ('0' <= char && char <= '9') {
			cleanedUsername.WriteRune(char)
		}
	}
	username := cleanedUsername.String()

	// 4. Handle empty username after cleaning.
	if username == "" {
		return generateFallbackUsername(normalizedEmail)
	}

	// 5. Shorten if necessary and add a hash suffix to increase uniqueness.
	if len(username) > 20 {
		username = username[:20]
	}

	hash := md5.Sum([]byte(normalizedEmail))
	hashString := hex.EncodeToString(hash[:])
	suffix := hashString[:6] // Use a short hash to avoid extremely long usernames.

	return fmt.Sprintf("%s%s", username, suffix)
}

func generateFallbackUsername(input string) string {
	// Generate a fallback username using a hash of the input.
	hash := md5.Sum([]byte(input))
	hashString := hex.EncodeToString(hash[:])
	return "user" + hashString[:10]
}
