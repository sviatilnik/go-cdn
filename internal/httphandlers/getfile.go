package httphandlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
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

// GetFile godoc
// @Summary Получить файл
// @Description Получает файл по пути
// @Tags files
// @Accept json
// @Produce application/octet-stream
// @Param folder path string true "Папка файла"
// @Param filename path string true "Имя файла"
// @Success 200 "Файл успешно получен"
// @Failure 400 "Неверный запрос"
// @Failure 404 "Файл не найден"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /{folder}/{filename} [get]
func (h *GetFileHandler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		folder := chi.URLParam(r, "folder")
		filename := chi.URLParam(r, "filename")
		if folder == "" || filename == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Folder and filename are required"))
			return
		}

		file, err := h.storage.GetFile(r.Context(), fmt.Sprintf("%s/%s", folder, filename))
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
