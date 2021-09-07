package logger

import (
	"io/ioutil"
	"os"

	"github.com/google/logger"
)

var logs *logger.Logger
var console *logger.Logger

func Init() {
	lf, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	// defer lf.Close()

	logs = logger.Init("PullHandler", false, false, lf)
	console = logger.Init("PullHandler", true, false, ioutil.Discard)
	// defer logs.Close()
}

// LogMessage - function to log colored messages
func LogMessage(level string, message string, args interface{}) {
	switch level {
	case "info":
		console.Infof("\033[1;34m%s -> %+v\033[0m", message, args)
		logs.Infof("%s -> %+v", message, args)
	case "debug":
		console.Warningf("\033[1;33m----------------%s----------------\033[0m", message)
		logs.Warningf("----------------%s----------------", message)
	case "error":
		console.Errorf("\033[1;31m%s -> %+v\033[0m", message, args)
		logs.Errorf("%s -> %+v", message, args)
	default:
		console.Infof("\033[1;36m%s -> %+v", message, args)
		logs.Infof("%s -> %+v", message, args)
	}
}
