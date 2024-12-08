package routes

import (
	"backend/internal/container"
	"backend/internal/http/handlers"
	"backend/internal/http/middlewares"
	"github.com/labstack/echo/v4"
)

type ActivitiesRoute struct {
	activitiesHandler *handlers.ActivitiesHandler
	router            *echo.Group
}

func NewActivitiesRoute(container *container.Container, router *echo.Group) *ActivitiesRoute {
	return &ActivitiesRoute{
		activitiesHandler: container.ActivitiesHandler,
		router:            router,
	}
}

func (r *ActivitiesRoute) Register() {
	activities := r.router.Group("/activities")
	admin := r.router.Group("/admin/activities")

	activities.Use(middlewares.RequireAuth)
	admin.Use(middlewares.RequireAuth)

	// activities routes
	activities.GET("", r.activitiesHandler.FindAll)
	activities.GET("/:id", r.activitiesHandler.FindByID)

	// admin activities routes
	admin.POST("", r.activitiesHandler.Save)
	admin.PATCH("/:id", r.activitiesHandler.Update)
	admin.DELETE("/:id", r.activitiesHandler.Delete)

}
