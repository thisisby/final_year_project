package handlers

import (
	"backend/internal/constants"
	"backend/internal/helpers"
	"backend/internal/http/data_transfers"
	"backend/internal/services"
	"backend/pkg/convert"
	"backend/pkg/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type WorkoutsHandler struct {
	service *services.WorkoutsService
}

func NewWorkoutsHandler(service *services.WorkoutsService) *WorkoutsHandler {
	return &WorkoutsHandler{
		service: service,
	}
}

func (h *WorkoutsHandler) FindAllByOwnerID(ctx echo.Context) error {
	var workouts []data_transfers.WorkoutsResponse

	userIDStr := ctx.Param("userID")
	userID, err := convert.StringToInt(userIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID")
	}

	workouts, statusCode, err := h.service.FindAllByOwnerID(userID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "workouts fetched successfully", workouts)
}

func (h *WorkoutsHandler) Save(ctx echo.Context) error {
	var createWorkoutRequest data_transfers.CreateWorkoutRequest
	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)

	err := helpers.BindAndValidate(ctx, &createWorkoutRequest)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	err = createWorkoutRequest.Validate()
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	createWorkoutRequest.OwnerID = jwtClaims.UserID
	statusCode, err := h.service.Save(createWorkoutRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "workout created successfully", nil)
}
