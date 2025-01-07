package handlers

import (
	"backend/internal/constants"
	"backend/internal/services"
	"backend/pkg/jwt"
	"backend/pkg/logger"
	"github.com/labstack/echo/v4"
	"time"
)

type AnalyticsHandler struct {
	service *services.AnalyticsService
}

func NewAnalyticsHandler(service *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}

func (h *AnalyticsHandler) GetDayWiseAnalytics(ctx echo.Context) error {

	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)

	params := ctx.QueryParam("date")
	layout := "2006-01-02"

	date, err := time.Parse(layout, params)
	if err != nil {
		return NewErrorResponse(ctx, 400, "Invalid date format. Expected format: YYYY-MM-DD")
	}

	dayWiseAnalyticsResponse, statusCode, err := h.service.FindTotalAnalytics(jwtClaims.UserID, date)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	logger.ZeroLogger.Info().Msg("Day Wise Analytics")

	return NewSuccessResponse(ctx, statusCode, "day wise analytics fetched successfully", dayWiseAnalyticsResponse)
}

func (h *AnalyticsHandler) GetTrainedDates(ctx echo.Context) error {

	jwtClaims := ctx.Get(constants.CtxAuthenticatedUserKey).(*jwt.Claims)

	paramsStartDate := ctx.QueryParam("start_date")
	paramsEndDate := ctx.QueryParam("end_date")
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, paramsStartDate)
	if err != nil {
		return NewErrorResponse(ctx, 400, "Invalid startDate format. Expected format: YYYY-MM-DD")
	}

	endDate, err := time.Parse(layout, paramsEndDate)
	if err != nil {
		return NewErrorResponse(ctx, 400, "Invalid endDate format. Expected format: YYYY-MM-DD")
	}

	exercisesSetsResponse, sessionsResponse, statusCode, err := h.service.FindTrainedDatesInDateRange(jwtClaims.UserID, startDate, endDate)
	if err != nil {
		return NewErrorResponse(ctx, statusCode, err.Error())
	}

	dates := make([]time.Time, 0)
	for _, exercisesSet := range exercisesSetsResponse {
		dates = append(dates, exercisesSet.CreatedAt)
	}
	for _, session := range sessionsResponse {
		dates = append(dates, session.StartTime)
	}

	return NewSuccessResponse(ctx, statusCode, "day wise analytics fetched successfully", dates)
}
