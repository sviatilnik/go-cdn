package httphandlers

import (
	"encoding/json"
	"net/http"

	"github.com/sviatilnik/go-cdn/internal/storage"
)

type deleteFileRequest struct {
	Path string `json:"path"`
}

type DeleteFileHandler struct {
	storage storage.Storage
}

func NewDeleteFileHandler(storage storage.Storage) *DeleteFileHandler {
	return &DeleteFileHandler{
		storage: storage,
	}
}

func (h *DeleteFileHandler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var req deleteFileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid request body"))
			return
		}

		err := h.storage.DeleteFile(r.Context(), req.Path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File deleted successfully"))
	}
}
