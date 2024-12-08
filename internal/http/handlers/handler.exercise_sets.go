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

type ExerciseSetsHandler struct {
	service *services.ExerciseSetsService
}

func NewExerciseSetsHandler(service *services.ExerciseSetsService) *ExerciseSetsHandler {
	return &ExerciseSetsHandler{
		service: service,
	}
}

func (h *ExerciseSetsHandler) Save(ctx echo.Context) error {
	var createExerciseSetsRequest data_transfers.CreateExerciseSetsRequest
	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)

	if err := helpers.BindAndValidate(ctx, &createExerciseSetsRequest); err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	createExerciseSetsRequest.OwnerID = jwtClaims.UserID
	id, statusCode, err := h.service.Save(createExerciseSetsRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "exercise set saved successfully", map[string]int{"id": id})
}

func (h *ExerciseSetsHandler) FindAllByWorkoutExerciseID(ctx echo.Context) error {
	var exerciseSets []data_transfers.ExerciseSetsResponse

	workoutExerciseIDStr := ctx.Param("workoutExerciseID")
	workoutExerciseID, err := convert.StringToInt(workoutExerciseIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid workout exercise ID")
	}

	exerciseSets, statusCode, err := h.service.FindByWorkoutExerciseID(workoutExerciseID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "exercise sets fetched successfully", exerciseSets)
}

func (h *ExerciseSetsHandler) Update(ctx echo.Context) error {
	var updateExerciseSetsRequest data_transfers.UpdateExerciseSetsRequest
	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)

	exerciseSetIDStr := ctx.Param("exerciseSetID")
	exerciseSetID, err := convert.StringToInt(exerciseSetIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid exercise-set ID")
	}

	if err := helpers.BindAndValidate(ctx, &updateExerciseSetsRequest); err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	exerciseSet, statusCode, err := h.service.FindByID(exerciseSetID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	if exerciseSet.OwnerID != jwtClaims.UserID && !jwtClaims.IsAdmin {
		return NewErrorResponse(ctx, http.StatusForbidden, "You are not allowed to update this exercise set")
	}

	statusCode, err = h.service.Update(exerciseSetID, updateExerciseSetsRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "exercise set updated successfully", nil)
}

func (h *ExerciseSetsHandler) Delete(ctx echo.Context) error {
	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)

	exerciseSetIDStr := ctx.Param("exerciseSetID")
	exerciseSetID, err := convert.StringToInt(exerciseSetIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid exercise-set ID")
	}

	exerciseSet, statusCode, err := h.service.FindByID(exerciseSetID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	if exerciseSet.OwnerID != jwtClaims.UserID && !jwtClaims.IsAdmin {
		return NewErrorResponse(ctx, http.StatusForbidden, "You are not allowed to delete this exercise set")
	}

	statusCode, err = h.service.Delete(exerciseSetID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "exercise set deleted successfully", nil)
}
