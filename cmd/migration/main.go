package main

import (
	"backend/internal/config"
	"backend/internal/datasources/drivers"
	"backend/internal/utils"
	"errors"
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	dir = "cmd/migration/migrations"
)

var (
	up   bool
	down bool
)

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Logger initialized")

	config.Config.MustInitializeConfig()
}

func main() {
	flag.BoolVar(&up, "up", false, "involves creating new tables, columns, or other database structures")
	flag.BoolVar(&down, "down", false, "involves dropping tables, columns, or other structures")
	flag.Parse()

	sqlxOptions := utils.GetSqlxDriverOptions()
	conn, err := drivers.ConnectWithSQLX(sqlxOptions)
	if err != nil {
		slog.Error("[Migration] failed to connect to db", err)
		return
	}
	defer conn.Close()

	if up {
		err = migrate(conn, "up")
		if err != nil {
			slog.Error("[Migration] failed to migrate up", err)
			return
		}
	}

	if down {
		err = migrate(conn, "down")
		if err != nil {
			slog.Error("[Migration] failed to migrate down", err)
			return
		}
	}
}

func migrate(db *sqlx.DB, action string) (err error) {
	slog.Info(fmt.Sprintf("[Migration]  running %s migration", action))

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(cwd, dir, fmt.Sprintf("*.%s.sql", action)))
	if err != nil {
		return errors.New("error when get files name")
	}

	for _, file := range files {
		slog.Info(fmt.Sprintf("[Migration] %s migration file: %v", action, file))
		data, err := os.ReadFile(file)
		if err != nil {
			return errors.New("error when read file")
		}

		_, err = db.Exec(string(data))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("error when exec query in file:%v", file)
		}
	}

	slog.Info(fmt.Sprintf("[Migration] %s migration success", action))

	return
}
