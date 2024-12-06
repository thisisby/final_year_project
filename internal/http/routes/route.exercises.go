package routes

import (
	"backend/internal/container"
	"backend/internal/http/handlers"
	"backend/internal/http/middlewares"
	"github.com/labstack/echo/v4"
)

type ExercisesRoute struct {
	exercisesHandler *handlers.ExercisesHandler
	router           *echo.Group
}

func NewExercisesRoute(container *container.Container, router *echo.Group) *ExercisesRoute {
	return &ExercisesRoute{
		exercisesHandler: container.ExercisesHandler,
		router:           router,
	}
}

func (r *ExercisesRoute) Register() {
	exercises := r.router.Group("/exercises")

	exercises.Use(middlewares.RequireAuth)
	exercises.GET("", r.exercisesHandler.FindAll)
	exercises.GET("/:id", r.exercisesHandler.FindByID)
	exercises.POST("", r.exercisesHandler.Save)
	exercises.PATCH("/:id", r.exercisesHandler.Update)
	exercises.DELETE("/:id", r.exercisesHandler.Delete)
}
