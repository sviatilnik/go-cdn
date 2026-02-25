package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/sviatilnik/go-cdn/internal/config"
	"github.com/sviatilnik/go-cdn/internal/server"
)

var (
	version string = "unknown"
	commit  string = "unknown"
	date    string = "unknown"
)

// @title Go CDN API
// @version 1.0
// @description API CDN
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	slog.Info(fmt.Sprintf("Version: %s, Commit: %s, Date: %s", version, commit, date))

	ctg, err := config.GetConfig()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	server := server.NewServer(ctg)
	err = server.Start(context.Background())
	if err != nil {
		slog.Error(err.Error())
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
