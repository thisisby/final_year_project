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
	users := r.router.Group("/users")
	admin := r.router.Group("/admin/exercises")
	workouts := r.router.Group("/workouts")

	exercises.Use(middlewares.RequireAuth)
	users.Use(middlewares.RequireAuth)
	admin.Use(middlewares.RequireAuth)
	workouts.Use(middlewares.RequireAuth)

	// users routes
	users.GET("/:userID/exercises", r.exercisesHandler.FindAllUserExercises)

	// exercises routes
	exercises.GET("", r.exercisesHandler.FindAll)
	exercises.GET("/:id", r.exercisesHandler.FindByID)
	exercises.POST("", r.exercisesHandler.CreateCustomExercise)
	exercises.PATCH("/:id", r.exercisesHandler.Update)
	exercises.DELETE("/:id", r.exercisesHandler.Delete)

	// admin routes
	admin.POST("", r.exercisesHandler.Save)

	// workouts routes
	workouts.GET("/:workoutID/exercises-included", r.exercisesHandler.FindAllWithWorkoutCheck)
}
