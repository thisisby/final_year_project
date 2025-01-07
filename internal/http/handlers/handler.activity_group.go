package handlers

import (
	"backend/internal/helpers"
	"backend/internal/http/data_transfers"
	"backend/internal/services"
	"backend/pkg/convert"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ActivityGroupsHandler struct {
	service *services.ActivityGroupsService
}

func NewActivityGroupsHandler(service *services.ActivityGroupsService) *ActivityGroupsHandler {
	return &ActivityGroupsHandler{
		service: service,
	}
}

func (h *ActivityGroupsHandler) FindAll(ctx echo.Context) error {
	var activityGroupsResponse []data_transfers.ActivityGroupResponse

	activityGroupsResponse, statusCode, err := h.service.FindAll()
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "activity groups fetched successfully", activityGroupsResponse)
}

func (h *ActivityGroupsHandler) FindByID(ctx echo.Context) error {
	var activityGroupResponse data_transfers.ActivityGroupResponse

	activityGroupIDSrt := ctx.Param("id")
	activityGroupID, err := convert.StringToInt(activityGroupIDSrt)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	activityGroupResponse, statusCode, err := h.service.FindByID(activityGroupID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "activity group fetched successfully", activityGroupResponse)
}

func (h *ActivityGroupsHandler) Save(ctx echo.Context) error {
	var createActivityGroupRequest data_transfers.CreateActivityGroupRequest

	if err := helpers.BindAndValidate(ctx, &createActivityGroupRequest); err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	id, statusCode, err := h.service.Save(createActivityGroupRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "activity group saved successfully", map[string]int{"id": id})
}

func (h *ActivityGroupsHandler) Update(ctx echo.Context) error {
	var updateActivityGroupRequest data_transfers.UpdateActivityGroupRequest

	if err := helpers.BindAndValidate(ctx, &updateActivityGroupRequest); err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	activityGroupIDSrt := ctx.Param("id")
	activityGroupID, err := convert.StringToInt(activityGroupIDSrt)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	statusCode, err := h.service.Update(activityGroupID, updateActivityGroupRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "activity group updated successfully", nil)
}

func (h *ActivityGroupsHandler) Delete(ctx echo.Context) error {
	activityGroupIDSrt := ctx.Param("id")
	activityGroupID, err := convert.StringToInt(activityGroupIDSrt)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	statusCode, err := h.service.Delete(activityGroupID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "activity group deleted successfully", nil)
}
