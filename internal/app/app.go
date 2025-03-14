package app

import (
	"log/slog"

	"github.com/c4erries/raptor-gateway/internal/app/server"
)

type App struct {
	Server *server.Server
}

func New(log *slog.Logger) *App {

	server := server.New(log)

	return &App{
		Server: server,
	}
}
