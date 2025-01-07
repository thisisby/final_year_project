package handlers

import (
	"backend/internal/helpers"
	"backend/internal/http/data_transfers"
	"backend/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) SignIn(ctx echo.Context) error {
	var signInRequest data_transfers.SignInRequest

	err := helpers.BindAndValidate(ctx, &signInRequest)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	accessToken, refreshToken, statusCode, err := h.service.SignIn(signInRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	helpers.WriteCookie(ctx, "refresh_token", refreshToken)

	return NewSuccessResponse(ctx, statusCode, "user signed in successfully", map[string]string{
		"access_token": accessToken,
	})
}

func (h *AuthHandler) SignUp(ctx echo.Context) error {
	var signUpRequest data_transfers.SignUpRequest

	err := helpers.BindAndValidate(ctx, &signUpRequest)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	statusCode, err := h.service.SignUp(signUpRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "user signed up successfully", nil)
}

func (h *AuthHandler) SignOut(ctx echo.Context) error {
	refreshToken, err := helpers.ReadCookie(ctx, "refresh_token")
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	statusCode, err := h.service.SignOut(refreshToken)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	helpers.WriteCookie(ctx, "refresh_token", "")

	return NewSuccessResponse(ctx, http.StatusOK, "user signed out successfully", nil)
}
