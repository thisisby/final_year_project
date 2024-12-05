package container

import (
	"backend/internal/datasources/repositories/postgres"
	"backend/internal/http/handlers"
	"backend/internal/services"
	"github.com/jmoiron/sqlx"
)

type Container struct {
	DB *sqlx.DB

	// Repositories
	UsersRepository services.UsersRepository

	// Services
	UsersService *services.UsersService
	AuthService  *services.AuthService

	// Handlers
	UsersHandler *handlers.UsersHandler
	AuthHandler  *handlers.AuthHandler
}

func NewContainer(db *sqlx.DB) *Container {
	// Initialize repositories
	usersRepository := postgres.NewPostgresUsersRepository(db)
	tokenRepository := postgres.NewPostgresTokensRepository(db)

	// Initialize services
	usersService := services.NewUsersService(usersRepository)
	tokenService := services.NewTokensService(tokenRepository)
	authService := services.NewAuthService(usersService, tokenService)

	// Initialize handlers
	usersHandler := handlers.NewUsersHandler(usersService)
	authHandler := handlers.NewAuthHandler(authService)

	return &Container{
		DB: db,

		// Repositories
		UsersRepository: usersRepository,

		// Services
		UsersService: usersService,
		AuthService:  authService,

		// Handlers
		UsersHandler: usersHandler,
		AuthHandler:  authHandler,
	}
}
