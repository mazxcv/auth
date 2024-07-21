package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"syscall"
)

// собирает в себе все модули
func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting application:")

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// чтение из канала блокирующая операция,  пока он пуст заставит висеть пока что-то туда не придет,
	sign := <-stop
	log.Info("stoppping application", slog.String("signal", sign.String()))
	application.GRPCSrv.Stop()
	log.Info("Application stop")
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
