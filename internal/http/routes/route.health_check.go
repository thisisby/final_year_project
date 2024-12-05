package routes

import (
	"backend/internal/http/handlers"
	"github.com/labstack/echo/v4"
)

type HealthCheckRoute struct {
	healthCheckHandler *handlers.HealthCheckHandler
	router             *echo.Group
}

func NewHealthCheckRoute(router *echo.Group) *HealthCheckRoute {
	healthCheckHandler := handlers.NewHealthCheckHandler()

	return &HealthCheckRoute{
		healthCheckHandler: healthCheckHandler,
		router:             router,
	}
}

func (r *HealthCheckRoute) Register() {
	healthCheck := r.router.Group("/health")

	healthCheck.GET("", r.healthCheckHandler.HealthCheck)
}
