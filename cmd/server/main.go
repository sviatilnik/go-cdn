package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/sviatilnik/go-cdn/internal/auth"
	"github.com/sviatilnik/go-cdn/internal/config"
	"github.com/sviatilnik/go-cdn/internal/server"
	"github.com/sviatilnik/go-cdn/internal/user"
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

	auth := auth.NewAuthService(ctg.Auth.Issuer, ctg.Auth.Secret, ctg.Auth.Exp)
	token, err := auth.CreateAccessToken(&user.User{Name: "Sviatoslav", Email: "tregubov.sv@yandex.ru"})

	fmt.Println(token, err)

	server := server.NewServer(ctg)
	server.Start(context.Background())
}
