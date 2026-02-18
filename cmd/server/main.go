package main

import (
	"context"

	"github.com/sviatilnik/go-cdn/internal/config"
	"github.com/sviatilnik/go-cdn/internal/server"
)

// @title Go CDN API
// @version 1.0
// @description API CDN
// @host localhost:8080
// @BasePath /
func main() {
	ctg := &config.Config{
		Port: ":8080",
	}

	server := server.NewServer(ctg)
	server.Start(context.Background())
}
