package httphandlers

import (
	"net/http"
)

func Healthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO add critical parts healthcheck (db, cache, etc.)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service is healthy"))
	}
}
