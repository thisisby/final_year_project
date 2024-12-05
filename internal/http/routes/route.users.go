package routes

import (
	"backend/internal/datasources/repositories/postgres"
	"backend/internal/http/handlers"
	"backend/internal/services"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type UsersRoute struct {
	usersHandler *handlers.UsersHandler
	router       *echo.Group
}

func NewUsersRoute(
	router *echo.Group,
	db *sqlx.DB,
) *UsersRoute {

	usersRepository := postgres.NewPostgresUsersRepository(db)
	usersService := services.NewUsersService(usersRepository)
	usersHandler := handlers.NewUsersHandler(usersService)

	return &UsersRoute{
		usersHandler: usersHandler,
		router:       router,
	}
}

func (r *UsersRoute) Register() {
	users := r.router.Group("/users")

	users.GET("", r.usersHandler.FindAll)
	users.GET("/:id", r.usersHandler.FindByID)
	users.POST("", r.usersHandler.Save)
	users.PATCH("/:id", r.usersHandler.Update)
	users.DELETE("/:id", r.usersHandler.Delete)
}
