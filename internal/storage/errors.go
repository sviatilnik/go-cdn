package storage

import "errors"

var ErrFileNotFound = errors.New("file not found")
var ErrUnknownStorageType = errors.New("unknown storage type")
