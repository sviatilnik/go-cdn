package httphandlers

import (
	"fmt"
	"net/http"
)

func Healthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO add critical parts healthcheck (db, cache, etc.)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Service is healthy")
	}
}
