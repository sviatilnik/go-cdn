package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sviatilnik/go-cdn/internal/config"
	"github.com/sviatilnik/go-cdn/internal/httphandlers"
	"github.com/sviatilnik/go-cdn/internal/storage"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cnf *config.Config) *Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	//fsStorage := storage.NewFSStorage("./files")
	s3Storage, err := storage.NewS3Storage(context.Background())
	if err != nil {
		log.Fatalf("failed to create s3 storage: %v", err)
	}

	r.Post("/api/v1/files/save", httphandlers.NewSaveFileHandler(s3Storage).Handle())
	r.Delete("/api/v1/files/delete", httphandlers.NewDeleteFileHandler(s3Storage).Handle())
	r.Get("/{folder:.*}/{filename:.*}", httphandlers.NewGetFileHandler(s3Storage).Handle())

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
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-notifyContext.Done()

	timeoutContext, timeoutContextCancel := context.WithTimeout(ctx, 10*time.Second)
	defer timeoutContextCancel()

	if err := s.httpServer.Shutdown(timeoutContext); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
		return err
	}

	return nil
}
