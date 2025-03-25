package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MosinFAM/tarantool-kv/internal/db"
	"github.com/MosinFAM/tarantool-kv/internal/handlers"
	"github.com/MosinFAM/tarantool-kv/internal/logger"
	"github.com/MosinFAM/tarantool-kv/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func setupTest(t *testing.T) (*handlers.Handler, *db.MockStorage, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	if ctrl == nil {
		t.Fatal("Failed to create gomock Controller")
	}

	mockStorage := db.NewMockStorage(ctrl)
	if mockStorage == nil {
		t.Fatal("Failed to create mockStorage")
	}

	logger.Init()
	h := handlers.NewHandler(mockStorage)
	if h == nil {
		t.Fatal("handlers.NewHandler returned nil")
	}

	gin.SetMode(gin.TestMode)
	return h, mockStorage, ctrl
}

func TestCreateKeyValue_Success(t *testing.T) {
	h, mockStorage, ctrl := setupTest(t)
	t.Logf("Handler: %+v", h)
	t.Logf("MockStorage: %+v", mockStorage)

	defer ctrl.Finish()

	validRequest := models.KeyValue{
		Key:   "testKey",
		Value: map[string]interface{}{"data": "testValue"},
	}

	mockStorage.EXPECT().Create(gomock.Any()).Return(&validRequest, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	requestBody, _ := json.Marshal(validRequest)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")
	t.Logf("Request body: %s", requestBody) // Логируем тело запроса
	h.CreateKeyValue(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestCreateKeyValue_Conflict(t *testing.T) {
	h, mockStorage, ctrl := setupTest(t)
	defer ctrl.Finish()

	validRequest := models.KeyValue{
		Key:   "testKey",
		Value: map[string]interface{}{"data": "testValue"},
	}

	mockStorage.EXPECT().Create(gomock.Any()).Return(nil, errors.New("key already exists"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	requestBody, _ := json.Marshal(validRequest)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateKeyValue(c)

	if w.Code != http.StatusConflict {
		t.Errorf("expected status 409, got %d", w.Code)
	}
}

func TestGetKeyValue_Success(t *testing.T) {
	h, mockStorage, ctrl := setupTest(t)
	defer ctrl.Finish()

	validRequest := models.KeyValue{
		Key:   "testKey",
		Value: map[string]interface{}{"data": "testValue"},
	}

	mockStorage.EXPECT().Get("testKey").Return(&validRequest, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "testKey"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/testKey", nil)

	h.GetKeyValue(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGetKeyValue_NotFound(t *testing.T) {
	h, mockStorage, ctrl := setupTest(t)
	defer ctrl.Finish()

	mockStorage.EXPECT().Get("missingKey").Return(nil, errors.New("key not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "missingKey"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/missingKey", nil)

	h.GetKeyValue(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestUpdateKeyValue_Success(t *testing.T) {
	h, mockStorage, ctrl := setupTest(t)
	defer ctrl.Finish()

	validRequest := models.KeyValue{
		Key:   "testKey",
		Value: map[string]interface{}{"data": "testValue"},
	}

	mockStorage.EXPECT().Update(gomock.Any()).Return(&validRequest, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "testKey"}}
	requestBody, _ := json.Marshal(validRequest)
	c.Request = httptest.NewRequest(http.MethodPut, "/testKey", bytes.NewBuffer(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	h.UpdateKeyValue(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestUpdateKeyValue_NotFound(t *testing.T) {
	h, mockStorage, ctrl := setupTest(t)
	defer ctrl.Finish()

	mockStorage.EXPECT().Update(gomock.Any()).Return(nil, errors.New("key not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "missingKey"}}
	requestBody, _ := json.Marshal(models.KeyValue{Key: "missingKey", Value: map[string]interface{}{"data": "newValue"}})
	c.Request = httptest.NewRequest(http.MethodPut, "/missingKey", bytes.NewBuffer(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	h.UpdateKeyValue(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestDeleteKeyValue_Success(t *testing.T) {
	h, mockStorage, ctrl := setupTest(t)
	defer ctrl.Finish()

	validKey := "testKey"
	validRequest := models.KeyValue{
		Key:   validKey,
		Value: map[string]interface{}{"data": "testValue"},
	}

	// Ожидаем вызов Delete с аргументом validKey и возвращаем успешный результат
	mockStorage.EXPECT().Delete(validKey).Return(&validRequest, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: validKey}}
	c.Request = httptest.NewRequest(http.MethodDelete, "/testKey", nil)

	h.DeleteKeyValue(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestDeleteKeyValue_NotFound(t *testing.T) {
	h, mockStorage, ctrl := setupTest(t)
	defer ctrl.Finish()

	invalidKey := "missingKey"

	// Ожидаем вызов Delete с аргументом invalidKey и возвращаем ошибку
	mockStorage.EXPECT().Delete(invalidKey).Return(nil, errors.New("key not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: invalidKey}}
	c.Request = httptest.NewRequest(http.MethodDelete, "/missingKey", nil)

	h.DeleteKeyValue(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestCreateKeyValue_InvalidBody(t *testing.T) {
	h, _, ctrl := setupTest(t)
	defer ctrl.Finish()

	// Некорректный JSON в теле запроса
	invalidJSON := `{"key": "testKey", "value": }` // Ошибка в JSON

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateKeyValue(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	// Проверяем, что в ответе есть сообщение об ошибке
	var response models.Response
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Error != "Invalid body" {
		t.Errorf("expected error message 'Invalid body', got '%s'", response.Error)
	}
}

func TestCreateKeyValue_MissingKey(t *testing.T) {
	h, _, ctrl := setupTest(t)
	defer ctrl.Finish()

	// Запрос без ключа
	invalidRequest := models.KeyValue{
		Value: map[string]interface{}{"data": "testValue"},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	requestBody, _ := json.Marshal(invalidRequest)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateKeyValue(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	// Проверяем, что в ответе есть сообщение об ошибке
	var response models.Response
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Error != "Key is required" {
		t.Errorf("expected error message 'Key is required', got '%s'", response.Error)
	}
}

func TestCreateKeyValue_EmptyValue(t *testing.T) {
	h, _, ctrl := setupTest(t)
	defer ctrl.Finish()

	// Запрос с пустым значением
	invalidRequest := models.KeyValue{
		Key:   "testKey",
		Value: map[string]interface{}{},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	requestBody, _ := json.Marshal(invalidRequest)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateKeyValue(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	// Проверяем, что в ответе есть сообщение об ошибке
	var response models.Response
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Error != "Value must be a non-empty object" {
		t.Errorf("expected error message 'Value must be a non-empty object', got '%s'", response.Error)
	}
}

func TestCreateKeyValue_InternalServerError(t *testing.T) {
	h, mockStorage, ctrl := setupTest(t)
	defer ctrl.Finish()

	// Запрос с правильными данными
	validRequest := models.KeyValue{
		Key:   "testKey",
		Value: map[string]interface{}{"data": "testValue"},
	}

	// Смоделируем ошибку при создании
	mockStorage.EXPECT().Create(gomock.Any()).Return(nil, errors.New("some internal error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	requestBody, _ := json.Marshal(validRequest)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateKeyValue(c)

	// Проверяем статус ответа
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	// Проверяем, что в ответе есть сообщение об ошибке
	var response models.Response
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Error != "Internal server error" {
		t.Errorf("expected error message 'Internal server error', got '%s'", response.Error)
	}
}

func TestCreateKeyValue_InvalidBody2(t *testing.T) {
	h, _, ctrl := setupTest(t)
	defer ctrl.Finish()

	// Некорректный JSON в теле запроса
	invalidJSON := `{"key": "testKey", "value": }` // Ошибка в JSON

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateKeyValue(c)

	// Проверяем статус ответа
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	// Проверяем, что в ответе есть сообщение об ошибке
	var response models.Response
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Error != "Invalid body" {
		t.Errorf("expected error message 'Invalid body', got '%s'", response.Error)
	}
}
