package logging

import (
	"fmt"
	"os"
	"runtime"

	"github.com/mattmattox/go-smtp-relay/pkg/config"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// SetupLogging initializes the logger with the appropriate settings.
func SetupLogging() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()
		logger.SetOutput(os.Stdout)
		logger.SetReportCaller(true)

		customFormatter := new(logrus.TextFormatter)
		customFormatter.TimestampFormat = "2006-01-02T15:04:05-0700"
		customFormatter.FullTimestamp = true
		customFormatter.CallerPrettyfier = func(frame *runtime.Frame) (function string, file string) {
			// Only show the full file path
			return "", frame.File
		}
		logger.SetFormatter(customFormatter)

		fmt.Printf("Debug flag value: %v\n", config.CFG.Debug)
		if config.CFG.Debug {
			logger.SetLevel(logrus.DebugLevel)
		} else {
			logger.SetLevel(logrus.InfoLevel)
		}

	}

	return logger
}
