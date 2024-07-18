package log

import "github.com/sirupsen/logrus"

var Logger *logrus.Logger

func InitLogger() {
	logger := logrus.New()

	logrus.SetLevel(logrus.InfoLevel)

	logger.SetFormatter(&logrus.JSONFormatter{})

	Logger = logger
}
