package storage

import (
	"io"
	"time"
)

// File represents a file in the storage
// @Description Информация о файле в хранилище
type File struct {
	// Имя файла
	Filename string `json:"filename" example:"6b141e65-901f-47b7-a1ed-282fd80fc7c6.pdf"`
	// Путь к файлу
	Path string `json:"path" example:"98ecf8427e/6b141e65-901f-47b7-a1ed-282fd80fc7c6.pdf"`
	// Размер файла в байтах
	Size int64 `json:"size" example:"123456"`
	// MIME тип файла
	ContentType string `json:"content_type" example:"application/pdf"`
	// Время создания/модификации файла
	Timestamp time.Time `json:"timestamp" example:"2024-01-15T10:30:00Z"`

	// Файловый дескриптор (не возвращается в JSON)
	File io.ReadSeekCloser `json:"-"`
}
