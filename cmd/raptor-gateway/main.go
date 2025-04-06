package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/c4erries/raptor-gateway/internal/app"
	"github.com/c4erries/raptor-gateway/internal/config"
)

const (
	logLocal = "local"
	logDev   = "dev"
	logProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	fmt.Println(cfg)

	app := app.New(log)
	go app.Server.Start(cfg)
	//graceful shutdown
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
