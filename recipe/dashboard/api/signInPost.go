package api

import (
	"encoding/json"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type signInPostResponse struct {
	Status    string `json:"status"`
	SessionId string `json:"sessionId,omitempty"`
	Message   string `json:"message,omitempty"`
}

type signInRequestBody struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func SignInPost(apiInterface dashboardmodels.APIInterface, options dashboardmodels.APIOptions) (signInPostResponse, error) {
	body, err := supertokens.ReadFromRequest(options.Req)

	if err != nil {
		return signInPostResponse{}, err
	}

	var readBody signInRequestBody
	err = json.Unmarshal(body, &readBody)
	if err != nil {
		return signInPostResponse{}, err
	}

	if readBody.Email == nil {
		return signInPostResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'email' is missing",
		}
	}

	if readBody.Password == nil {
		return signInPostResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'password' is missing",
		}
	}

	querier, querierErr := supertokens.GetNewQuerierInstanceOrThrowError("dashboard")

	if querierErr != nil {
		return signInPostResponse{}, querierErr
	}

	apiResponse, apiErr := querier.SendPostRequest("/recipe/dashboard/signin", map[string]interface{}{
		"email":    *readBody.Email,
		"password": *readBody.Password,
	})

	if apiErr != nil {
		return signInPostResponse{}, apiErr
	}

	status, ok := apiResponse["status"]

	if ok && status == "OK" {
		return signInPostResponse{
			Status:    "OK",
			SessionId: apiResponse["sessionId"].(string),
		}, nil
	}

	if status == "USER_SUSPENDED_ERROR" {
		return signInPostResponse{
			Status:  "USER_SUSPENDED_ERROR",
			Message: apiResponse["message"].(string),
		}, nil
	}

	return signInPostResponse{
		Status: "INVAlID_CREDENTIALS_ERROR",
	}, nil
}
