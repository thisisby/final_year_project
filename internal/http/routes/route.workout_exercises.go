package routes

import (
	"backend/internal/container"
	"backend/internal/http/handlers"
	"backend/internal/http/middlewares"
	"github.com/labstack/echo/v4"
)

type WorkoutExercisesRoute struct {
	workoutExercisesHandler *handlers.WorkoutExercisesHandler
	router                  *echo.Group
}

func NewWorkoutExercisesRoute(container *container.Container, router *echo.Group) *WorkoutExercisesRoute {
	return &WorkoutExercisesRoute{
		workoutExercisesHandler: container.WorkoutExercisesHandler,
		router:                  router,
	}
}

func (r *WorkoutExercisesRoute) Register() {
	workoutExercises := r.router.Group("/workout-exercises")

	workoutExercises.Use(middlewares.RequireAuth)

	// workout_exercises routes
	workoutExercises.POST("", r.workoutExercisesHandler.Save)
	workoutExercises.GET("", r.workoutExercisesHandler.FindAll)
	workoutExercises.GET("/:workoutExerciseID", r.workoutExercisesHandler.FindByID)
	workoutExercises.PATCH("/:workoutExerciseID", r.workoutExercisesHandler.Update)
	workoutExercises.DELETE("/:workoutExerciseID", r.workoutExercisesHandler.Delete)
}
