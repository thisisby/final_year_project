package handlers

import (
	"backend/internal/constants"
	"backend/internal/helpers"
	"backend/internal/http/data_transfers"
	"backend/internal/services"
	"backend/pkg/convert"
	"backend/pkg/jwt"
	"backend/pkg/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UsersHandler struct {
	service *services.UsersService
}

func NewUsersHandler(service *services.UsersService) *UsersHandler {
	return &UsersHandler{
		service: service,
	}
}

func (h *UsersHandler) FindAll(ctx echo.Context) error {
	var users []data_transfers.UsersResponse

	users, statusCode, err := h.service.FindAll()
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "users fetched successfully", users)

}

func (h *UsersHandler) FindByID(ctx echo.Context) error {
	var user data_transfers.UsersResponse

	idStr := ctx.Param("id")
	userId, err := convert.StringToInt(idStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	user, statusCode, err := h.service.FindByID(userId)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "user fetched successfully", user)

}

func (h *UsersHandler) Save(ctx echo.Context) error {
	var createUserRequest data_transfers.CreateUsersRequest

	err := helpers.BindAndValidate(ctx, &createUserRequest)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	statusCode, err := h.service.Save(createUserRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "user created successfully", nil)
}

func (h *UsersHandler) Update(ctx echo.Context) error {
	var updateUserRequest data_transfers.UpdateUsersRequest

	idStr := ctx.Param("id")
	userId, err := convert.StringToInt(idStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	err = helpers.BindAndValidate(ctx, &updateUserRequest)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	logger.ZeroLogger.Info().Msgf("updateUserRequest: %+v", updateUserRequest)

	statusCode, err := h.service.Update(userId, updateUserRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "user updated successfully", nil)
}

func (h *UsersHandler) Delete(ctx echo.Context) error {
	idStr := ctx.Param("id")
	userId, err := convert.StringToInt(idStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	err = h.service.Delete(userId)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return NewSuccessResponse(ctx, http.StatusOK, "user deleted successfully", nil)
}

func (h *UsersHandler) Me(ctx echo.Context) error {
	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)

	user, statusCode, err := h.service.FindByID(jwtClaims.UserID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "user fetched successfully", user)
}
