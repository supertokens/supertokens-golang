package supertokens

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

const supertokens_namespace = "com.supertokens"

/*
 The debug logger below can be used to log debug messages in the following format
    com.supertokens {t: "2022-03-21T17:10:42+05:30", message: "Test Message", file: "/home/supertokens-golang/supertokens/supertokens.go:51" sdkVer: "0.5.2"}
*/

// Logger interface exposes the Log() method which is implemented by the client.
// We don't need the AddRequestIDToContext() method anymore, the same thing is done from the
// MakeDefaultUserContextFromAPI(), instead we can just take the key from in the config and add it to the UserContext
// in the sdk itself
type Logger interface {
	Log(msg string)
}

var loggerClt Logger

func SetLogger(logger Logger) {
	loggerClt = logger
}

func LogNewDebugMessage(ctx UserContext, message string) {
	_, exists := os.LookupEnv("SUPERTOKENS_DEBUG")
	if exists {
		reqID, ok := (*ctx)[RequestIDKey].(string)
		if !ok {
			loggerClt.Log(newFormatMessage(message, nil))
		} else {
			loggerClt.Log(newFormatMessage(message, &reqID))
		}
	}
}

// newFormatMessage adds requestID to the log if it is non-nil
func newFormatMessage(message string, reqID *string) string {
	_, file, line, _ := runtime.Caller(2)
	if reqID == nil {
		return fmt.Sprintf(" {t: \"%s\", message: \"%s\", file: \"%s:%d\" sdkVer: \"%s\"}\n\n", time.Now().Format(time.RFC3339), message, file, line, VERSION)
	}
	return fmt.Sprintf(" {t: \"%s\", message: \"%s\", file: \"%s:%d\" sdkVer: \"%s\", requestID: \"%s\"}\n\n", time.Now().Format(time.RFC3339), message, file, line, VERSION, *reqID)
}

// Old code left untouched
var (
	logger = log.New(os.Stdout, supertokens_namespace, 0)
)

func formatMessage(message string) string {
	_, file, line, _ := runtime.Caller(2)
	return fmt.Sprintf(" {t: \"%s\", message: \"%s\", file: \"%s:%d\" sdkVer: \"%s\"}\n\n", time.Now().Format(time.RFC3339), message, file, line, VERSION)
}

func LogDebugMessage(message string) {
	_, exists := os.LookupEnv("SUPERTOKENS_DEBUG")
	if exists {
		logger.Printf(formatMessage(message))
	}
}
