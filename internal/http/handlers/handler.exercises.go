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

type ExercisesHandler struct {
	service *services.ExercisesService
}

func NewExercisesHandler(service *services.ExercisesService) *ExercisesHandler {
	return &ExercisesHandler{
		service: service,
	}
}

func (h *ExercisesHandler) FindAll(ctx echo.Context) error {
	var exercises []data_transfers.ExercisesResponse

	exercises, statusCode, err := h.service.FindAll()
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "exercises fetched successfully", exercises)
}

func (h *ExercisesHandler) FindByID(ctx echo.Context) error {
	var exercise data_transfers.ExercisesResponse

	idStr := ctx.Param("id")
	exerciseId, err := convert.StringToInt(idStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	exercise, statusCode, err := h.service.FindByID(exerciseId)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "exercise fetched successfully", exercise)
}

func (h *ExercisesHandler) FindByName(ctx echo.Context) error {
	var exercise data_transfers.ExercisesResponse

	name := ctx.Param("name")

	exercise, statusCode, err := h.service.FindByName(name)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "exercise fetched successfully", exercise)
}

func (h *ExercisesHandler) Save(ctx echo.Context) error {
	var createExerciseRequest data_transfers.CreateExercisesRequest

	err := helpers.BindAndValidate(ctx, &createExerciseRequest)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	statusCode, err := h.service.Save(createExerciseRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "exercise saved successfully", nil)
}

func (h *ExercisesHandler) Update(ctx echo.Context) error {
	var updateExerciseRequest data_transfers.UpdateExercisesRequest

	idStr := ctx.Param("id")
	exerciseId, err := convert.StringToInt(idStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	err = helpers.BindAndValidate(ctx, &updateExerciseRequest)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	statusCode, err := h.service.Update(exerciseId, updateExerciseRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "exercise updated successfully", nil)
}

func (h *ExercisesHandler) Delete(ctx echo.Context) error {
	idStr := ctx.Param("id")
	exerciseId, err := convert.StringToInt(idStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	statusCode, err := h.service.Delete(exerciseId)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "exercise deleted successfully", nil)
}

func (h *ExercisesHandler) CreateCustomExercise(ctx echo.Context) error {
	var createExerciseRequest data_transfers.CreateExercisesRequest
	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)

	err := helpers.BindAndValidate(ctx, &createExerciseRequest)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	id, statusCode, err := h.service.CreateCustomExercise(createExerciseRequest, jwtClaims.UserID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "custom exercise created successfully", map[string]int{"id": id})
}

func (h *ExercisesHandler) FindAllUserExercises(ctx echo.Context) error {
	var exercises []data_transfers.ExercisesResponse
	userIDS := ctx.Param("userID")
	userID, err := convert.StringToInt(userIDS)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	exercises, statusCode, err := h.service.FindAllUserExercises(userID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "exercises fetched successfully", exercises)
}

func (h *ExercisesHandler) FindAllWithWorkoutCheck(ctx echo.Context) error {
	var exercises []data_transfers.ExercisesResponseWithWorkoutCheckResponse

	workoutIDS := ctx.Param("workoutID")
	workoutID, err := convert.StringToInt(workoutIDS)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	exercises, statusCode, err := h.service.FindAllWithWorkoutCheck(workoutID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "exercises fetched successfully", exercises)
}
