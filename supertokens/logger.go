package supertokens

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	logger = log.New(os.Stdout, "com.supertokens ", 0)
)

func formatMessage(message string) string {
	_, file, line, _ := runtime.Caller(3)
	return fmt.Sprintf("{t: \"%d\", message: \"%s\", file: \"%s:%d\" sdkVer: \"%s\"}", time.Now().UnixMilli(), message, file, line, VERSION)
}

func LogDebugMessage(message string) {
	namespace, exists := os.LookupEnv("SUPERTOKENS_DEBUG")
	if exists {
		if strings.Contains(namespace, "com.supertokens") {
			logger.Printf(formatMessage(message))
		}

	}
}
