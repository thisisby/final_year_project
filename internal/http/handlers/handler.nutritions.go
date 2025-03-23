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

type NutritionsHandler struct {
	service *services.NutritionsService
}

func NewNutritionsHandler(service *services.NutritionsService) *NutritionsHandler {
	return &NutritionsHandler{service}
}

func (h *NutritionsHandler) FindAllByOwnerID(ctx echo.Context) error {
	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)
	ownerID := jwtClaims.UserID

	nutritionsResponse, statusCode, err := h.service.FindAllByOwnerID(ownerID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "nutritions fetched successfully", nutritionsResponse)
}

func (h *NutritionsHandler) FindByID(ctx echo.Context) error {
	nutritionIDStr := ctx.Param("id")
	nutritionID, err := convert.StringToInt(nutritionIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	nutritionResponse, statusCode, err := h.service.FindByID(nutritionID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "nutrition fetched successfully", nutritionResponse)
}

func (h *NutritionsHandler) Save(ctx echo.Context) error {
	var createNutritionRequest data_transfers.CreateNutritionsRequest
	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)

	if err := helpers.BindAndValidate(ctx, &createNutritionRequest); err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	createNutritionRequest.OwnerID = jwtClaims.UserID

	id, statusCode, err := h.service.Save(createNutritionRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "nutrition saved successfully", map[string]int{"id": id})
}

func (h *NutritionsHandler) Update(ctx echo.Context) error {
	nutritionIDStr := ctx.Param("id")
	nutritionID, err := convert.StringToInt(nutritionIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	var updateNutritionRequest data_transfers.UpdateNutritionsRequest
	if err := helpers.BindAndValidate(ctx, &updateNutritionRequest); err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	statusCode, err := h.service.Update(nutritionID, updateNutritionRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "nutrition updated successfully", nil)
}

func (h *NutritionsHandler) Delete(ctx echo.Context) error {
	nutritionIDStr := ctx.Param("id")
	nutritionID, err := convert.StringToInt(nutritionIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	err = h.service.Delete(nutritionID)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return NewSuccessResponse(ctx, http.StatusOK, "nutrition deleted successfully", nil)
}
