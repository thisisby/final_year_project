package services

import (
	"backend/internal/datasources/records"
	"backend/internal/helpers"
	"backend/internal/http/data_transfers"
	"backend/pkg/jwt"
	"fmt"
	"github.com/jinzhu/copier"
	"net/http"
)

type AuthService struct {
	UsersService  *UsersService
	TokensService *TokensService
}

func NewAuthService(usersService *UsersService, tokensService *TokensService) *AuthService {
	return &AuthService{
		UsersService:  usersService,
		TokensService: tokensService,
	}
}

func (s *AuthService) SignIn(signInRequest data_transfers.SignInRequest) (string, string, int, error) {
	user, statusCode, err := s.UsersService.FindByEmail(signInRequest.Email)
	if err != nil {
		return "", "", statusCode, err
	}

	if !helpers.ValidateHash(signInRequest.Password, user.Password) {
		return "", "", http.StatusBadRequest, nil
	}

	accessToken, refreshToken, err := jwt.GenerateTokenPair(user.ID, false)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	tokenRecord := records.Tokens{
		Token:  refreshToken,
		UserID: user.ID,
	}
	statusCode, err = s.TokensService.Save(tokenRecord)
	if err != nil {
		return "", "", statusCode, err
	}

	return accessToken, refreshToken, http.StatusOK, nil
}

func (s *AuthService) SignUp(signUpRequest data_transfers.SignUpRequest) (int, error) {
	var createUserRequest data_transfers.CreateUsersRequest

	err := copier.Copy(&createUserRequest, &signUpRequest)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("service - SignUp - copier.Copy: %w", err)
	}

	statusCode, err := s.UsersService.Save(createUserRequest)
	if err != nil {
		return statusCode, err
	}

	return http.StatusCreated, nil
}

func (s *AuthService) SignOut(token string) (int, error) {
	statusCode, err := s.TokensService.Delete(token)
	if err != nil {
		return statusCode, err
	}

	return http.StatusOK, nil
}
