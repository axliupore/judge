package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
)

var Logger *logrus.Logger

func InitLogger() {
	logger := logrus.New()

	logrus.SetLevel(logrus.InfoLevel)

	logger.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			fullPath := f.File
			_, fileName := filepath.Split(fullPath)
			parentDir := filepath.Base(filepath.Dir(fullPath))
			callerInfo := fmt.Sprintf("%s/%s:%d", parentDir, fileName, f.Line)
			return "", callerInfo
		},
	})

	logger.SetReportCaller(true)

	Logger = logger
}
