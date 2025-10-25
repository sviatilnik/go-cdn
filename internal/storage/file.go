package storage

import (
	"io"
	"time"
)

type File struct {
	Filename    string            `json:"filename"`
	Path        string            `json:"path"`
	Size        int64             `json:"size"`
	ContentType string            `json:"content_type"`
	Timestamp   time.Time         `json:"timestamp"`
	File        io.ReadSeekCloser `json:"-"`
}
