package main

import (
	"context"

	"github.com/sviatilnik/go-cdn/internal/config"
	"github.com/sviatilnik/go-cdn/internal/server"
)

func main() {
	ctg := &config.Config{
		Port: ":8090",
	}

	server := server.NewServer(ctg)
	server.Start(context.Background())
}
