package supertokens

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

const SUPERTOKENS_LOGGER_NAMESPACE = "com.supertokens:"

type LoggerCodes int

const (
	API_RESPONSE LoggerCodes = iota + 1
	API_CALLED
)

var (
	iLogger             = log.New(os.Stdout, "com.supertokens:info ", 0)
	dLogger             = log.New(os.Stdout, "com.supertokens:debug ", 0)
	DebugLoggerWithCode = map[int]func(...string){
		int(API_RESPONSE): func(input ...string) {
			debugLoggerHelper(int(API_RESPONSE), input[0]+" replied with status: "+input[1])
		},
	}
	InfoLoggerWithCode = map[int]func(...string){
		int(API_CALLED): func(input ...string) {
			infoLoggerHelper(int(API_CALLED), input[0]+" was called")
		},
	}
)

func logMessage(message string, code *int) string {
	_, file, line, _ := runtime.Caller(3)
	messageWithOptionalCode := fmt.Sprintf("msg: \"%s\"", message)
	if code != nil {
		messageWithOptionalCode = fmt.Sprintf("%s, code: %d", messageWithOptionalCode, *code)
	}
	return fmt.Sprintf("{t: \"%d\", %s, file: \"%s:%d\" sdkVer: \"%s\"}", time.Now().UnixMilli(), messageWithOptionalCode, file, line, VERSION)
}

func infoLoggerHelper(code int, message string) {
	if isNamespacePassed("info") {
		iLogger.Printf(logMessage(message, &code))
	}
}

func debugLoggerHelper(code int, message string) {
	if isNamespacePassed("debug") {
		dLogger.Printf(logMessage(message, &code))
	}

}

func isNamespacePassed(id string) bool {
	namespace, exists := os.LookupEnv("SUPERTOKENS_DEBUG")
	// if the SUPERTOKENS_DEBUG env variable is passed and it contains either com.supertokens:* or com.supertokens:{id} where id is the id of the logger, then the logger should execute the log.
	if exists {
		namespaceWithStar := fmt.Sprintf("%s*", SUPERTOKENS_LOGGER_NAMESPACE)
		namespaceWithId := fmt.Sprintf("%s%s", SUPERTOKENS_LOGGER_NAMESPACE, id)

		if strings.Contains(namespace, namespaceWithStar) || strings.Contains(namespace, namespaceWithId) {
			return true
		}
	}
	return false
}
