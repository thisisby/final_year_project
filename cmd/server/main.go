package main

import (
	"backend/internal/bootstrap"
	"backend/internal/config"
	"backend/pkg/logger"
)

func init() {
	logger.InitZeroLogger()
	config.Config.MustInitializeConfig()
}

func main() {
	bootstrap.MustRun()
}
