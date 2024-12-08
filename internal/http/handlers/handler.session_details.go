package handlers

import (
	"backend/internal/helpers"
	"backend/internal/http/data_transfers"
	"backend/internal/services"
	"backend/pkg/convert"
	"github.com/labstack/echo/v4"
	"net/http"
)

type SessionDetailsHandler struct {
	service *services.SessionDetailsService
}

func NewSessionDetailsHandler(service *services.SessionDetailsService) *SessionDetailsHandler {
	return &SessionDetailsHandler{
		service: service,
	}
}

func (h *SessionDetailsHandler) Save(ctx echo.Context) error {
	var createSessionDetailsRequest data_transfers.CreateSessionDetailsRequest

	if err := helpers.BindAndValidate(ctx, &createSessionDetailsRequest); err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	id, statusCode, err := h.service.Save(createSessionDetailsRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "session detail saved successfully", map[string]int{"id": id})
}

func (h *SessionDetailsHandler) FindAll(ctx echo.Context) error {
	var sessionDetails []data_transfers.CreateSessionDetailsResponse

	sessionDetails, statusCode, err := h.service.FindAll()
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "session details fetched successfully", sessionDetails)
}

func (h *SessionDetailsHandler) FindByID(ctx echo.Context) error {
	var sessionDetail data_transfers.CreateSessionDetailsResponse

	sessionDetailIDStr := ctx.Param("id")
	sessionDetailID, err := convert.StringToInt(sessionDetailIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid session detail ID")
	}

	sessionDetail, statusCode, err := h.service.FindByID(sessionDetailID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "session detail fetched successfully", sessionDetail)
}

func (h *SessionDetailsHandler) Update(ctx echo.Context) error {
	var updateSessionDetailsRequest data_transfers.UpdateSessionDetailsRequest

	sessionDetailIDStr := ctx.Param("id")
	sessionDetailID, err := convert.StringToInt(sessionDetailIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid session detail ID")
	}

	if err = helpers.BindAndValidate(ctx, &updateSessionDetailsRequest); err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	statusCode, err := h.service.Update(sessionDetailID, updateSessionDetailsRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "session detail updated successfully", nil)
}

func (h *SessionDetailsHandler) Delete(ctx echo.Context) error {
	sessionDetailIDStr := ctx.Param("id")
	sessionDetailID, err := convert.StringToInt(sessionDetailIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid session detail ID")
	}

	statusCode, err := h.service.Delete(sessionDetailID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "session detail deleted successfully", nil)
}

func (h *SessionDetailsHandler) FindBySessionID(ctx echo.Context) error {
	var sessionDetails []data_transfers.CreateSessionDetailsResponse

	sessionIDStr := ctx.Param("sessionID")
	sessionID, err := convert.StringToInt(sessionIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid session ID")
	}

	sessionDetails, statusCode, err := h.service.FindAllBySessionID(sessionID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "session details fetched successfully", sessionDetails)
}
