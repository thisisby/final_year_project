package routes

import (
	"backend/internal/container"
	"backend/internal/http/handlers"
	"backend/internal/http/middlewares"
	"github.com/labstack/echo/v4"
)

type SessionDetailsRoute struct {
	sessionDetailsHandler *handlers.SessionDetailsHandler
	router                *echo.Group
}

func NewSessionDetailsRoute(container *container.Container, router *echo.Group) *SessionDetailsRoute {
	return &SessionDetailsRoute{
		sessionDetailsHandler: container.SessionDetailsHandler,
		router:                router,
	}
}

func (r *SessionDetailsRoute) Register() {
	sessionDetails := r.router.Group("/session-details")
	admin := r.router.Group("/admin/session-details")
	sessions := r.router.Group("/sessions")

	sessionDetails.Use(middlewares.RequireAuth)
	admin.Use(middlewares.RequireAuth)
	sessions.Use(middlewares.RequireAuth)

	// session_details routes
	sessionDetails.GET("", r.sessionDetailsHandler.FindAll)
	sessionDetails.GET("/:id", r.sessionDetailsHandler.FindByID)

	// sessions routes
	sessions.GET("/:sessionID/session-details", r.sessionDetailsHandler.FindBySessionID)

	// admin session_details routes
	admin.POST("", r.sessionDetailsHandler.Save)
	admin.PATCH("/:id", r.sessionDetailsHandler.Update)
	admin.DELETE("/:id", r.sessionDetailsHandler.Delete)
}
