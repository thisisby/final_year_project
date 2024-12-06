package routes

import (
	"backend/internal/container"
	"backend/internal/http/handlers"
	"backend/internal/http/middlewares"
	"github.com/labstack/echo/v4"
)

type UsersRoute struct {
	usersHandler *handlers.UsersHandler
	router       *echo.Group
}

func NewUsersRoute(container *container.Container, router *echo.Group) *UsersRoute {
	return &UsersRoute{
		usersHandler: container.UsersHandler,
		router:       router,
	}
}

func (r *UsersRoute) Register() {
	users := r.router.Group("/users")

	users.Use(middlewares.RequireAuth)
	users.GET("", r.usersHandler.FindAll)
	users.GET("/:id", r.usersHandler.FindByID)
	users.POST("", r.usersHandler.Save)
	users.PATCH("/:id", r.usersHandler.Update)
	users.DELETE("/:id", r.usersHandler.Delete)
}
