package bootstrap

import (
	"backend/internal/config"
	"backend/internal/container"
	"backend/internal/datasources/drivers"
	"backend/internal/helpers"
	"backend/internal/http/routes"
	"backend/internal/utils"
	"backend/pkg/httpserver"
	"backend/pkg/jwt"
	"backend/pkg/logger"
	"backend/third_party/io"
	"backend/third_party/s3"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func MustRun() {
	logger.ZeroLogger.Info().Msg("Setting up default postgres connection...")
	sqlxOptions := utils.GetSqlxDriverOptions()
	conn, err := drivers.ConnectWithSQLX(sqlxOptions)
	if err != nil {
		logger.ZeroLogger.Fatal().Msgf("bootstrap - MustRun - drivers.SqlxConnect: %v", err)
	}
	defer conn.Close()

	logger.ZeroLogger.Info().Msg("Default postgres connection established.")

	s3Client, err := s3.NewClient(config.Config.AWSRegion, config.Config.AWSAccessKeyID, config.Config.AWSAccessKeySecret)
	if err != nil {
		logger.ZeroLogger.Fatal().Msgf("bootstrap - MustRun - s3.NewClient: %v", err)
	}

	ionet := io.NewClient(config.Config.IOAPIKey)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:  true,
		LogLatency: true,
		LogURI:     true,
		LogStatus:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.ZeroLogger.Info().
				Str("method", v.Method).
				Str("URI", v.URI).
				Int("status", v.Status).
				Str("latency", fmt.Sprintf("%dms", v.Latency.Milliseconds())).
				Msg("Request -> ")

			return nil
		},
	}))
	jwt.MustInitializeConfig(config.Config.JwtSecretKey, time.Minute*time.Duration(config.Config.JwtAccessTokenExpiresIn), time.Hour*time.Duration(config.Config.JwtRefreshTokenExpiresIn))
	e.Validator = helpers.NewValidator()

	v1 := e.Group("/api/v1")

	setupRoutes(v1, conn, s3Client, ionet)

	// running server
	logger.ZeroLogger.Info().Msg("Starting http server...")
	httpServer := httpserver.New(e, httpserver.Port(config.Config.Port))

	// waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.ZeroLogger.Info().Msg(fmt.Sprintf("app - Run - signal: " + s.String()))
	case err = <-httpServer.Notify():
		logger.ZeroLogger.Info().Msg(fmt.Errorf("app - Run - httpServer.Notify: %w", err).Error())
	}

	// Graceful shutdown
	logger.ZeroLogger.Info().Msg("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		logger.ZeroLogger.Fatal().Msg(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error())
	}
}

func setupRoutes(e *echo.Group, conn *sqlx.DB, s3Client *s3.Client, ionet *io.Client) {
	cont := container.NewContainer(conn, s3Client, ionet)

	// Register routes
	routes.NewUsersRoute(cont, e).Register()
	routes.NewAuthRoute(cont, e).Register()
	routes.NewExercisesRoute(cont, e).Register()
	routes.NewWorkoutsRoute(cont, e).Register()
	routes.NewWorkoutExercisesRoute(cont, e).Register()
	routes.NewExerciseSetsRoute(cont, e).Register()
	routes.NewActivityGroupsRoute(cont, e).Register()
	routes.NewActivitiesRoute(cont, e).Register()
	routes.NewSessionsRoute(cont, e).Register()
	routes.NewSessionDetailsRoute(cont, e).Register()
	routes.NewAnalyticsRoute(cont, e).Register()
	routes.NewNutritionsRoute(cont, e).Register()

	routes.NewHealthCheckRoute(e).Register()
}
