package routes

import (
	"backend/internal/container"
	"backend/internal/http/handlers"
	"backend/internal/http/middlewares"
	"github.com/labstack/echo/v4"
)

type WorkoutsRoute struct {
	workoutHandler *handlers.WorkoutsHandler
	router         *echo.Group
}

func NewWorkoutsRoute(container *container.Container, router *echo.Group) *WorkoutsRoute {
	return &WorkoutsRoute{
		workoutHandler: container.WorkoutsHandler,
		router:         router,
	}
}

func (r *WorkoutsRoute) Register() {
	workouts := r.router.Group("/workouts")
	users := r.router.Group("/users")

	users.Use(middlewares.RequireAuth)
	workouts.Use(middlewares.RequireAuth)

	// users routes
	users.GET("/:userID/workouts", r.workoutHandler.FindAllByOwnerID)

	// workouts routes
	workouts.POST("", r.workoutHandler.Save)
	workouts.GET("", r.workoutHandler.FindAllWithFilters)
	workouts.GET("/:id", r.workoutHandler.FindByID)
	workouts.PATCH("/:id", r.workoutHandler.Update)
	workouts.DELETE("/:id", r.workoutHandler.Delete)
	workouts.GET("/:workoutID/copy", r.workoutHandler.Copy)
	workouts.GET("/:workoutID/purchase", r.workoutHandler.PurchaseWorkout)
}
