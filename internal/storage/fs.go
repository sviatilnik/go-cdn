package storage

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"os"
	"path/filepath"
	"time"

	"crypto/md5"

	"github.com/google/uuid"
)

type FSStorage struct {
	path string
}

func NewFSStorage(path string) *FSStorage {
	return &FSStorage{
		path: path,
	}
}

func (s *FSStorage) SaveFile(ctx context.Context, file io.Reader, filename string) (*File, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := s.ensureDirectoryExists(); err != nil {
		return nil, err
	}

	ext := filepath.Ext(filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	hash := hex.EncodeToString(md5.New().Sum([]byte(time.Now().Format("200601021504"))))
	hash = string(hash[len(hash)-10:])
	dirPath := filepath.Join(s.path, hash)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, err
	}

	filePath := filepath.Join(dirPath, newFileName)

	outFile, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := outFile.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	_, err = io.Copy(outFile, file)
	if err != nil {
		return nil, err
	}

	info, err := outFile.Stat()
	if err != nil {
		return nil, err
	}

	fileInfo := &File{
		Filename:    info.Name(),
		Path:        filepath.Join(hash, info.Name()),
		Size:        info.Size(),
		ContentType: mime.TypeByExtension(ext),
		Timestamp:   info.ModTime(),
	}

	return fileInfo, nil
}

func (s *FSStorage) GetFile(ctx context.Context, relativePath string) (*File, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	filePath := filepath.Join(s.path, relativePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, ErrFileNotFound
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &File{
		Path:        relativePath,
		Filename:    filepath.Base(filePath),
		Size:        info.Size(),
		ContentType: mime.TypeByExtension(filepath.Ext(filePath)),
		Timestamp:   info.ModTime(),
		File:        f,
	}, nil
}

func (s *FSStorage) DeleteFile(ctx context.Context, relativePath string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	filePath := filepath.Join(s.path, relativePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ErrFileNotFound
	}

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}

func (s *FSStorage) GetFileMaxSize() int64 {
	return 10 << 20
}

func (s *FSStorage) ensureDirectoryExists() error {
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		if err := os.MkdirAll(s.path, 0755); err != nil {
			return err
		}
	}
	return nil
}
