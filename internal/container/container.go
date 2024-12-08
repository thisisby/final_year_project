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
	ExerciseSetsRepository     services.ExerciseSetsRepository
	ActivityGroupsRepository   services.ActivityGroupsRepository
	ActivitiesRepository       services.ActivitiesRepository
	SessionsRepository         services.SessionsRepository
	SessionDetailsRepository   services.SessionDetailsRepository

	// Services
	UsersService            *services.UsersService
	AuthService             *services.AuthService
	ExercisesService        *services.ExercisesService
	WorkoutsService         *services.WorkoutsService
	WorkoutExercisesService *services.WorkoutExercisesService
	ExerciseSetsService     *services.ExerciseSetsService
	ActivityGroupsService   *services.ActivityGroupsService
	ActivitiesService       *services.ActivitiesService
	SessionsService         *services.SessionsService
	SessionDetailsService   *services.SessionDetailsService

	// Handlers
	UsersHandler            *handlers.UsersHandler
	AuthHandler             *handlers.AuthHandler
	ExercisesHandler        *handlers.ExercisesHandler
	WorkoutsHandler         *handlers.WorkoutsHandler
	WorkoutExercisesHandler *handlers.WorkoutExercisesHandler
	ExerciseSetsHandler     *handlers.ExerciseSetsHandler
	ActivityGroupsHandler   *handlers.ActivityGroupsHandler
	ActivitiesHandler       *handlers.ActivitiesHandler
	SessionsHandler         *handlers.SessionsHandler
	SessionDetailsHandler   *handlers.SessionDetailsHandler
}

func NewContainer(db *sqlx.DB) *Container {
	// Initialize repositories
	usersRepository := postgres.NewPostgresUsersRepository(db)
	tokenRepository := postgres.NewPostgresTokensRepository(db)
	exercisesRepository := postgres.NewPostgresExercisesRepository(db)
	workoutsRepository := postgres.NewPostgresWorkoutsRepository(db)
	workoutExercisesRepository := postgres.NewPostgresWorkoutExercisesRepository(db)
	exerciseSetsRepository := postgres.NewPostgresExerciseSetsRepository(db)
	activityGroupsRepository := postgres.NewPostgresActivityGroupsRepository(db)
	activitiesRepository := postgres.NewPostgresActivitiesRepository(db)
	sessionsRepository := postgres.NewPostgresSessionsRepository(db)
	sessionDetailsRepository := postgres.NewPostgresSessionDetailsRepository(db)

	// Initialize services
	usersService := services.NewUsersService(usersRepository)
	tokenService := services.NewTokensService(tokenRepository)
	authService := services.NewAuthService(usersService, tokenService)
	exercisesService := services.NewExercisesService(exercisesRepository)
	workoutExercisesService := services.NewWorkoutExercisesService(workoutExercisesRepository)
	workoutsService := services.NewWorkoutsService(workoutsRepository, workoutExercisesService)
	exerciseSetsService := services.NewExerciseSetsService(exerciseSetsRepository)
	activityGroupsService := services.NewActivityGroupsService(activityGroupsRepository)
	activitiesService := services.NewActivitiesService(activitiesRepository)
	sessionsService := services.NewSessionsService(sessionsRepository)
	sessionDetailsService := services.NewSessionDetailsService(sessionDetailsRepository)

	// Initialize handlers
	usersHandler := handlers.NewUsersHandler(usersService)
	authHandler := handlers.NewAuthHandler(authService)
	exercisesHandler := handlers.NewExercisesHandler(exercisesService)
	workoutsHandler := handlers.NewWorkoutsHandler(workoutsService)
	workoutExercisesHandler := handlers.NewWorkoutExercisesHandler(workoutExercisesService)
	exerciseSetsHandler := handlers.NewExerciseSetsHandler(exerciseSetsService)
	activityGroupsHandler := handlers.NewActivityGroupsHandler(activityGroupsService)
	activitiesHandler := handlers.NewActivitiesHandler(activitiesService)
	sessionsHandler := handlers.NewSessionsHandler(sessionsService)
	sessionDetailsHandler := handlers.NewSessionDetailsHandler(sessionDetailsService)

	return &Container{
		DB: db,

		// Repositories
		UsersRepository:            usersRepository,
		ExercisesRepository:        exercisesRepository,
		WorkoutsRepository:         workoutsRepository,
		WorkoutExercisesRepository: workoutExercisesRepository,
		ExerciseSetsRepository:     exerciseSetsRepository,
		ActivityGroupsRepository:   activityGroupsRepository,
		ActivitiesRepository:       activitiesRepository,
		SessionsRepository:         sessionsRepository,
		SessionDetailsRepository:   sessionDetailsRepository,

		// Services
		UsersService:            usersService,
		AuthService:             authService,
		ExercisesService:        exercisesService,
		WorkoutsService:         workoutsService,
		WorkoutExercisesService: workoutExercisesService,
		ExerciseSetsService:     exerciseSetsService,
		ActivityGroupsService:   activityGroupsService,
		ActivitiesService:       activitiesService,
		SessionsService:         sessionsService,
		SessionDetailsService:   sessionDetailsService,

		// Handlers
		UsersHandler:            usersHandler,
		AuthHandler:             authHandler,
		ExercisesHandler:        exercisesHandler,
		WorkoutsHandler:         workoutsHandler,
		WorkoutExercisesHandler: workoutExercisesHandler,
		ExerciseSetsHandler:     exerciseSetsHandler,
		ActivityGroupsHandler:   activityGroupsHandler,
		ActivitiesHandler:       activitiesHandler,
		SessionsHandler:         sessionsHandler,
		SessionDetailsHandler:   sessionDetailsHandler,
	}
}
