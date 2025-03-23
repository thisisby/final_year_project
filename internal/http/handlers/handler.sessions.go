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

type SessionsHandler struct {
	service *services.SessionsService
}

func NewSessionsHandler(service *services.SessionsService) *SessionsHandler {
	return &SessionsHandler{service}
}

func (h *SessionsHandler) FindAll(ctx echo.Context) error {
	sessionResponses, statusCode, err := h.service.FindAll()
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "sessions fetched successfully", sessionResponses)
}

func (h *SessionsHandler) FindByID(ctx echo.Context) error {
	sessionIDStr := ctx.Param("id")
	sessionID, err := convert.StringToInt(sessionIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	sessionResponse, statusCode, err := h.service.FindByID(sessionID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "session fetched successfully", sessionResponse)
}

func (h *SessionsHandler) FindAllByOwnerID(ctx echo.Context) error {
	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)
	ownerID := jwtClaims.UserID

	sessionResponses, statusCode, err := h.service.FindAllByOwnerID(ownerID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "sessions fetched successfully", sessionResponses)
}

func (h *SessionsHandler) Save(ctx echo.Context) error {
	var createSessionRequest data_transfers.CreateSessionRequest
	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)

	if err := helpers.BindAndValidate(ctx, &createSessionRequest); err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	createSessionRequest.OwnerID = jwtClaims.UserID

	id, statusCode, err := h.service.Save(createSessionRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "session saved successfully", map[string]int{"id": id})
}

func (h *SessionsHandler) Update(ctx echo.Context) error {
	sessionIDStr := ctx.Param("id")
	sessionID, err := convert.StringToInt(sessionIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	var updateSessionRequest data_transfers.UpdateSessionRequest
	if err := helpers.BindAndValidate(ctx, &updateSessionRequest); err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	statusCode, err := h.service.Update(sessionID, updateSessionRequest)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "session updated successfully", nil)
}

func (h *SessionsHandler) Delete(ctx echo.Context) error {
	sessionIDStr := ctx.Param("id")
	sessionID, err := convert.StringToInt(sessionIDStr)
	if err != nil {
		return NewErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	statusCode, err := h.service.Delete(sessionID)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	return NewSuccessResponse(ctx, statusCode, "session deleted successfully", nil)
}
