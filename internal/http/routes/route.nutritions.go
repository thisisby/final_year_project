package routes

import (
	"backend/internal/container"
	"backend/internal/http/handlers"
	"backend/internal/http/middlewares"
	"github.com/labstack/echo/v4"
)

type NutritionsRoute struct {
	nutritionsHandler *handlers.NutritionsHandler
	router            *echo.Group
}

func NewNutritionsRoute(container *container.Container, router *echo.Group) *NutritionsRoute {
	return &NutritionsRoute{
		nutritionsHandler: container.NutritionsHandler,
		router:            router,
	}
}

func (r *NutritionsRoute) Register() {
	nutritions := r.router.Group("/nutritions")

	nutritions.Use(middlewares.RequireAuth)
	// nutritions routes
	nutritions.GET("", r.nutritionsHandler.FindAllByOwnerID)
	nutritions.GET("/:id", r.nutritionsHandler.FindByID)
	nutritions.POST("", r.nutritionsHandler.Save)
	nutritions.PATCH("/:id", r.nutritionsHandler.Update)
	nutritions.DELETE("/:id", r.nutritionsHandler.Delete)
}
