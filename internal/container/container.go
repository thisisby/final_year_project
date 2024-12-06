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
	UsersRepository     services.UsersRepository
	ExercisesRepository services.ExercisesRepository

	// Services
	UsersService     *services.UsersService
	AuthService      *services.AuthService
	ExercisesService *services.ExercisesService

	// Handlers
	UsersHandler     *handlers.UsersHandler
	AuthHandler      *handlers.AuthHandler
	ExercisesHandler *handlers.ExercisesHandler
}

func NewContainer(db *sqlx.DB) *Container {
	// Initialize repositories
	usersRepository := postgres.NewPostgresUsersRepository(db)
	tokenRepository := postgres.NewPostgresTokensRepository(db)
	exercisesRepository := postgres.NewPostgresExercisesRepository(db)

	// Initialize services
	usersService := services.NewUsersService(usersRepository)
	tokenService := services.NewTokensService(tokenRepository)
	authService := services.NewAuthService(usersService, tokenService)
	exercisesService := services.NewExercisesService(exercisesRepository)

	// Initialize handlers
	usersHandler := handlers.NewUsersHandler(usersService)
	authHandler := handlers.NewAuthHandler(authService)
	exercisesHandler := handlers.NewExercisesHandler(exercisesService)

	return &Container{
		DB: db,

		// Repositories
		UsersRepository:     usersRepository,
		ExercisesRepository: exercisesRepository,

		// Services
		UsersService:     usersService,
		AuthService:      authService,
		ExercisesService: exercisesService,

		// Handlers
		UsersHandler:     usersHandler,
		AuthHandler:      authHandler,
		ExercisesHandler: exercisesHandler,
	}
}
