package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/derekstavis/go-qs"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
)

func MakeAPIImplementation() models.APIImplementation {
	return models.APIImplementation{
		AuthorisationUrlGET: func(provider models.TypeProvider, options models.APIOptions) models.AuthorisationUrlGETResponse {
			providerInfo := provider.Get(nil, nil)
			params := make(map[string]string)
			for key, value := range providerInfo.AuthorisationRedirect.Params {
				if reflect.ValueOf(value).Kind() == reflect.String {
					params[key] = value.(string)
				} else {
					call, ok := value.(func(req *http.Request) string)
					if ok {
						params[key] = call(options.Req)
					}
				}

			}
			paramsString, _ := getParamString(params)
			url := providerInfo.AuthorisationRedirect.URL + "?" + paramsString
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
			userInfo, err := providerInfo.GetProfileInfo(accessTokenAPIResponse)
			if err != nil {
				fmt.Println("GetProfileInfo err")
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
				fmt.Println("SignInUp err")
				return models.SignInUpPOSTResponse{
					Status: response.Status,
					Error:  response.Error,
				}
			}

			action := "signin"
			if response.CreatedNewUser {
				action = "signup"
			}

			jwtPayload := options.Config.SessionFeature.SetJwtPayload(response.User, accessTokenAPIResponse, action)
			sessionData := options.Config.SessionFeature.SetSessionData(response.User, accessTokenAPIResponse, action)

			_, err = session.CreateNewSession(options.Res, response.User.ID, jwtPayload, sessionData)
			if err != nil {
				fmt.Println("CreateNewSession err", err)
				return models.SignInUpPOSTResponse{
					Status: "FIELD_ERROR",
					Error:  err,
				}
			}
			return models.SignInUpPOSTResponse{
				Status:           "OK",
				CreatedNewUser:   response.CreatedNewUser,
				User:             response.User,
				AuthCodeResponse: accessTokenAPIResponse,
			}
		},
	}
}

func postRequest(providerInfo models.TypeProviderGetResponse) (map[string]interface{}, error) {
	querystring, err := getParamString(providerInfo.AccessTokenAPI.Params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", providerInfo.AccessTokenAPI.URL, bytes.NewBuffer([]byte(querystring)))
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

func getParamString(paramsMap map[string]string) (string, error) {
	params := make(map[string]interface{})
	for key, value := range paramsMap {
		params[key] = value
	}
	return qs.Marshal(params)
}
