package logger

import (
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func Init() {
	Logger = logrus.New()

	// Устанавливаем формат вывода
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	Logger.SetLevel(logrus.InfoLevel)
}

func LogInfo(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Info(message)
}

func LogError(message string, err error, fields logrus.Fields) {
	Logger.WithFields(fields).WithError(err).Error(message)
}
