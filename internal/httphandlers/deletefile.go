package httphandlers

import (
	"encoding/json"
	"net/http"

	"github.com/sviatilnik/go-cdn/internal/storage"
)

// deleteFileRequest represents a request to delete a file
type deleteFileRequest struct {
	Path string `json:"path" example:"98ecf8427e/6b141e65-901f-47b7-a1ed-282fd80fc7c6.pdf"`
}

type DeleteFileHandler struct {
	storage storage.Storage
}

func NewDeleteFileHandler(storage storage.Storage) *DeleteFileHandler {
	return &DeleteFileHandler{
		storage: storage,
	}
}

// DeleteFile godoc
// @Summary Удалить файл
// @Description Удаляет файл по пути
// @Tags files
// @Accept json
// @Produce json
// @Param request body deleteFileRequest true "Путь к файлу"
// @Success 200 "Файл успешно удален"
// @Failure 400 "Неверный запрос"
// @Failure 401 "Неавторизован"
// @Failure 500 "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /api/v1/files/delete [delete]
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
