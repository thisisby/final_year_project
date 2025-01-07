package routes

import (
	"backend/internal/container"
	"backend/internal/http/handlers"
	"backend/internal/http/middlewares"
	"github.com/labstack/echo/v4"
)

type AnalyticsRoute struct {
	analyticsHandler *handlers.AnalyticsHandler
	router           *echo.Group
}

func NewAnalyticsRoute(container *container.Container, router *echo.Group) *AnalyticsRoute {
	return &AnalyticsRoute{
		analyticsHandler: container.AnalyticsHandler,
		router:           router,
	}
}

func (r *AnalyticsRoute) Register() {
	analytics := r.router.Group("/analytics")

	analytics.Use(middlewares.RequireAuth)
	analytics.GET("/day-wise", r.analyticsHandler.GetDayWiseAnalytics)
	analytics.GET("/training-days", r.analyticsHandler.GetTrainedDates)
}
