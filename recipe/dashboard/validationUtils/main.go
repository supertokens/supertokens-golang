package validationUtils

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"net/http"
	"strings"
)

func ValidateApiKey(req *http.Request, config dashboardmodels.TypeNormalisedInput, usercontext supertokens.UserContext) (bool, error) {
	apiKeyHeaderValue := req.Header.Get("authorization")

	// We receive the api key as `Bearer API_KEY`, this retrieves just the key
	keyParts := strings.Split(apiKeyHeaderValue, " ")
	apiKeyHeaderValue = keyParts[len(keyParts)-1]

	if apiKeyHeaderValue == "" {
		return false, nil
	}

	return apiKeyHeaderValue == *config.ApiKey, nil
}
