package service

import (
	"github.com/w1nsec/collector/internal/storage"
	"github.com/w1nsec/collector/internal/storage/filestorage"
)

type Service interface {
	Start() error
	Stop() error

	SetupLogger(level string) error
	CheckDB() error

	storage.Storage
	filestorage.FileStorageInterface
}
