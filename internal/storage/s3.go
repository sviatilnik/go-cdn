package storage

import "io"

type S3Storage struct{}

func NewS3Storage() *S3Storage {
	return &S3Storage{}
}

func (s *S3Storage) SaveFile(file io.Reader, filename string) (*File, error) {
	return nil, nil
}

func (s *S3Storage) GetFile(relativePath string) (*File, error) {
	return nil, nil
}

func (s *S3Storage) DeleteFile(filename string) error {
	return nil
}
