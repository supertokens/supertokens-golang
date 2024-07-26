package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/supertokens"
)

var info = struct {
	coreCallCount int
}{}

func GetFunc(evalStr string) (interface{}, error) {
	if strings.HasPrefix(evalStr, "supertokens.init.supertokens.networkInterceptor") {
		return func(request *http.Request, userContext supertokens.UserContext) (*http.Request, error) {
			info.coreCallCount += 1
			return request, nil
		}, nil
	}

	return nil, fmt.Errorf("Unknown eval string")
}
