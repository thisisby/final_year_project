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
	WorkoutsRepository  services.WorkoutsRepository

	// Services
	UsersService     *services.UsersService
	AuthService      *services.AuthService
	ExercisesService *services.ExercisesService
	WorkoutsService  *services.WorkoutsService

	// Handlers
	UsersHandler     *handlers.UsersHandler
	AuthHandler      *handlers.AuthHandler
	ExercisesHandler *handlers.ExercisesHandler
	WorkoutsHandler  *handlers.WorkoutsHandler
}

func NewContainer(db *sqlx.DB) *Container {
	// Initialize repositories
	usersRepository := postgres.NewPostgresUsersRepository(db)
	tokenRepository := postgres.NewPostgresTokensRepository(db)
	exercisesRepository := postgres.NewPostgresExercisesRepository(db)
	workoutsRepository := postgres.NewPostgresWorkoutsRepository(db)

	// Initialize services
	usersService := services.NewUsersService(usersRepository)
	tokenService := services.NewTokensService(tokenRepository)
	authService := services.NewAuthService(usersService, tokenService)
	exercisesService := services.NewExercisesService(exercisesRepository)
	workoutsService := services.NewWorkoutsService(workoutsRepository)

	// Initialize handlers
	usersHandler := handlers.NewUsersHandler(usersService)
	authHandler := handlers.NewAuthHandler(authService)
	exercisesHandler := handlers.NewExercisesHandler(exercisesService)
	workoutsHandler := handlers.NewWorkoutsHandler(workoutsService)

	return &Container{
		DB: db,

		// Repositories
		UsersRepository:     usersRepository,
		ExercisesRepository: exercisesRepository,
		WorkoutsRepository:  workoutsRepository,

		// Services
		UsersService:     usersService,
		AuthService:      authService,
		ExercisesService: exercisesService,
		WorkoutsService:  workoutsService,

		// Handlers
		UsersHandler:     usersHandler,
		AuthHandler:      authHandler,
		ExercisesHandler: exercisesHandler,
		WorkoutsHandler:  workoutsHandler,
	}
}
