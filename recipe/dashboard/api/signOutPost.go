package api

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"strings"
)

type signOutPostResponse struct {
	Status string `json:"status"`
}

func SignOutPost(apiInterface dashboardmodels.APIInterface, options dashboardmodels.APIOptions) (signOutPostResponse, error) {
	if options.Config.AuthMode == dashboardmodels.AuthModeAPIKey {
		return signOutPostResponse{
			Status: "OK",
		}, nil
	}

	sessionIdFromHeader := options.Req.Header.Get("authorization")

	// We receive the api key as `Bearer API_KEY`, this retrieves just the key
	keyParts := strings.Split(sessionIdFromHeader, " ")
	sessionIdFromHeader = keyParts[len(keyParts)-1]

	querier, querierErr := supertokens.GetNewQuerierInstanceOrThrowError("dashboard")

	if querierErr != nil {
		return signOutPostResponse{}, querierErr
	}

	_, apiError := querier.SendDeleteRequest("/recipe/dashboard/session", map[string]interface{}{}, map[string]string{
		"sessionId": sessionIdFromHeader,
	})

	if apiError != nil {
		return signOutPostResponse{}, apiError
	}

	return signOutPostResponse{
		Status: "OK",
	}, nil
}
