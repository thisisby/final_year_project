package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HealthCheckHandler struct {
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) HealthCheck(ctx echo.Context) error {
	return NewSuccessResponse(
		ctx,
		http.StatusOK,
		"OK",
		nil,
	)
}
