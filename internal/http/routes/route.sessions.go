package routes

import (
	"backend/internal/container"
	"backend/internal/http/handlers"
	"backend/internal/http/middlewares"
	"github.com/labstack/echo/v4"
)

type SessionsRoute struct {
	sessionsHandler *handlers.SessionsHandler
	router          *echo.Group
}

func NewSessionsRoute(container *container.Container, router *echo.Group) *SessionsRoute {
	return &SessionsRoute{
		sessionsHandler: container.SessionsHandler,
		router:          router,
	}
}

func (r *SessionsRoute) Register() {
	sessions := r.router.Group("/sessions")

	sessions.Use(middlewares.RequireAuth)

	// sessions routes
	sessions.GET("", r.sessionsHandler.FindAllByOwnerID)
	sessions.GET("/:id", r.sessionsHandler.FindByID)
	sessions.POST("", r.sessionsHandler.Save)
	sessions.PATCH("/:id", r.sessionsHandler.Update)
	sessions.DELETE("/:id", r.sessionsHandler.Delete)
}
