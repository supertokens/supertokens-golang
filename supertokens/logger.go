package supertokens

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const SUPERTOKENS_LOGGER_NAMESPACE = "com.supertokens:"

var (
	iLogger             = log.New(os.Stdout, "com.supertokens:info ", 0)
	dLogger             = log.New(os.Stdout, "com.supertokens:debug ", 0)
	DebugLoggerWithCode = map[int]func(string){
		1: func(param string) {
			debugLoggerHelper(1, "API replied with status: "+param)
		},
	}
)

func logMessage(message string) string {
	return fmt.Sprintf("{t: \"%d\", msg: %s, sdkVer: \"%s\"}", time.Now().UnixMilli(), message, VERSION)
}

func InfoLogger(message string) {
	if isNamespacePassed("info") {
		iLogger.Printf(logMessage("\"%s\""), message)
	}
}

func debugLoggerHelper(errorCode int, message string) {
	if isNamespacePassed("debug") {
		formattedMessage := fmt.Sprintf("\"%s\", debugCode: %d", message, errorCode)
		dLogger.Printf(logMessage((formattedMessage)))
	}

}

func isNamespacePassed(id string) bool {
	namespace, exists := os.LookupEnv("SUPERTOKENS_DEBUG")

	if exists {
		namespaceWithStar := fmt.Sprintf("%s*", SUPERTOKENS_LOGGER_NAMESPACE)
		namespaceWithId := fmt.Sprintf("%s%s", SUPERTOKENS_LOGGER_NAMESPACE, id)

		if strings.Contains(namespace, namespaceWithStar) || strings.Contains(namespace, namespaceWithId) {
			return true
		}
	}
	return false
}
