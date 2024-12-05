package routes

import (
	"backend/internal/container"
	"backend/internal/http/handlers"
	"github.com/labstack/echo/v4"
)

type AuthRoute struct {
	authHandler *handlers.AuthHandler
	router      *echo.Group
}

func NewAuthRoute(container *container.Container, router *echo.Group) *AuthRoute {
	return &AuthRoute{
		authHandler: container.AuthHandler,
		router:      router,
	}
}

func (r *AuthRoute) Register() {
	auth := r.router.Group("/auth")

	auth.POST("/sign-in", r.authHandler.SignIn)
	auth.POST("/sign-up", r.authHandler.SignUp)
}
