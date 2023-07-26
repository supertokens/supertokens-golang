package validationUtils

import (
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func ValidateApiKey(req *http.Request, config dashboardmodels.TypeNormalisedInput, userContext ...supertokens.UserContext) (bool, error) {
	apiKeyHeaderValue := req.Header.Get("authorization")

	// We receive the api key as `Bearer API_KEY`, this retrieves just the key
	keyParts := strings.Split(apiKeyHeaderValue, " ")
	apiKeyHeaderValue = keyParts[len(keyParts)-1]

	if apiKeyHeaderValue == "" {
		return false, nil
	}

	return apiKeyHeaderValue == config.ApiKey, nil
}
