package service

import (
	"github.com/w1nsec/collector/internal/storage/dbstorage"
	"github.com/w1nsec/collector/internal/storage/filestorage"
	"github.com/w1nsec/collector/internal/storage/memstorage"
)

type Service interface {
	Start() error
	Stop() error

	SetupLogger(level string) error
	CheckDB() error

	memstorage.Storage
	filestorage.FileStorageInterface
	dbstorage.DBStorage
}
