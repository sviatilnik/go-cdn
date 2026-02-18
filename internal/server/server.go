package server

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sviatilnik/go-cdn/internal/config"
	"github.com/sviatilnik/go-cdn/internal/httphandlers"
	"github.com/sviatilnik/go-cdn/internal/middlewares"
	"github.com/sviatilnik/go-cdn/internal/storage"

	_ "github.com/sviatilnik/go-cdn/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cnf *config.Config) *Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middlewares.GzipCompress)

	storage := storage.NewFSStorage("./files")
	/* storage, err := storage.NewS3Storage(context.Background())
	if err != nil {
		log.Fatalf("failed to create s3 storage: %v", err)
	} */

	r.Get("/api/v1/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/api/v1/docs/doc.json"),
	))

	r.Post("/api/v1/files/save", httphandlers.NewSaveFileHandler(storage).Handle())
	r.Delete("/api/v1/files/delete", httphandlers.NewDeleteFileHandler(storage).Handle())
	r.Get("/{folder:.*}/{filename:.*}", httphandlers.NewGetFileHandler(storage).Handle())

	serv := &http.Server{
		Addr:    cnf.Port,
		Handler: r,
	}

	return &Server{
		httpServer: serv,
	}
}

func (s *Server) Start(ctx context.Context) error {

	notifyContext, cancel := signal.NotifyContext(ctx, os.Kill, os.Interrupt)
	defer cancel()

	go func() {
		slog.Info("server starting ...")
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-notifyContext.Done()

	slog.Info("server shuting down ...")
	timeoutContext, timeoutContextCancel := context.WithTimeout(ctx, 10*time.Second)
	defer timeoutContextCancel()

	if err := s.httpServer.Shutdown(timeoutContext); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
		return err
	}

	slog.Info("server stopped")

	return nil
}
