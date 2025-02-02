package main

import (
	"github.com/c4erries/raptor-gateway/internal/server"
)

func main() {
	config := server.NewConfig()
	server.Start(config)
}
