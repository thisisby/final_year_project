package handlers

import (
	"backend/internal/config"
	"backend/internal/constants"
	"backend/internal/helpers"
	"backend/internal/http/data_transfers"
	"backend/internal/services"
	"backend/internal/utils"
	"backend/pkg/convert"
	"backend/pkg/jwt"
	"backend/pkg/logger"
	"backend/third_party/s3"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UsersHandler struct {
	service  *services.UsersService
	s3Client *s3.Client
}

func NewUsersHandler(service *services.UsersService, s3Client *s3.Client) *UsersHandler {
	return &UsersHandler{
		service:  service,
		s3Client: s3Client,
	}
}

func (h *UsersHandler) FindAll(ctx echo.Context) error {
	var users []data_transfers.UsersResponse

	params, err := utils.ExtractQueryParams(ctx.QueryParams())
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	users, total, statusCode, err := h.service.FindAllWithFilters(params)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "users fetched successfully", map[string]interface{}{
		"data":  users,
		"total": total,
	})

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

func (h *UsersHandler) FindByUsername(ctx echo.Context) error {
	var user data_transfers.UsersResponse

	username := ctx.Param("username")

	user, statusCode, err := h.service.FindByUsername(username)
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

	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)
	idStr := ctx.Param("id")
	userId, err := convert.StringToInt(idStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	if jwtClaims.UserID != userId {
		return NewErrorResponse(ctx, http.StatusForbidden, "You are not allowed to update other user's data")
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

func (h *UsersHandler) ChangeAvatar(ctx echo.Context) error {
	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)
	var bucketName = config.Config.AWSBucketName

	avatar, err := ctx.FormFile("avatar")
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid avatar")
	}

	src, err := avatar.Open()
	if err != nil {
		fmt.Println(err)
		return NewErrorResponse(ctx, http.StatusInternalServerError, "Failed to open avatar file")
	}
	defer src.Close()

	url, err := h.s3Client.UploadFile(src, avatar.Filename, bucketName)
	if err != nil {
		fmt.Println(err)
		return NewErrorResponse(ctx, http.StatusInternalServerError, "Failed to upload avatar to S3")
	}

	statusCode, err := h.service.ChangeAvatar(jwtClaims.UserID, url)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "avatar updated successfully", map[string]string{"avatar": url})
}
