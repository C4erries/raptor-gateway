package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/c4erries/raptor-gateway/internal/app"
	"github.com/c4erries/raptor-gateway/internal/app/server"
)

const (
	logLocal = "local"
	logDev   = "dev"
	logProd  = "prod"
)

func main() {
	log := setupLogger("local")

	config := server.NewConfig()
	app := app.New(log)
	go app.Server.Start(config)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.Server.Stop()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case logLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case logDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case logProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
