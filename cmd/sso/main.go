package main

import (
	"log/slog"
	"os"
	"sso/internal/config"
)

// собирает в себе все модули
func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting application:")

	// TODO: Инициализировать  приложение (/app)

	// TODO: запустить gRPC-сервер приложения
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// логгер надо выносить в отдельную функцию
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
