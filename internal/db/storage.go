package db

import "kv-storage/internal/models"

type Storage interface {
	Create(in *models.KeyValue) (*models.KeyValue, error)
	Get(key string) (*models.KeyValue, error)
	Update(in *models.KeyValue) (*models.KeyValue, error)
	Delete(key string) (*models.KeyValue, error)
}
