package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func SetupLogger() {
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.DebugLevel)
}
