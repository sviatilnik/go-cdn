package httphandlers

import (
	"log/slog"
	"net/http"
)

func Healthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO add critical parts healthcheck (db, cache, etc.)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Service is healthy"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error(err.Error())
			return
		}
	}
}
