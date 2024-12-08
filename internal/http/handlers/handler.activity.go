package handlers

import (
	"backend/internal/helpers"
	"backend/internal/http/data_transfers"
	"backend/internal/services"
	"backend/pkg/convert"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ActivitiesHandler struct {
	service *services.ActivitiesService
}

func NewActivitiesHandler(service *services.ActivitiesService) *ActivitiesHandler {
	return &ActivitiesHandler{
		service: service,
	}
}

func (h *ActivitiesHandler) FindAll(ctx echo.Context) error {
	var activitiesResponse []data_transfers.ActivityResponse

	activitiesResponse, statusCode, err := h.service.FindAll()
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "activities fetched successfully", activitiesResponse)
}

func (h *ActivitiesHandler) FindByID(ctx echo.Context) error {
	var activityResponse data_transfers.ActivityResponse

	activityIDSrt := ctx.Param("id")
	activityID, err := convert.StringToInt(activityIDSrt)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	activityResponse, statusCode, err := h.service.FindByID(activityID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "activity fetched successfully", activityResponse)
}

func (h *ActivitiesHandler) Save(ctx echo.Context) error {
	var createActivityRequest data_transfers.CreateActivityRequest

	if err := helpers.BindAndValidate(ctx, &createActivityRequest); err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	id, statusCode, err := h.service.Save(createActivityRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "activity created successfully", id)
}

func (h *ActivitiesHandler) Update(ctx echo.Context) error {
	var updateActivityRequest data_transfers.UpdateActivityRequest

	if err := helpers.BindAndValidate(ctx, &updateActivityRequest); err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	activityIDSrt := ctx.Param("id")
	activityID, err := convert.StringToInt(activityIDSrt)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	statusCode, err := h.service.Update(activityID, updateActivityRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "activity updated successfully", nil)
}

func (h *ActivitiesHandler) Delete(ctx echo.Context) error {
	activityIDSrt := ctx.Param("id")
	activityID, err := convert.StringToInt(activityIDSrt)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	statusCode, err := h.service.Delete(activityID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "activity deleted successfully", nil)
}
