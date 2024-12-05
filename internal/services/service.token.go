package services

import (
	"backend/internal/datasources/records"
	"backend/internal/datasources/repositories"
	"errors"
	"fmt"
	"net/http"
)

type TokensRepository interface {
	Save(token records.Tokens) error
	FindByToken(token string) (records.Tokens, error)
	Delete(token string) error
}

type TokensService struct {
	repository TokensRepository
}

func NewTokensService(repository TokensRepository) *TokensService {
	return &TokensService{repository}
}

func (s *TokensService) Save(token records.Tokens) (int, error) {
	err := s.repository.Save(token)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("service - Save - repository.Save: %w", err)
	}

	return http.StatusCreated, nil
}

func (s *TokensService) FindByToken(token string) (records.Tokens, int, error) {
	tokenRecord, err := s.repository.FindByToken(token)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return records.Tokens{}, http.StatusNotFound, errors.New("refresh token not found")
		}
		return records.Tokens{}, http.StatusInternalServerError, fmt.Errorf("service - FindByToken - repository.FindByToken: %w", err)
	}

	return tokenRecord, http.StatusOK, nil
}

func (s *TokensService) Delete(token string) (int, error) {
	err := s.repository.Delete(token)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("service - Delete - repository.Delete: %w", err)
	}

	return http.StatusOK, nil
}
