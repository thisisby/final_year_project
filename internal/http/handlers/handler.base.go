package handlers

import "github.com/labstack/echo/v4"

type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

func NewSuccessResponse(ctx echo.Context, statusCode int, message string, payload any) error {
	return ctx.JSON(statusCode, BaseResponse{
		Success: true,
		Message: message,
		Payload: payload,
	})
}

func NewErrorResponse(ctx echo.Context, statusCode int, message string) error {
	return ctx.JSON(statusCode, BaseResponse{
		Success: false,
		Message: message,
	})
}
