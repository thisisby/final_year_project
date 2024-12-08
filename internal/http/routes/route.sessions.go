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
	admin := r.router.Group("/admin/sessions")

	sessions.Use(middlewares.RequireAuth)
	admin.Use(middlewares.RequireAuth)

	// sessions routes
	sessions.GET("", r.sessionsHandler.FindAll)
	sessions.GET("/:id", r.sessionsHandler.FindByID)

	// admin sessions routes
	admin.POST("", r.sessionsHandler.Save)
	admin.PATCH("/:id", r.sessionsHandler.Update)
	admin.DELETE("/:id", r.sessionsHandler.Delete)
}
