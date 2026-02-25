package httphandlers

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/sviatilnik/go-cdn/internal/storage"
)

type SaveFileHandler struct {
	storage storage.Storage
}

func NewSaveFileHandler(storage storage.Storage) *SaveFileHandler {
	return &SaveFileHandler{
		storage: storage,
	}
}

// SaveFile godoc
// @Summary Загрузить файл
// @Description Загружает новый файл в хранилище
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл для загрузки"
// @Success 200 {object} storage.File "Информация о загруженном файле"
// @Failure 400 "Неверный запрос"
// @Failure 401 "Неавторизован"
// @Failure 413 "Файл слишком большой"
// @Failure 500 "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /api/v1/files/save [post]
func (h *SaveFileHandler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer func() {
			if err := file.Close(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error(err.Error())
			}
		}()

		if header.Size > h.storage.GetFileMaxSize() {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			_, err = w.Write([]byte(fmt.Sprintf("File size exceeds maximum allowed size: %d bytes", h.storage.GetFileMaxSize())))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error(err.Error())
				return
			}
			return
		}

		fileInfo, err := h.storage.SaveFile(r.Context(), file, header.Filename)
		if err != nil {
			log.Printf("failed to save file: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(fileInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonData)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error(err.Error())
		}
	}
}
