package main

import (
	"os"

	"github.com/MosinFAM/tarantool-kv/internal/db"
	"github.com/MosinFAM/tarantool-kv/internal/handlers"

	"github.com/MosinFAM/tarantool-kv/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Инициализация логирования
	logger.Init()

	conn, err := db.ConnectTarantool()
	if err != nil {
		logger.LogError("Failed to connect to Tarantool", err, logrus.Fields{
			"host": os.Getenv("TARANTOOL_HOST"),
			"port": os.Getenv("TARANTOOL_PORT"),
		})
		os.Exit(1)
	}

	kvManager := db.NewKeyValueManager(conn)
	handler := handlers.NewHandler(kvManager)

	r := gin.Default()

	r.POST("/kv", handler.CreateKeyValue)
	r.PUT("/kv/:id", handler.UpdateKeyValue)
	r.GET("/kv/:id", handler.GetKeyValue)
	r.DELETE("/kv/:id", handler.DeleteKeyValue)

	if err := r.Run(":8080"); err != nil {
		logger.LogError("Failed to start server", err, nil)
		os.Exit(1)
	}
}
