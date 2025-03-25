package db

import "github.com/MosinFAM/tarantool-kv/internal/models"

// go install go.uber.org/mock/mockgen@latest
//
//go:generate mockgen -source=storage.go -destination=storage_mock.go -package=db StorageRepo
type Storage interface {
	Create(in *models.KeyValue) (*models.KeyValue, error)
	Get(key string) (*models.KeyValue, error)
	Update(in *models.KeyValue) (*models.KeyValue, error)
	Delete(key string) (*models.KeyValue, error)
}
