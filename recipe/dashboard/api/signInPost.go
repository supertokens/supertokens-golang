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

func SignInPost(apiInterface dashboardmodels.APIInterface, options dashboardmodels.APIOptions, userContext supertokens.UserContext) error {
	body, err := supertokens.ReadFromRequest(options.Req)

	if err != nil {
		return err
	}

	var readBody signInRequestBody
	err = json.Unmarshal(body, &readBody)
	if err != nil {
		return err
	}

	if readBody.Email == nil {
		return supertokens.BadInputError{
			Msg: "Required parameter 'email' is missing",
		}
	}

	if readBody.Password == nil {
		return supertokens.BadInputError{
			Msg: "Required parameter 'password' is missing",
		}
	}

	querier, querierErr := supertokens.GetNewQuerierInstanceOrThrowError("dashboard")

	if querierErr != nil {
		return querierErr
	}

	apiResponse, apiErr := querier.SendPostRequest("/recipe/dashboard/signin", map[string]interface{}{
		"email":    *readBody.Email,
		"password": *readBody.Password,
	})

	if apiErr != nil {
		return apiErr
	}

	status := apiResponse["status"]

	if status == "OK" {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status":    "OK",
			"sessionId": apiResponse["sessionId"].(string),
		})
	}

	if status == "USER_SUSPENDED_ERROR" {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status":  "USER_SUSPENDED_ERROR",
			"message": apiResponse["message"].(string),
		})
	}

	return supertokens.Send200Response(options.Res, map[string]interface{}{
		"status": "INVAlID_CREDENTIALS_ERROR",
	})
}
