package storage

import (
	"context"

	"github.com/sviatilnik/go-cdn/internal/config"
)

func GetStorage(ctx context.Context, cnf *config.StorageConfig) (Storage, error) {
	if cnf.Type == config.FSStorageType {
		return NewFSStorage(cnf.Path), nil
	}

	if cnf.Type == config.S3StorageType {
		return NewS3Storage(ctx, cnf)
	}

	return nil, ErrUnknownStorageType
}
