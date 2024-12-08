package routes

import (
	"backend/internal/container"
	"backend/internal/http/handlers"
	"backend/internal/http/middlewares"
	"github.com/labstack/echo/v4"
)

type ActivityGroupsRoute struct {
	activityGroupsHandler *handlers.ActivityGroupsHandler
	router                *echo.Group
}

func NewActivityGroupsRoute(container *container.Container, router *echo.Group) *ActivityGroupsRoute {
	return &ActivityGroupsRoute{
		activityGroupsHandler: container.ActivityGroupsHandler,
		router:                router,
	}
}

func (r *ActivityGroupsRoute) Register() {
	activityGroups := r.router.Group("/activity-groups")
	admin := r.router.Group("/admin/activity-groups")

	activityGroups.Use(middlewares.RequireAuth)
	admin.Use(middlewares.RequireAuth)

	// activity_groups routes
	activityGroups.GET("", r.activityGroupsHandler.FindAll)
	activityGroups.GET("/:id", r.activityGroupsHandler.FindByID)

	// admin activity_groups routes
	admin.POST("", r.activityGroupsHandler.Save)
	admin.PATCH("/:id", r.activityGroupsHandler.Update)
	admin.DELETE("/:id", r.activityGroupsHandler.Delete)
}
