package storage

import (
	"context"
	"io"
)

type Storage interface {
	SaveFile(ctx context.Context, file io.Reader, filename string) (*File, error)
	GetFile(ctx context.Context, relativePath string) (*File, error)
	DeleteFile(ctx context.Context, relativePath string) error
	GetFileMaxSize() int64
}
