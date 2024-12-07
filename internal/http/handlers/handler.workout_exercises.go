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

type WorkoutExercisesHandler struct {
	service *services.WorkoutExercisesService
}

func NewWorkoutExercisesHandler(service *services.WorkoutExercisesService) *WorkoutExercisesHandler {
	return &WorkoutExercisesHandler{
		service: service,
	}
}

func (h *WorkoutExercisesHandler) Save(ctx echo.Context) error {
	var createWorkoutExercisesRequest data_transfers.CreateWorkoutExercisesRequest
	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)

	if err := helpers.BindAndValidate(ctx, &createWorkoutExercisesRequest); err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	createWorkoutExercisesRequest.OwnerID = jwtClaims.UserID
	id, statusCode, err := h.service.Save(createWorkoutExercisesRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "workout exercise saved successfully", map[string]int{"id": id})
}

func (h *WorkoutExercisesHandler) FindAll(ctx echo.Context) error {
	var workoutExercises []data_transfers.WorkoutExercisesResponse

	workoutExercises, statusCode, err := h.service.FindAll()
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "workout exercises fetched successfully", workoutExercises)
}

func (h *WorkoutExercisesHandler) FindByID(ctx echo.Context) error {
	var workoutExercise data_transfers.WorkoutExercisesResponse

	workoutExerciseIDStr := ctx.Param("workoutExerciseID")
	workoutExerciseID, err := convert.StringToInt(workoutExerciseIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid workout exercise ID")
	}

	workoutExercise, statusCode, err := h.service.FindByID(workoutExerciseID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "workout exercise fetched successfully", workoutExercise)
}

func (h *WorkoutExercisesHandler) FindAllByWorkoutID(ctx echo.Context) error {
	var workoutExercises []data_transfers.WorkoutExercisesResponse

	workoutIDStr := ctx.Param("workoutID")
	workoutID, err := convert.StringToInt(workoutIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid workout ID")
	}

	workoutExercises, statusCode, err := h.service.FindAllByWorkoutID(workoutID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "workout exercises fetched successfully", workoutExercises)
}

func (h *WorkoutExercisesHandler) Update(ctx echo.Context) error {
	var updateWorkoutExercisesRequest data_transfers.UpdateWorkoutExercisesRequest
	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)

	workoutExerciseIDStr := ctx.Param("workoutExerciseID")
	workoutExerciseID, err := convert.StringToInt(workoutExerciseIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid workout exercise ID")
	}

	if err = helpers.BindAndValidate(ctx, &updateWorkoutExercisesRequest); err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	workoutExercise, statusCode, err := h.service.FindByID(workoutExerciseID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	if workoutExercise.OwnerID != jwtClaims.UserID && !jwtClaims.IsAdmin {
		return NewErrorResponse(ctx, http.StatusForbidden, "You are not authorized to update this workout exercise")
	}

	statusCode, err = h.service.Update(workoutExerciseID, updateWorkoutExercisesRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "workout exercise updated successfully", nil)
}

func (h *WorkoutExercisesHandler) Delete(ctx echo.Context) error {
	workoutExerciseIDStr := ctx.Param("workoutExerciseID")
	workoutExerciseID, err := convert.StringToInt(workoutExerciseIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid workout exercise ID")
	}

	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)
	workoutExercise, statusCode, err := h.service.FindByID(workoutExerciseID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	if workoutExercise.OwnerID != jwtClaims.UserID && !jwtClaims.IsAdmin {
		return NewErrorResponse(ctx, http.StatusForbidden, "You are not authorized to delete this workout exercise")
	}

	statusCode, err = h.service.Delete(workoutExerciseID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "workout exercise deleted successfully", nil)
}
