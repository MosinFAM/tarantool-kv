package db

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/MosinFAM/tarantool-kv/internal/logger"
	"github.com/MosinFAM/tarantool-kv/internal/models"

	"github.com/sirupsen/logrus"
	tarantool "github.com/tarantool/go-tarantool"
)

type KeyValueManager struct {
	tConn *tarantool.Connection
}

func NewKeyValueManager(conn *tarantool.Connection) *KeyValueManager {
	return &KeyValueManager{tConn: conn}
}

func ConnectTarantool() (*tarantool.Connection, error) {
	host := os.Getenv("TARANTOOL_HOST")
	port := os.Getenv("TARANTOOL_PORT")

	if host == "" {
		host = "tarantool"
	}
	if port == "" {
		port = "3301"
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	opts := tarantool.Opts{User: "guest"}

	conn, err := tarantool.Connect(addr, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Tarantool: %w", err)
	}

	logger.LogInfo("Connected to Tarantool at", logrus.Fields{"addr": addr})
	return conn, nil
}

// Create добавляет новую пару ключ-значение в Tarantool
func (kv *KeyValueManager) Create(in *models.KeyValue) (*models.KeyValue, error) {
	logger.LogInfo("Start creating key-value", logrus.Fields{"key-value": in})
	dataSerialized, err := json.Marshal(in.Value)
	if err != nil {
		logger.LogError("Data serialization failed", err, logrus.Fields{"key": in.Key})
		return nil, fmt.Errorf("data serialization failed: %w", err)
	}

	_, err = kv.tConn.Call("insert_kv", []interface{}{in.Key, string(dataSerialized)})
	if err != nil {
		re := regexp.MustCompile(`key already exists`)
		if re.MatchString(err.Error()) {
			logger.LogError("Key already exists during insert", err, logrus.Fields{"key": in.Key})
			return nil, fmt.Errorf("key already exists")
		}

		logger.LogError("Failed to insert key", err, logrus.Fields{"key": in.Key})
		return nil, fmt.Errorf("failed to insert key: %w", err)
	}

	logger.LogInfo("Key successfully created", logrus.Fields{"key": in.Key})
	return in, nil
}

// Get получает значение по ключу
func (kv *KeyValueManager) Get(key string) (*models.KeyValue, error) {
	logger.LogInfo("Start getting key", logrus.Fields{"key": key})
	resp, err := kv.tConn.Call("get_kv", []interface{}{key})
	if err != nil {
		logger.LogError("Failed to get key", err, logrus.Fields{"key": key})
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	if len(resp.Data) == 0 {
		logger.LogInfo("Key not found", logrus.Fields{"key": key})
		return nil, fmt.Errorf("key not found")
	}

	firstItem := resp.Data[0].([]interface{})
	if len(firstItem) == 0 || firstItem[0] == nil {
		logger.LogInfo("Value is nil or missing", logrus.Fields{"key": key})
		return nil, fmt.Errorf("key not found")
	}

	rawValue := firstItem[1].(string)
	var value map[string]interface{}
	if err := json.Unmarshal([]byte(rawValue), &value); err != nil {
		logger.LogError("Failed to unmarshal value", err, logrus.Fields{"key": key})
		return nil, fmt.Errorf("failed to deserialize value: %w", err)
	}

	logger.LogInfo("Key successfully getted", logrus.Fields{"key": key, "Value": value})
	return &models.KeyValue{Key: key, Value: value}, nil
}

// Delete удаляет ключ
func (kv *KeyValueManager) Delete(key string) (*models.KeyValue, error) {
	logger.LogInfo("Start deleting key", logrus.Fields{"key": key})
	existing, err := kv.Get(key)
	if err != nil {
		logger.LogInfo("Key not found during delete", logrus.Fields{"key": key})
		return nil, fmt.Errorf("key not found")
	}

	resp, err := kv.tConn.Call("delete_kv", []interface{}{key})
	if err != nil {
		logger.LogError("Failed to delete key", err, logrus.Fields{"key": key})
		return nil, fmt.Errorf("failed to delete key: %w", err)
	}

	if len(resp.Data) == 0 {
		logger.LogInfo("Key not found during delete", logrus.Fields{"key": key})
		return nil, fmt.Errorf("key not found")
	}

	logger.LogInfo("Key successfully deleted", logrus.Fields{"key": key})
	return existing, nil
}

// Update обновляет значение для ключа
func (kv *KeyValueManager) Update(in *models.KeyValue) (*models.KeyValue, error) {
	logger.LogInfo("Start updating key-value", logrus.Fields{"key-value": in})
	dataSerialized, err := json.Marshal(in.Value)
	if err != nil {
		logger.LogError("Data serialization failed during update", err, logrus.Fields{"key": in.Key})
		return nil, fmt.Errorf("data serialization failed: %w", err)
	}

	resp, err := kv.tConn.Call("update_kv", []interface{}{in.Key, string(dataSerialized)})
	if err != nil {
		logger.LogError("Failed to update key", err, logrus.Fields{"key": in.Key})
		return nil, fmt.Errorf("failed to update key: %w", err)
	}

	data := resp.Data[0].([]interface{})
	if data[0] == nil {
		logger.LogInfo("Key not found during update", logrus.Fields{"key": in.Key})
		return nil, fmt.Errorf("key not found")
	}

	logger.LogInfo("Key successfully updated", logrus.Fields{"key": in.Key})
	return in, nil
}
