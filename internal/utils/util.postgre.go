package utils

import (
	"backend/internal/config"
	"backend/internal/datasources/drivers"
	"fmt"
	"time"
)

func GetSqlxDriverOptions() *drivers.DriverOptions {
	sqlxOptions := drivers.NewDriverOptions(
		"postgres",
		getPostgresDSN(
			config.Config.DBHost,
			config.Config.DBPort,
			config.Config.DBUser,
			config.Config.DBPassword,
			config.Config.DBName,
		),
		config.Config.DBMaxConn,
		config.Config.DBMaxIdle,
		time.Duration(config.Config.DBConnMaxLifeTime)*time.Minute,
	)

	return sqlxOptions
}

func GetDefaultPostgresDSN() string {
	return "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
}

func getPostgresDSN(
	host,
	port,
	user,
	password,
	dbname string,
) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}
