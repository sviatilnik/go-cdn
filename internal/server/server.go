package server

import (
	"context"
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

	fsStorage := storage.NewFSStorage("./files")

	r.Post("/save-file", httphandlers.NewSaveFileHandler(fsStorage).Handle())
	r.Get("/get", httphandlers.NewGetFileHandler(fsStorage).Handle())
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

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
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-notifyContext.Done()

	timeoutContext, timeoutContextCancel := context.WithTimeout(ctx, 10*time.Second)
	defer timeoutContextCancel()

	if err := s.httpServer.Shutdown(timeoutContext); err != nil {
		return err
	}

	return nil
}
