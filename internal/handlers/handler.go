package handlers

import (
	"net/http"

	"github.com/MosinFAM/tarantool-kv/internal/db"
	"github.com/MosinFAM/tarantool-kv/internal/logger"
	"github.com/MosinFAM/tarantool-kv/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const keyNotFoundError = "key not found"

type Handler struct {
	storage db.Storage
}

func NewHandler(storage db.Storage) *Handler {
	return &Handler{storage: storage}
}

// CreateKeyValue создает новый ключ-значение
func (h *Handler) CreateKeyValue(c *gin.Context) {
	var request models.KeyValue

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.LogError("Invalid request body", err, logrus.Fields{"body": c.Request.Body})
		c.JSON(http.StatusBadRequest, models.Response{
			Error: "Invalid body",
		})
		return
	}

	if request.Key == "" {
		logger.LogInfo("Key is required", logrus.Fields{"key": request.Key})
		c.JSON(http.StatusBadRequest, models.Response{
			Error: "Key is required",
		})
		return
	}

	// Проверка на пустое значение
	if len(request.Value) == 0 {
		logger.LogInfo("Value must be a non-empty object", logrus.Fields{"key": request.Key})
		c.JSON(http.StatusBadRequest, models.Response{
			Error: "Value must be a non-empty object",
		})
		return
	}

	createdItem, err := h.storage.Create(&request)
	if err != nil {
		if err.Error() == "key already exists" {
			logger.LogError("Key already exists", err, logrus.Fields{"key": request.Key})
			c.JSON(http.StatusConflict, models.Response{
				Error: "Key already exists",
			})
			return
		}

		logger.LogError("Error creating key", err, logrus.Fields{"key": request.Key})
		c.JSON(http.StatusInternalServerError, models.Response{
			Error: "Internal server error",
		})
		return
	}

	logger.LogInfo("Created key successfully", logrus.Fields{"key": request.Key})
	c.JSON(http.StatusOK, models.Response{
		Result:  createdItem,
		Message: "Key created successfully",
	})
}

// GetKeyValue получает значение для ключа
func (h *Handler) GetKeyValue(c *gin.Context) {
	key := c.Param("id")

	gettedItem, err := h.storage.Get(key)
	if err != nil {
		logger.LogError("Error getting key", err, logrus.Fields{"key": key})
		if err.Error() == keyNotFoundError {
			c.JSON(http.StatusNotFound, models.Response{
				Error: keyNotFoundError,
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Error: "Internal server error",
			})
		}
		return
	}

	logger.LogInfo("Fetched key successfully", logrus.Fields{"key": key})
	c.JSON(http.StatusOK, models.Response{
		Result:  gettedItem,
		Message: "Key getted successfully",
	})
}

// DeleteKeyValue удаляет ключ
func (h *Handler) DeleteKeyValue(c *gin.Context) {
	key := c.Param("id")

	deletedItem, err := h.storage.Delete(key)
	if err != nil {
		logger.LogError("Error deleting key", err, logrus.Fields{"key": key})
		if err.Error() == keyNotFoundError {
			c.JSON(http.StatusNotFound, models.Response{
				Error: keyNotFoundError,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.Response{
			Error: "Internal server error",
		})
		return
	}

	logger.LogInfo("Deleted key successfully", logrus.Fields{"key": key})
	c.JSON(http.StatusOK, models.Response{
		Deleted: deletedItem,
		Message: "Key deleted successfully",
	})
}

// UpdateKeyValue обновляет значение для ключа
func (h *Handler) UpdateKeyValue(c *gin.Context) {
	var request models.KeyValue
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.LogError("Invalid request body", err, logrus.Fields{"body": c.Request.Body})
		c.JSON(http.StatusBadRequest, models.Response{
			Error: "Invalid body",
		})
		return
	}

	key := c.Param("id")
	request.Key = key

	updatedItem, err := h.storage.Update(&request)
	if err != nil {
		logger.LogError("Error updating key", err, logrus.Fields{"key": key})
		if err.Error() == keyNotFoundError {
			c.JSON(http.StatusNotFound, models.Response{
				Error: keyNotFoundError,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.Response{
			Error: "Internal server error",
		})
		return
	}

	logger.LogInfo("Updated key successfully", logrus.Fields{"key": key})
	c.JSON(http.StatusOK, models.Response{
		Result:  updatedItem,
		Message: "Key updated successfully",
	})
}
