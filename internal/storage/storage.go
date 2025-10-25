package storage

import (
	"io"
)

type Storage interface {
	SaveFile(file io.Reader, filename string) (*File, error)
	GetFile(relativePath string) (*File, error)
	DeleteFile(relativePath string) error
	GetFileMaxSize() int64
}
