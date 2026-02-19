package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/sviatilnik/go-cdn/internal/config"
	"github.com/sviatilnik/go-cdn/internal/server"
)

// @title Go CDN API
// @version 1.0
// @description API CDN
// @host localhost:8080
// @BasePath /
func main() {
	ctg, err := config.GetConfig()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	server := server.NewServer(ctg)
	server.Start(context.Background())
}
