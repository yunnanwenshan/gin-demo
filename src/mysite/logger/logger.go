package logger

import (
	"github.com/Sirupsen/logrus"
	"os"
)

var logger = logrus.New()

func init()  {
	file, err := os.OpenFile("logrus.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		logger.Out = file
	} else {
		logger.Info("Failed to log to file, using default stderr")
	}

	logger.Level = logrus.DebugLevel
}

func GetLogger() *logrus.Logger  {
	return logger;
}
