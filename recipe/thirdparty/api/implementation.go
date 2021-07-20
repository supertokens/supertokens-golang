package api

import (
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
)

func MakeAPIImplementation() models.APIImplementation {
	return models.APIImplementation{
		AuthorisationUrlGET: func(provider models.TypeProvider, options models.APIOptions) models.AuthorisationUrlGETResponse {
			providerInfo := provider.Get(nil, nil)
			var params map[string]string
			paramsString := "?"
			for key, value := range providerInfo.AuthorisationRedirect.Params {
				params[key] = value
				paramsString = paramsString + key + "=" + value + "&"
			}
			paramsString = strings.TrimSuffix(paramsString, "&")
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
			userInfo := providerInfo.GetProfileInfo(nil)

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

			jwtPayload := options.Config.SessionFeature.SetJwtPayload(response.User, accessTokenAPIResponse, action)
			sessionData := options.Config.SessionFeature.SetSessionData(response.User, accessTokenAPIResponse, action)

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
				AuthCodeResponse: accessTokenAPIResponse,
			}
		},
	}
}

// TODO:
func postRequest(providerInfo models.TypeProviderGetResponse) (string, error) {
	return "", nil
}
