package handlers

import (
	"backend/internal/constants"
	"backend/internal/helpers"
	"backend/internal/http/data_transfers"
	"backend/internal/services"
	"backend/internal/utils"
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

func (h *WorkoutsHandler) FindAll(ctx echo.Context) error {
	var workouts []data_transfers.WorkoutsResponse

	workouts, statusCode, err := h.service.FindAll()
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "workouts fetched successfully", workouts)
}

func (h *WorkoutsHandler) FindByID(ctx echo.Context) error {
	var workout data_transfers.WorkoutsResponse

	workoutIDStr := ctx.Param("id")
	workoutID, err := convert.StringToInt(workoutIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid workout ID")
	}

	workout, statusCode, err := h.service.FindByID(workoutID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "workout fetched successfully", workout)
}

func (h *WorkoutsHandler) FindAllByOwnerID(ctx echo.Context) error {
	var workouts []data_transfers.WorkoutsResponse
	var statusCode int
	var err error

	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)

	userIDStr := ctx.Param("userID")
	userID, err := convert.StringToInt(userIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID")
	}

	if jwtClaims.UserID == userID {
		workouts, statusCode, err = h.service.FindAllByCurrentUserID(userID)
		if err != nil {
			return NewErrorResponse(ctx, statusCode, err.Error())
		}
	} else {
		workouts, statusCode, err = h.service.FindAllByOwnerID(userID)
		if err != nil {
			return NewErrorResponse(ctx, statusCode, err.Error())
		}
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

	createWorkoutRequest.OwnerID = jwtClaims.UserID
	id, statusCode, err := h.service.Save(createWorkoutRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "workout created successfully", map[string]int{"id": id})
}

func (h *WorkoutsHandler) Update(ctx echo.Context) error {
	var updateWorkoutRequest data_transfers.UpdateWorkoutRequest
	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)

	workoutIDStr := ctx.Param("id")
	workoutID, err := convert.StringToInt(workoutIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid workout ID")
	}

	err = helpers.BindAndValidate(ctx, &updateWorkoutRequest)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	workout, statusCode, err := h.service.FindByID(workoutID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	if workout.OwnerID != jwtClaims.UserID && !jwtClaims.IsAdmin {
		return NewErrorResponse(ctx, http.StatusForbidden, "You are not allowed to update this workout")
	}

	statusCode, err = h.service.Update(workoutID, updateWorkoutRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "workout updated successfully", nil)
}

func (h *WorkoutsHandler) Delete(ctx echo.Context) error {
	idStr := ctx.Param("id")
	workoutID, err := convert.StringToInt(idStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid workout ID")
	}

	statusCode, err := h.service.Delete(workoutID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "workout deleted successfully", nil)
}

func (h *WorkoutsHandler) Copy(ctx echo.Context) error {
	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)

	idStr := ctx.Param("workoutID")
	workoutID, err := convert.StringToInt(idStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid workout ID")
	}

	workout, statusCode, err := h.service.FindByID(workoutID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	if workout.Price > float64(0) || workout.IsPrivate || workout.OwnerID != jwtClaims.UserID {
		return NewErrorResponse(ctx, http.StatusForbidden, "You cannot copy a paid or private workout")
	}

	id, statusCode, err := h.service.Copy(workoutID, jwtClaims.UserID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "workout copied successfully", map[string]int{"id": id})
}

func (h *WorkoutsHandler) PurchaseWorkout(ctx echo.Context) error {
	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)

	idStr := ctx.Param("workoutID")
	workoutID, err := convert.StringToInt(idStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid workout ID")
	}

	workout, statusCode, err := h.service.FindByID(workoutID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	if workout.IsPrivate || workout.Price <= float64(0) {
		return NewErrorResponse(ctx, http.StatusForbidden, "You cannot purchase a free or private workout")
	}

	id, statusCode, err := h.service.PurchaseWorkout(workoutID, jwtClaims.UserID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "workout purchased successfully", map[string]int{"id": id})
}

func (h *WorkoutsHandler) FindAllWithFilters(ctx echo.Context) error {
	var workouts []data_transfers.WorkoutsResponse

	params, err := utils.ExtractQueryParams(ctx.QueryParams())
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	workouts, total, statusCode, err := h.service.FindAllWithFilters(params)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, http.StatusOK, "Workouts fetched successfully", map[string]interface{}{
		"data":  workouts,
		"total": total,
	})
}
