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
	UsersRepository            services.UsersRepository
	ExercisesRepository        services.ExercisesRepository
	WorkoutsRepository         services.WorkoutsRepository
	WorkoutExercisesRepository services.WorkoutExercisesRepository

	// Services
	UsersService            *services.UsersService
	AuthService             *services.AuthService
	ExercisesService        *services.ExercisesService
	WorkoutsService         *services.WorkoutsService
	WorkoutExercisesService *services.WorkoutExercisesService

	// Handlers
	UsersHandler            *handlers.UsersHandler
	AuthHandler             *handlers.AuthHandler
	ExercisesHandler        *handlers.ExercisesHandler
	WorkoutsHandler         *handlers.WorkoutsHandler
	WorkoutExercisesHandler *handlers.WorkoutExercisesHandler
}

func NewContainer(db *sqlx.DB) *Container {
	// Initialize repositories
	usersRepository := postgres.NewPostgresUsersRepository(db)
	tokenRepository := postgres.NewPostgresTokensRepository(db)
	exercisesRepository := postgres.NewPostgresExercisesRepository(db)
	workoutsRepository := postgres.NewPostgresWorkoutsRepository(db)
	workoutExercisesRepository := postgres.NewPostgresWorkoutExercisesRepository(db)

	// Initialize services
	usersService := services.NewUsersService(usersRepository)
	tokenService := services.NewTokensService(tokenRepository)
	authService := services.NewAuthService(usersService, tokenService)
	exercisesService := services.NewExercisesService(exercisesRepository)
	workoutExercisesService := services.NewWorkoutExercisesService(workoutExercisesRepository)
	workoutsService := services.NewWorkoutsService(workoutsRepository, workoutExercisesService)

	// Initialize handlers
	usersHandler := handlers.NewUsersHandler(usersService)
	authHandler := handlers.NewAuthHandler(authService)
	exercisesHandler := handlers.NewExercisesHandler(exercisesService)
	workoutsHandler := handlers.NewWorkoutsHandler(workoutsService)
	workoutExercisesHandler := handlers.NewWorkoutExercisesHandler(workoutExercisesService)

	return &Container{
		DB: db,

		// Repositories
		UsersRepository:            usersRepository,
		ExercisesRepository:        exercisesRepository,
		WorkoutsRepository:         workoutsRepository,
		WorkoutExercisesRepository: workoutExercisesRepository,

		// Services
		UsersService:            usersService,
		AuthService:             authService,
		ExercisesService:        exercisesService,
		WorkoutsService:         workoutsService,
		WorkoutExercisesService: workoutExercisesService,

		// Handlers
		UsersHandler:            usersHandler,
		AuthHandler:             authHandler,
		ExercisesHandler:        exercisesHandler,
		WorkoutsHandler:         workoutsHandler,
		WorkoutExercisesHandler: workoutExercisesHandler,
	}
}
