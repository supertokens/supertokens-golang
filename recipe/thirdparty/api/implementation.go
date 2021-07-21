package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
)

func MakeAPIImplementation() models.APIImplementation {
	return models.APIImplementation{
		AuthorisationUrlGET: func(provider models.TypeProvider, options models.APIOptions) models.AuthorisationUrlGETResponse {
			providerInfo := provider.Get(nil, nil)
			var params map[string]string
			for key, value := range providerInfo.AuthorisationRedirect.Params {
				// TODO: check for value as function
				params[key] = value
			}
			paramsString := getParamString(providerInfo.AuthorisationRedirect.Params)
			url := providerInfo.AuthorisationRedirect.URL + paramsString
			return models.AuthorisationUrlGETResponse{
				Status: "OK",
				URL:    url,
			}
		},
		SignInUpPOST: func(provider models.TypeProvider, code, redirectURI string, options models.APIOptions) models.SignInUpPOSTResponse {
			providerInfo := provider.Get(&redirectURI, &code)
			accessTokenAPIResponse, err := postRequest(providerInfo)
			if err != nil {
				return models.SignInUpPOSTResponse{
					Status: "FIELD_ERROR",
					Error:  err,
				}
			}
			userInfo, err := providerInfo.GetProfileInfo(accessTokenAPIResponse["data"])
			if err != nil {
				return models.SignInUpPOSTResponse{
					Status: "FIELD_ERROR",
					Error:  err,
				}
			}

			emailInfo := userInfo.Email
			if emailInfo == nil {
				return models.SignInUpPOSTResponse{
					Status: "NO_EMAIL_GIVEN_BY_PROVIDER",
				}
			}

			response := options.RecipeImplementation.SignInUp(provider.ID, userInfo.ID, *emailInfo)
			if response.Status == "FIELD_ERROR" {
				return models.SignInUpPOSTResponse{
					Status: response.Status,
					Error:  response.Error,
				}
			}

			action := "signin"
			if response.CreatedNewUser {
				action = "signup"
			}

			jwtPayload := options.Config.SessionFeature.SetJwtPayload(response.User, accessTokenAPIResponse["data"], action)
			sessionData := options.Config.SessionFeature.SetSessionData(response.User, accessTokenAPIResponse["data"], action)

			_, err = session.CreateNewSession(options.Res, response.User.ID, jwtPayload, sessionData)
			if err != nil {
				return models.SignInUpPOSTResponse{
					Status: "FIELD_ERROR",
					Error:  err,
				}
			}
			return models.SignInUpPOSTResponse{
				Status:           "OK",
				CreatedNewUser:   response.CreatedNewUser,
				User:             response.User,
				AuthCodeResponse: accessTokenAPIResponse["data"],
			}
		},
	}
}

func postRequest(providerInfo models.TypeProviderGetResponse) (map[string]interface{}, error) {
	paramsString := getParamString(providerInfo.AuthorisationRedirect.Params)
	req, err := http.NewRequest("POST", providerInfo.AccessTokenAPI.URL, bytes.NewBuffer([]byte(paramsString)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("accept", "application/json") // few providers like github don't send back json response by default

	client := &http.Client{}
	response, err := client.Do(req)
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getParamString(params map[string]string) string {
	paramsString := "?"
	for key, value := range params {
		paramsString = paramsString + key + "=" + value + "&"
	}
	paramsString = strings.TrimSuffix(paramsString, "&")
	return paramsString
}
