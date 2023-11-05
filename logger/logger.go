package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	Logger *log.Logger
)

func CreateFileLogger(appName, filePath string) *log.Logger {
	logFile := os.Stdout

	if filePath != "" {
		_logFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[logger.createFileLogger]: %s\n", err.Error())
			os.Exit(1)
		} else {
			logFile = _logFile
		}
	}

	logger := log.New(logFile, "", log.Ldate|log.Ltime)
	prefix := fmt.Sprintf("[%s] ", appName)
	logger.SetPrefix(prefix)

	Logger = logger
	return logger
}
