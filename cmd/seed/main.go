package main

import (
	"backend/internal/config"
	"backend/internal/datasources/drivers"
	"backend/internal/utils"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	dir = "cmd/seed/seeds"
)

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Logger initialized")

	config.Config.MustInitializeConfig()
}

func main() {

	sqlxOptions := utils.GetSqlxDriverOptions()
	conn, err := drivers.ConnectWithSQLX(sqlxOptions)
	if err != nil {
		slog.Error("[Migration] failed to connect to db", err)
		return
	}
	defer conn.Close()

	err = seed(conn)

}

func seed(db *sqlx.DB) (err error) {
	slog.Info(fmt.Sprintf("[Seed]  running seed"))

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(cwd, dir, fmt.Sprintf("*.seed.sql")))
	if err != nil {
		return errors.New("error when get files name")
	}

	for _, file := range files {
		slog.Info(fmt.Sprintf("[Seed] seed file: %v", file))
		data, err := os.ReadFile(file)
		if err != nil {
			return errors.New("error when read file")
		}

		_, err = db.Exec(string(data))
		if err != nil {
			slog.Error(fmt.Sprintf("[Seed] failed to seed file: %v", file), err)
		}
	}

	slog.Info(fmt.Sprintf("[Seed] seed success"))

	return
}
