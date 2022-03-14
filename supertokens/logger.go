package supertokens

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	iLogger             = log.New(os.Stdout, "com.supertokens:info ", 0)
	eLogger             = log.New(os.Stderr, "com.supertokens:error ", 0)
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
		eLogger.Printf(logMessage((formattedMessage)))
	}

}

func isNamespacePassed(id string) bool {
	namespace, exists := os.LookupEnv(DEBUG_FLAG)

	if exists {
		namespaceWithStar := fmt.Sprintf("%s*", SUPERTOKENS_NAMESPACE)
		namespaceWithId := fmt.Sprintf("%s%s", SUPERTOKENS_NAMESPACE, id)

		if strings.Contains(namespace, namespaceWithStar) || strings.Contains(namespace, namespaceWithId) {
			return true
		}
	}
	return false
}
