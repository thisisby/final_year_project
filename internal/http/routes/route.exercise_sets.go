package routes

import (
	"backend/internal/container"
	"backend/internal/http/handlers"
	"backend/internal/http/middlewares"
	"github.com/labstack/echo/v4"
)

type ExerciseSetsRoute struct {
	exerciseSetsHandler *handlers.ExerciseSetsHandler
	router              *echo.Group
}

func NewExerciseSetsRoute(container *container.Container, router *echo.Group) *ExerciseSetsRoute {
	return &ExerciseSetsRoute{
		exerciseSetsHandler: container.ExerciseSetsHandler,
		router:              router,
	}
}

func (r *ExerciseSetsRoute) Register() {
	exerciseSets := r.router.Group("/exercise-sets")
	workoutExercises := r.router.Group("/workout-exercises")

	exerciseSets.Use(middlewares.RequireAuth)
	workoutExercises.Use(middlewares.RequireAuth)

	// exercise_sets routes
	exerciseSets.POST("", r.exerciseSetsHandler.Save)
	exerciseSets.PATCH("/:exerciseSetID", r.exerciseSetsHandler.Update)
	exerciseSets.DELETE("/:exerciseSetID", r.exerciseSetsHandler.Delete)

	// workout_exercises routes
	workoutExercises.GET("/:workoutExerciseID/exercise-sets", r.exerciseSetsHandler.FindAllByWorkoutExerciseID)
}
