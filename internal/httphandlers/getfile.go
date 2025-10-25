package httphandlers

import (
	"net/http"

	"github.com/sviatilnik/go-cdn/internal/storage"
)

type GetFileHandler struct {
	storage storage.Storage
}

func NewGetFileHandler(storage storage.Storage) *GetFileHandler {
	return &GetFileHandler{
		storage: storage,
	}
}

func (h *GetFileHandler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		path := r.URL.Query().Get("path")
		if path == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Path is required"))
			return
		}

		file, err := h.storage.GetFile(path)
		if err != nil {
			if err == storage.ErrFileNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.File.Close()

		http.ServeContent(w, r, file.Filename, file.Timestamp, file.File)
	}
}
